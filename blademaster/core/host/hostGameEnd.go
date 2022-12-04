package host

import (
	"math/rand"
	"net"
	"sync"
	"time"

	. "github.com/KouKouChan/CSO2-Server/blademaster/Exp"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/inventory"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/room"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

var (
	randSeed     = 0
	randSeedLock sync.Mutex
)

func OnHostGameEnd(p *PacketData, client net.Conn) {
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : A user request to send GameEnd but not in server!")
		return
	}
	//找到玩家的房间
	rm := GetRoomFromID(uPtr.GetUserChannelServerID(),
		uPtr.GetUserChannelID(),
		uPtr.GetUserRoomID())
	if rm == nil ||
		rm.Id <= 0 {
		DebugInfo(2, "Error : User", uPtr.UserName, "try to send GameEnd but in a null room !")
		return
	}
	//是不是房主
	if rm.HostUserID != uPtr.Userid {
		DebugInfo(2, "Error : User", uPtr.UserName, "try to send GameEnd but isn't host !")
		return
	}
	//修改房间信息
	rm.SetStatus(StatusWaiting)
	header := BuildGameResultHeader(rm)

	for _, v := range rm.Users {
		//修改用户状态
		v.SetUserStatus(UserNotReady)
		v.AddMatches()
		//发送房间状态
		rst := BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeRoom), BuildRoomSetting(rm, 0xFFFFFFFFFFFFFFFF))
		SendPacket(rst, v.CurrentConnection)
		//检查是否还在游戏内
		if v.CurrentIsIngame {
			//发送游戏结束数据包
			rst = BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeHost), BuildHostStop())
			SendPacket(rst, v.CurrentConnection)
			//发送游戏战绩
			var boxids []uint32
			bl, numBox := canGetBox(rm, v)
			if bl {
				for i := 0; i < numBox; i++ {
					boxid := GetRandomBox()
					boxids = append(boxids, boxid)
					idx := v.AddItem(boxid, 1, 0)
					rst = BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeInventory_Create),
						BuildInventoryInfoSingle(v, 0, idx))
					SendPacket(rst, v.CurrentConnection)
				}
			}
			rst = BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeRoom), header, BuildGameResult(v, boxids))
			SendPacket(rst, v.CurrentConnection)
			DebugInfo(2, "Sent game result to User", v.UserName)
			//修改用户状态
			v.SetUserIngame(false)
		}
	}
	//给每个人发送房间内所有人的准备状态
	for _, v := range rm.Users {
		rst := BuildUserReadyStatus(v)
		//UserInfo部分
		rst = BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeUserInfo), BuildUserInfo(0XFFFFFFFF, NewUserInfo(v), v.Userid, true))
		SendPacket(rst, v.CurrentConnection)
		for _, k := range rm.Users {
			rst = BytesCombine(BuildHeader(k.CurrentSequence, PacketTypeRoom), rst)
			SendPacket(rst, k.CurrentConnection)
		}
	}

	rm.ResetRoomKillNum()
	rm.ResetRoomScore()
	rm.ResetRoomWinner()
}

func BuildHostStop() []byte {
	return []byte{HostStop}
}
func BuildGameResult(u *User, boxids []uint32) []byte {
	buf := make([]byte, 128)
	offset := 0
	WriteUint64(&buf, u.CurrentExp+LevelExpTotal[u.Level-1], &offset) //now total EXP
	WriteUint64(&buf, u.Points, &offset)                              //now total point
	WriteUint8(&buf, 0, &offset)                                      //unk18
	WriteUint8(&buf, 0, &offset)                                      // str len
	WriteUint8(&buf, 0, &offset)                                      // str len
	//WriteString(&buf, []byte("Good"), &offset)
	//WriteString(&buf, []byte("Good"), &offset)
	WriteUint8(&buf, uint8(len(boxids)), &offset) //num of gifts
	for _, v := range boxids {
		WriteUint32(&buf, v, &offset) //item id
		WriteUint16(&buf, 1, &offset) //item count
		WriteUint64(&buf, 0, &offset) //unk22
		WriteUint16(&buf, 0, &offset) //unk23 ，maybe 2 bytes
	}
	WriteUint8(&buf, 0, &offset)  //unk24
	WriteUint16(&buf, 0, &offset) //unk25 ，maybe 2 bytes
	return buf[:offset]
}

func BuildGameResultHeader(rm *Room) []byte {
	buf := make([]byte, 30)
	offset := 0
	WriteUint8(&buf, OUTSetGameResult, &offset)
	WriteUint8(&buf, 0, &offset)                     //unk01
	WriteUint8(&buf, rm.Setting.GameModeID, &offset) //game mod？ 0x02 0x01
	switch rm.Setting.GameModeID {
	case ModeOriginal, ModePig, ModeGiant:
		WriteUint8(&buf, rm.WinnerTeam, &offset) //winner team？ 0x02 ，生化模式貌似没有？
		WriteUint8(&buf, rm.CtScore, &offset)    //CT winNum
		WriteUint8(&buf, rm.TrScore, &offset)    //TR winNum
		WriteUint8(&buf, 0, &offset)             //上半局CT winNum，开启阵营互换情况
		WriteUint8(&buf, 0, &offset)             //上半局TR winNum
	case ModeStealth:
		WriteUint8(&buf, rm.WinnerTeam, &offset) //winner team？ 0x02 ，生化模式貌似没有？
		WriteUint8(&buf, rm.CtScore, &offset)    //CT winNum
		WriteUint8(&buf, rm.TrScore, &offset)    //TR winNum
		WriteUint8(&buf, 0, &offset)             //上半场CT winNum?
		WriteUint8(&buf, 0, &offset)             //上半场TR winNum?
	case ModeDeathmatch, ModeTeamdeath, ModeTeamdeath_mutation, ModeEventmod01:
		WriteUint8(&buf, rm.WinnerTeam, &offset) //winner team？ 0x02 ，生化模式貌似没有？
		WriteUint32(&buf, rm.CtKillNum, &offset) //CT killnum
		WriteUint32(&buf, rm.TrKillNum, &offset) //TR killnum
		WriteUint64(&buf, 0, &offset)            //unk02
	case ModeGhost, ModeTag:
		WriteUint8(&buf, 0, &offset)
		WriteUint32(&buf, 0, &offset)
	case ModeZombie, ModeZombiecraft, ModeZombie_commander, ModeZombie_prop, ModeZombie_zeta, ModeZ_scenario, ModeZ_scenario_side,
		ModeHide, ModeHide2, ModeHide_Item, ModeHide_ice, ModeHide_match, ModeHide_multi, ModeHide_origin:
	case ModeHeroes, ModeZd_boss1, ModeZd_boss2, ModeZd_boss3:
		WriteUint8(&buf, rm.WinnerTeam, &offset)
	default:
		WriteUint8(&buf, rm.WinnerTeam, &offset) //winner team？ 0x02 ，生化模式貌似没有？
		WriteUint8(&buf, rm.CtScore, &offset)    //CT winNum
		WriteUint8(&buf, rm.TrScore, &offset)    //TR winNum
		WriteUint16(&buf, 0, &offset)            //unk00
	}
	WriteUint8(&buf, uint8(rm.GetNumOfRealIngamePlayers()), &offset) //usernum？
	buf = buf[:offset]
	idx := 1
	for _, v := range rm.Users {
		if v.CurrentIsIngame {
			temp := make([]byte, 100)
			offset = 0
			WriteUint32(&temp, v.Userid, &offset) //userid

			WriteUint8(&temp, 0, &offset)                   //unk02
			WriteUint64(&temp, 0, &offset)                  //unk03
			WriteUint64(&temp, 0, &offset)                  //unk04
			WriteUint16(&temp, v.CurrentKillNum, &offset)   //killnum
			WriteUint16(&temp, v.CurrentAssistNum, &offset) //assistnum
			WriteUint16(&temp, v.CurrentDeathNum, &offset)  //deathnum
			WriteUint16(&temp, 0, &offset)                  //unk05 ，maybe 2 bytes
			WriteUint16(&temp, 0, &offset)                  //unk06 ，maybe 2 bytes 0x56 = 86
			WriteUint16(&temp, 0, &offset)                  //unk07 ，maybe 2 bytes 0x2b = 43

			gainexp := GetGainExp(v, rm.Setting.AreBotsEnabled)
			WriteUint64(&temp, gainexp, &offset) //gained EXP
			WriteUint32(&temp, 0, &offset)       //unk08 ，maybe 4 bytes
			WriteUint16(&temp, 0, &offset)       //unk09 ，maybe 2 bytes
			WriteUint8(&temp, 0, &offset)        //unk10 ，maybe 1 bytes

			points := GetGainPoints(v, rm.Setting.AreBotsEnabled)
			WriteUint64(&temp, points, &offset)        //gained point
			WriteUint32(&temp, 0, &offset)             //unk11 ，maybe 4 bytes
			WriteUint16(&temp, 0, &offset)             //unk12 ，maybe 2 bytes
			WriteUint8(&temp, 0, &offset)              //unk13 ，maybe 1 bytes
			WriteUint8(&temp, uint8(v.Level), &offset) //current level ?

			v.GetExp(gainexp)

			v.GetPoints(points)
			if rm.Setting.GameModeID != ModeZ_scenario &&
				rm.Setting.GameModeID != ModeZ_scenario_side &&
				rm.Setting.GameModeID != ModeHeroes {

				v.GetKills(uint32(v.CurrentKillNum))

			}
			v.GetDeathes(uint32(v.CurrentDeathNum))
			v.GetAssists(uint32(v.CurrentAssistNum))
			if v.CurrentTeam == rm.WinnerTeam {
				v.AddWins()
			}

			WriteUint8(&temp, uint8(v.Level), &offset)    //next level ？
			WriteUint8(&temp, 0, &offset)                 //unk15
			WriteUint8(&temp, uint8(idx), &offset)        //rank
			WriteUint16(&temp, v.CurrentKillNum, &offset) //连续击杀数
			WriteUint32(&temp, 0, &offset)                //unk16 ，maybe 4 bytes
			WriteUint8(&temp, v.CurrentTeam, &offset)     //user team
			switch rm.Setting.GameModeID {
			case ModeOriginal, ModePig, ModeGiant:
				WriteUint32(&temp, 0, &offset) //unk17
			case ModeDeathmatch, ModeTeamdeath, ModeTeamdeath_mutation, ModeEventmod01:
			case ModeStealth:
				WriteUint16(&temp, 0, &offset) //unk17
			case ModeZombie, ModeZombiecraft, ModeZombie_commander, ModeZombie_prop, ModeZombie_zeta, ModeZ_scenario, ModeZ_scenario_side, ModeHeroes,
				ModeHide, ModeHide2, ModeHide_Item, ModeHide_ice, ModeHide_match, ModeHide_multi, ModeHide_origin, ModeZd_boss1,
				ModeZd_boss2, ModeZd_boss3, ModeTag:
			default:
				WriteUint32(&temp, 0, &offset) //unk17,貌似有时候不用
			}
			idx++
			buf = BytesCombine(buf, temp[:offset])
		}
	}
	//log.Println(buf)
	return buf
}

func GetGainExp(u *User, bot uint8) uint64 {
	if bot == 0 {
		exp := uint64(u.CurrentKillNum*80 + u.CurrentAssistNum*10 - u.CurrentDeathNum*50)
		if exp > 100 {
			return exp
		}
		return 100
	}
	return 100
}

func GetGainPoints(u *User, bot uint8) uint64 {
	points := uint64(u.CurrentKillNum*80 + u.CurrentAssistNum*20 - u.CurrentDeathNum*40)
	if bot == 0 {
		if points > 400 && points <= 30000 {
			return points
		} else if points > 30000 {
			return 30000
		}
		return 400
	}
	if points > 400 && points <= 25000 {
		return points
	} else if points > 25000 {
		return 25000
	}
	return 400
}

func GetRandomBox() uint32 {
	randSeedLock.Lock()
	rand.Seed(int64(randSeed) + time.Now().UnixNano())
	idx := rand.Intn(len(BoxIDs))
	randSeed++
	if randSeed > 10000 {
		randSeed = 0
	}
	randSeedLock.Unlock()
	return BoxIDs[idx]
}

func canGetBox(rm *Room, u *User) (bool, int) {
	switch rm.Setting.GameModeID {
	case ModeZ_scenario, ModeZ_scenario_side, ModeHeroes:
		if rm.WinnerTeam == u.GetUserTeam() {
			return true, 2
		}
		return false, 0
	case ModeZombie, ModeZombiecraft, ModeZombie_commander, ModeZombie_prop, ModeZombie_zeta:
		if rm.NumPlayers >= 3 {
			return true, 1
		}
		return false, 0
	case ModeOriginal, ModePig, ModeGiant, ModeDeathmatch, ModeTeamdeath, ModeTeamdeath_mutation,
		ModeStealth, ModeHide, ModeHide2, ModeHide_Item, ModeHide_ice, ModeHide_match, ModeHide_multi,
		ModeHide_origin:
		if rm.NumPlayers >= 3 {
			return true, 1
		}
	default:
		return false, 0
	}
	return false, 0
}
