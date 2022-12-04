package room

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/core/udp"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnGameStart(p *PacketData, client net.Conn) {
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to start game but not in server !")
		return
	}
	//检查用户是不是房主
	rm := GetRoomFromID(uPtr.GetUserChannelServerID(),
		uPtr.GetUserChannelID(),
		uPtr.GetUserRoomID())
	if rm == nil ||
		rm.Id <= 0 {
		DebugInfo(2, "Error : User", uPtr.UserName, "try to start game in a null room !")
		return
	}
	//房主开始游戏,设置房间状态
	if rm.HostUserID == uPtr.Userid {
		rm.StopCountdown()
		rm.SetStatus(StatusIngame)
		rm.ResetRoomKillNum()
		rm.ResetRoomScore()
		rm.ResetRoomWinner()
		//设置用户状态
		uPtr.SetUserIngame(true)
		uPtr.ResetKillNum()
		uPtr.ResetDeadNum()
		uPtr.ResetAssistNum()
		//主机开始游戏
		rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeHost), BuildGameStart(uPtr.Userid))

		SendPacket(rst, uPtr.CurrentConnection)
		DebugInfo(2, "Host", uPtr.UserName, "started game in room", string(rm.Setting.RoomName))

	} else if rm.Setting.IsIngame != 0 {
		host := rm.RoomGetUser(rm.HostUserID)
		if host == nil ||
			host.Userid <= 0 {
			DebugInfo(2, "Error : User", uPtr.UserName, "try to start game but host is null !")
			return
		}
		//设置用户状态
		uPtr.ResetKillNum()
		uPtr.ResetDeadNum()
		uPtr.ResetAssistNum()
		uPtr.SetUserIngame(true)
		//发送房间数据
		rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeRoom), BuildRoomSetting(rm, 0xFFFFFFFFFFFFFFFF))
		SendPacket(rst, uPtr.CurrentConnection)
		//给主机发送其他人的数据
		rst = UDPBuild(host.CurrentSequence, 0, uPtr.Userid, uPtr.NetInfo.ExternalIpAddress, uPtr.NetInfo.ExternalClientPort)
		SendPacket(rst, host.CurrentConnection)
		//连接到主机
		rst = UDPBuild(uPtr.CurrentSequence, 1, host.Userid, host.NetInfo.ExternalIpAddress, host.NetInfo.ExternalServerPort)
		SendPacket(rst, uPtr.CurrentConnection)
		//加入主机
		rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeHost), BuildJoinHost(host.Userid))
		SendPacket(rst, uPtr.CurrentConnection)
		//给每个人发送房间内所有人的准备状态
		for _, v := range rm.Users {
			temp := BuildUserReadyStatus(v)
			for _, k := range rm.Users {
				rst = BytesCombine(BuildHeader(k.CurrentSequence, PacketTypeRoom), temp)
				SendPacket(rst, k.CurrentConnection)
			}
		}
		DebugInfo(2, "User", uPtr.UserName, "joined in game in room", string(rm.Setting.RoomName), "id", rm.Id)
	}
}

func BuildJoinHost(id uint32) []byte {
	buf := make([]byte, 13)
	offset := 0
	WriteUint8(&buf, HostJoin, &offset)
	WriteUint32(&buf, id, &offset)
	WriteUint64(&buf, 0, &offset)
	return buf[:offset]
}

func BuildGameStart(id uint32) []byte {
	buf := make([]byte, 5)
	offset := 0
	WriteUint8(&buf, GameStart, &offset)
	WriteUint32(&buf, id, &offset)
	return buf[:offset]
}
