package useitem

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/core/message"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnTryItemUse(p *PacketData, client net.Conn) {
	//检索数据包
	var pkt InTryItemUsePacket
	if !p.PraseTryItemUsePacket(&pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a error itemuse packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to request use item but not in server !")
		return
	}
	//发送数据
	itemID := uPtr.GetItemIDBySeq(pkt.ItemSeq)
	switch itemID {
	case 2001: //改名卡
		if IsExistsIngameName(pkt.NewName) {
			OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, CSO2_POPUP_NICKNAME_ALREADY_EXIST)
			DebugInfo(2, "User", uPtr.UserName, "try change nickname to", string(pkt.NewName), "but this name already exists")
			return
		}
		rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeUseItem),
			onBuildChangeName())
		SendPacket(rst, uPtr.CurrentConnection)
		DebugInfo(2, "User", uPtr.UserName, "try change nickname to", string(pkt.NewName))
	default:
		DebugInfo(2, "User", uPtr.UserName, "try using item but itemid is", itemID)
		return
	}

}

func onBuildChangeName() []byte {
	buf := make([]byte, 16)
	offset := 0
	WriteUint8(&buf, tryuseitem, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 1, &offset)
	return buf[:offset]
}
