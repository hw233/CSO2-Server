package host

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/core/udp"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnHostGameStart(client net.Conn) {
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Host from", client.RemoteAddr().String(), "try to start game server but not in server !")
		return
	}
	//检查用户是不是房主
	rm := GetRoomFromID(uPtr.GetUserChannelServerID(),
		uPtr.GetUserChannelID(),
		uPtr.GetUserRoomID())
	if rm == nil ||
		rm.Id <= 0 {
		DebugInfo(2, "Error : Host", uPtr.UserName, "try to start game server in a null room !")
		return
	}
	//房主开始游戏,设置房间状态
	setting := BuildRoomSetting(rm, 0xFFFFFFFFFFFFFFFF)
	for _, v := range rm.Users {
		if v.Userid != uPtr.Userid {
			rst := BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeRoom), setting)
			SendPacket(rst, v.CurrentConnection)
			if v.IsUserReady() {
				v.ResetAssistNum()
				v.ResetKillNum()
				v.ResetDeadNum()
				v.SetUserIngame(true)
				//给主机发送其他人的数据
				rst := UDPBuild(uPtr.CurrentSequence, 0, v.Userid, v.NetInfo.ExternalIpAddress, v.NetInfo.ExternalClientPort)
				SendPacket(rst, uPtr.CurrentConnection)
				//连接到主机
				rst = UDPBuild(v.CurrentSequence, 1, uPtr.Userid, uPtr.NetInfo.ExternalIpAddress, uPtr.NetInfo.ExternalServerPort)
				SendPacket(rst, v.CurrentConnection)
				//加入主机
				rst = BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeHost), buildJoinHost(uPtr.Userid))
				SendPacket(rst, v.CurrentConnection)
				DebugInfo(2, "Host", v.UserName, "join room", string(rm.Setting.RoomName), "id", rm.Id)
			}
		}
	}
	//给每个人发送房间内所有人的准备状态
	for _, v := range rm.Users {
		temp := buildUserReadyStatus(v)
		for _, k := range rm.Users {
			rst := BytesCombine(BuildHeader(k.CurrentSequence, PacketTypeRoom), temp)
			SendPacket(rst, k.CurrentConnection)
		}
	}
	DebugInfo(2, "Host", uPtr.UserName, "started game server in room", string(rm.Setting.RoomName), "id", rm.Id)
}

func buildJoinHost(id uint32) []byte {
	buf := make([]byte, 13)
	offset := 0
	WriteUint8(&buf, HostJoin, &offset)
	WriteUint32(&buf, id, &offset)
	WriteUint64(&buf, 0, &offset)
	return buf[:offset]
}

func buildUserReadyStatus(u *User) []byte {
	buf := make([]byte, 6)
	offset := 0
	WriteUint8(&buf, OUTSetPlayerReady, &offset)
	WriteUint32(&buf, u.Userid, &offset)
	WriteUint8(&buf, u.Currentstatus, &offset)
	return buf[:offset]
}
