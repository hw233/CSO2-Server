package room

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnLeaveRoom(client net.Conn, end bool) {
	//找到玩家
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		return
	}
	//找到玩家的房间
	rm := GetRoomFromID(uPtr.GetUserChannelServerID(),
		uPtr.GetUserChannelID(),
		uPtr.GetUserRoomID())
	if rm == nil ||
		rm.Id <= 0 {
		return
	}
	//检查玩家游戏状态，准备情况下并且开始倒计时了，那么就不允许离开房间
	if uPtr.IsUserReady() &&
		rm.IsGlobalCountdownInProgress() {
		DebugInfo(2, "Error : User", uPtr.UserName, "try to leave room but is started !")
		return
	}
	//房间移除玩家
	rm.RoomRemoveUser(uPtr.Userid)
	//检查房间是否为空
	if rm.NumPlayers <= 0 {
		DelChannelRoom(rm.Id,
			uPtr.GetUserChannelID(),
			uPtr.GetUserChannelServerID())

	} else {
		//向其他玩家发送离开信息
		SentUserLeaveMes(uPtr, rm)
	}
	// //扣除1000points
	// if uPtr.CurrentIsIngame {
	// 	uPtr.PunishPoints()
	// 	OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, GAME_ROOM_LEAVE_EARLY)
	// 	//UserInfo部分
	// 	rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeUserInfo), BuildUserInfo(0XFFFFFFFF, NewUserInfo(uPtr), uPtr.Userid, true))
	// 	SendPacket(rst, uPtr.CurrentConnection)
	// }
	//设置玩家状态
	uPtr.QuitRoom()
	//发送房间列表给玩家
	if !end {
		OnBroadcastRoomList(uPtr.GetUserChannelServerID(), uPtr.GetUserChannelID(), uPtr)
	}
	//房间状态
	rm.CheckIngameStatus()
	DebugInfo(2, "User", uPtr.UserName, "id", uPtr.Userid, "left room", string(rm.Setting.RoomName), "id", rm.Id)
}
func SentUserLeaveMes(uPtr *User, rm *Room) {
	//发送离开消息
	if rm.HostUserID == uPtr.Userid {
		//选出新房主
		for _, v := range rm.Users {
			rm.SetRoomHost(v)

			DebugInfo(2, "Set User", v.UserName, "id", v.Userid, "to host in room", string(rm.Setting.RoomName), "id", rm.Id)

			if !v.CurrentIsIngame {
				v.SetUserStatus(UserNotReady)
				temp := BuildUserReadyStatus(v)
				for _, k := range rm.Users {
					rst := BytesCombine(BuildHeader(k.CurrentSequence, PacketTypeRoom), temp)
					SendPacket(rst, k.CurrentConnection)
				}
			}
			break
		}
		sethost := BuildSetHost(rm.HostUserID, 1)

		numInGame := 0
		leave := BuildUserLeave(uPtr.Userid)
		//发送数据包
		for _, v := range rm.Users {
			if v.CurrentIsIngame {
				numInGame++
			}
			rst1 := append(BuildHeader(v.CurrentSequence, PacketTypeRoom), OUTPlayerLeave)
			rst1 = BytesCombine(rst1, leave)
			SendPacket(rst1, v.CurrentConnection)

			rst2 := append(BuildHeader(v.CurrentSequence, PacketTypeRoom), OUTSetHost)
			rst2 = BytesCombine(rst2, sethost)
			SendPacket(rst2, v.CurrentConnection)

		}

		if numInGame == 0 {
			rm.SetStatus(StatusWaiting)
		}
		DebugInfo(2, "Sent a set roomHost packet to other users")
		return
	} else {
		leave := BuildUserLeave(uPtr.Userid)
		for _, v := range rm.Users {
			rst1 := append(BuildHeader(v.CurrentSequence, PacketTypeRoom), OUTPlayerLeave)
			rst1 = BytesCombine(rst1, leave)
			SendPacket(rst1, v.CurrentConnection)
		}
		DebugInfo(2, "Sent a leave room packet to other users")
		return
	}
}

func BuildUserLeave(id uint32) []byte {
	buf := make([]byte, 4)
	offset := 0
	WriteUint32(&buf, id, &offset)
	return buf
}

func BuildSetHost(id uint32, isHost uint8) []byte {
	buf := make([]byte, 20)
	offset := 0
	WriteUint32(&buf, id, &offset)
	WriteUint8(&buf, isHost, &offset)
	return buf[:offset]
}
