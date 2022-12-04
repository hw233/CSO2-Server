package room

import (
	"log"
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnRoomList(p *PacketData, client net.Conn) {
	var pkt InRoomListRequestPacket
	if p.PraseChannelRequest(&pkt) {
		uPtr := GetUserFromConnection(client)
		if uPtr.Userid <= 0 {
			DebugInfo(2, "Error : A unknow Client from", client.RemoteAddr().String(), "request a RoomList !")
			return
		}

		//发送频道请求返回包
		chlsrv := GetChannelServerWithID(pkt.ChannelServerIndex)
		if chlsrv == nil {
			DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "request a unknown channelServer !")
			return
		}

		//发送频道请求所得房间列表
		chl := GetChannelWithID(pkt.ChannelIndex, chlsrv)
		if chl == nil {
			log.Println("Error : Client from", client.RemoteAddr().String(), "request a unknown channel !")
			return
		}

		rst := BuildLobbyReply(uPtr.CurrentSequence, *p, chlsrv.ServerIndex, chl.ChannelID)
		SendPacket(rst, uPtr.CurrentConnection)
		DebugInfo(2, "Sent a lobbyReply packet to", client.RemoteAddr().String())

		rst = BuildRoomList(uPtr.CurrentSequence, chl)
		SendPacket(rst, uPtr.CurrentConnection)
		DebugInfo(2, "Sent a roomList packet to", client.RemoteAddr().String())

		//设置用户所在频道
		uPtr.SetUserChannelServer(chlsrv.ServerIndex)
		uPtr.SetUserChannel(chl.ChannelID)
	} else {
		log.Println("Recived a damaged packet from", client.RemoteAddr().String())
	}
}

func OnBroadcastRoomList(chlsrvid uint8, chlid uint8, u *User) {
	if u == nil {
		return
	}
	//发送频道请求返回包
	chlsrv := GetChannelServerWithID(chlsrvid)
	if chlsrv == nil {
		DebugInfo(2, "Error : Client from", u.CurrentConnection.RemoteAddr().String(), "request a unknown channelServer !")
		return
	}
	//发送频道请求所得房间列表
	chl := GetChannelWithID(chlid, chlsrv)
	if chl == nil {
		log.Println("Error : Client from", u.CurrentConnection.RemoteAddr().String(), "request a unknown channel !")
		return
	}
	//检索房间是否正常
	chl.ChannelMutex.Lock()
	for k, v := range chl.Rooms {
		num, orgnum := 0, 0

		for i, user := range v.Users {
			if user.GetUserRoomID() != v.Id {
				delete(v.Users, i)
				continue
			}
			if user != nil {
				num++
			}
		}
		orgnum = int(v.NumPlayers)

		if num <= 0 {
			DelChannelRoomQuick(k, chl)
		} else if num != orgnum {
			v.SyncUserNum(uint8(num))
		}
	}
	chl.ChannelMutex.Unlock()
	//发送房间数据
	rst := BuildRoomList(u.CurrentSequence, chl)
	SendPacket(rst, u.CurrentConnection)
	DebugInfo(2, "Sent a roomList packet to", u.CurrentConnection.RemoteAddr().String())
}

func BuildLobbyReply(seq *uint8, p PacketData, chlsrvid, chlid uint8) []byte {
	rst := BuildHeader(seq, PacketTypeLobby)
	lob := OutLobbyJoinRoom{
		JoinLobby, chlsrvid, chlid,
	}
	rst = append(rst,
		SetLobby,
		lob.Unk00,
		lob.Unk01,
		lob.Unk02)
	return rst
}

func BuildCurrentLobby(seq *uint8, chlsrvid, chlid uint8) []byte {
	rst := BuildHeader(seq, PacketTypeLobby)
	buf := make([]byte, 1)
	offset := 0
	WriteUint8(&buf, 0, &offset)
	rst = BytesCombine(rst, buf[:offset], UsersManager.GetChannelUsers(chlsrvid, chlid))
	return rst
}

//暂定
func BuildRoomList(seq *uint8, chl *ChannelInfo) []byte {
	rst := BuildHeader(seq, PacketTypeRoomList)
	rst = append(rst,
		SendFullRoomList,
	)
	buf := make([]byte, 2)
	tempoffset := 0
	chl.ChannelMutex.Lock()
	WriteUint16(&buf, chl.RoomNum, &tempoffset)
	for _, v := range chl.Rooms {
		if v == nil {
			DebugInfo(1, "Waring! here is a null room in channelID", chl.ChannelID)
			continue
		}
		roombuf := make([]byte, 512)
		offset := 0
		WriteUint16(&roombuf, v.Id, &offset)
		WriteUint64(&roombuf, 0xFFFFFFFFFFFFFFFF, &offset)
		WriteString(&roombuf, v.Setting.RoomName, &offset)
		WriteUint8(&roombuf, v.RoomNumber, &offset)
		WriteUint8(&roombuf, v.PasswordProtected, &offset)
		WriteUint16(&roombuf, 0, &offset)
		WriteUint8(&roombuf, v.Setting.GameModeID, &offset)
		WriteUint8(&roombuf, v.Setting.MapID, &offset)
		WriteUint8(&roombuf, v.NumPlayers, &offset)
		WriteUint8(&roombuf, v.Setting.MaxPlayers, &offset)
		WriteUint8(&roombuf, v.Unk08, &offset)
		WriteUint32(&roombuf, v.HostUserID, &offset)
		WriteString(&roombuf, v.HostUserName, &offset)
		WriteUint8(&roombuf, v.Unk11, &offset)
		WriteUint8(&roombuf, v.Unk12, &offset)
		WriteUint32(&roombuf, v.Unk13, &offset)
		WriteUint16(&roombuf, v.Unk14, &offset)
		WriteUint16(&roombuf, v.Unk15, &offset)
		WriteUint32(&roombuf, v.Unk16, &offset)
		WriteUint16(&roombuf, v.Unk17, &offset)
		WriteUint16(&roombuf, v.Unk18, &offset)
		WriteUint8(&roombuf, v.Unk19, &offset)
		WriteUint8(&roombuf, v.Unk20, &offset)
		if v.Unk20 == 1 {
			WriteUint32(&roombuf, 0, &offset)
			WriteUint8(&roombuf, 0, &offset)
			WriteUint32(&roombuf, 0, &offset)
			WriteUint8(&roombuf, 0, &offset)
		}
		WriteUint8(&roombuf, v.Unk21, &offset)
		WriteUint8(&roombuf, v.Setting.Status, &offset)
		WriteUint8(&roombuf, v.Setting.AreBotsEnabled, &offset)
		WriteUint8(&roombuf, v.Unk24, &offset)
		WriteUint16(&roombuf, v.Setting.StartMoney, &offset)
		WriteUint8(&roombuf, v.Unk26, &offset)
		WriteUint8(&roombuf, 0, &offset)
		WriteUint8(&roombuf, v.Unk28, &offset)
		WriteUint8(&roombuf, v.Unk29, &offset)
		WriteUint64(&roombuf, v.Unk30, &offset)
		WriteUint8(&roombuf, v.Setting.WinLimit, &offset)
		WriteUint16(&roombuf, v.Setting.KillLimit, &offset)
		WriteUint8(&roombuf, v.Setting.ForceCamera, &offset)
		// WriteUint8(&roombuf, v.botEnabled, &offset)
		// if v.botEnabled == 1 {
		// 	WriteUint8(&roombuf, v.botDifficulty, &offset)
		// 	WriteUint8(&roombuf, v.numCtBots, &offset)
		// 	WriteUint8(&roombuf, v.numTrBots, &offset)
		// }
		WriteUint8(&roombuf, v.Unk31, &offset)
		WriteUint8(&roombuf, v.Unk35, &offset)
		WriteUint8(&roombuf, v.Setting.NextMapEnabled, &offset)
		WriteUint8(&roombuf, v.Setting.ChangeTeams, &offset)
		WriteUint8(&roombuf, v.AreFlashesDisabled, &offset)
		WriteUint8(&roombuf, v.CanSpec, &offset)
		WriteUint8(&roombuf, v.IsVipRoom, &offset)
		WriteUint8(&roombuf, v.VipRoomLevel, &offset)
		WriteUint8(&roombuf, v.Setting.Difficulty, &offset)
		buf = BytesCombine(buf, roombuf[:offset])
	}
	chl.ChannelMutex.Unlock()
	return BytesCombine(rst, buf)
}
