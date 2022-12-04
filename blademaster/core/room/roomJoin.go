package room

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/core/message"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnJoinRoom(p *PacketData, client net.Conn) {
	//检索数据包
	var pkt InJoinRoomPacket
	if !p.PraseJoinRoomPacket(&pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a error JoinRoom packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to join room but not in server !")
		return
	}
	//检索玩家房间
	if uPtr.CurrentRoomId != 0 {
		DebugInfo(2, "Error : User", uPtr.UserName, "try to join room but in another room !")
		OnLeaveRoom(uPtr.CurrentConnection, false)
		return
	}
	//找到对应房间
	rm := GetRoomFromID(uPtr.GetUserChannelServerID(),
		uPtr.GetUserChannelID(),
		uint16(pkt.RoomId))
	if rm == nil ||
		rm.Id <= 0 {
		OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox,
			GAME_ROOM_JOIN_FAILED_CLOSED)
		DebugInfo(2, "Error : User", uPtr.UserName, "try to join a null room ! ID=", pkt.RoomId)
		return
	}
	//检索密码
	if rm.Setting.LenOfPassWd > 0 {
		if !CompareBytes(pkt.PassWord, rm.Setting.PassWd) {
			OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox,
				GAME_ROOM_JOIN_FAILED_BAD_PASSWORD)
			DebugInfo(2, "User", uPtr.UserName, "try to join a room with error password!")
			return
		}
	}
	//检索房间状态
	if rm.GetFreeSlots() <= 0 {
		OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox,
			GAME_ROOM_JOIN_FAILED_FULL)
		DebugInfo(2, "User", uPtr.UserName, "try to join a full room ! ID=", pkt.RoomId)
		return
	}
	//玩家加进房间
	if !rm.JoinUser(uPtr) {
		OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox,
			GAME_ROOM_JOIN_ERROR)
		DebugInfo(2, "User", uPtr.UserName, "join  room ID=", rm.Id, "failed!")
		return
	}
	//发送数据
	rst := append(BuildHeader(uPtr.CurrentSequence, PacketTypeRoom), OUTCreateAndJoin)
	rst = BytesCombine(rst, BuildCreateAndJoin(rm))
	SendPacket(rst, uPtr.CurrentConnection)
	DebugInfo(2, "User", uPtr.UserName, "joined room", string(rm.Setting.RoomName), "id", rm.Id)
	rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeRoom), BuildRoomSetting(rm, 0xFFFFFFFFFFFFFFFF))
	SendPacket(rst, client)
	DebugInfo(2, "Sent a room setting packet to", uPtr.UserName)
	//发送玩家状态
	ustatus := BuildUserReadyStatus(uPtr)
	uplayjoin := BuildPlayerJoin(uPtr)

	for _, v := range rm.Users {
		rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeRoom), BuildUserReadyStatus(v))
		SendPacket(rst, uPtr.CurrentConnection)
		if v.Userid != uPtr.Userid {
			//发送给其他玩家该玩家信息
			rst = BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeRoom), uplayjoin)
			SendPacket(rst, v.CurrentConnection)
			rst = BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeRoom), ustatus)
			SendPacket(rst, v.CurrentConnection)
		}
	}
	DebugInfo(2, "Sync user status to all player in room", string(rm.Setting.RoomName), "id", rm.Id)
	return
}

func BuildPlayerJoin(u *User) []byte {
	buf := make([]byte, 5)
	offset := 0
	WriteUint8(&buf, OUTPlayerJoin, &offset)
	WriteUint32(&buf, u.Userid, &offset)
	buf = BytesCombine(buf, u.BuildUserNetInfo(),
		BuildUserInfo(0xFFFFFFFF, NewUserInfo(u), 0, false))
	return buf
}
