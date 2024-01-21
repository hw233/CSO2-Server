package achievement

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/core/message"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

const (
	unlock = 1
	lock   = 2
)

var num uint8 = 4 // packet testing.
var value uint8 = 0// packet testing.

func OnAchievementBackground(p *PacketData, client net.Conn) {
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent Achievement packet but not in server !")
		return
	}
	SendPacket(OnBuildBackground(num, uPtr), uPtr.CurrentConnection)
	num++
	//rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeAchievement), BuildBackground())
	//SendPacket(rst, uPtr.CurrentConnection)

	var pkt InAchievementBackgroundPacket
	if !p.PraseInAchievementBackgroundPacket(&pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a error Achievement background packet !")
		return
	}

	switch pkt.ACHType {
	case unlock:
		uPtr.BGType = pkt.SelectedCase

    // Unfortunately, I could not apply the selected background and could not find it.
		rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeAchievement), SetBackground(pkt.SelectedCase, num))
		SendPacket(rst, uPtr.CurrentConnection)

		rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeLobby), SetBackground(pkt.SelectedCase, num))
		SendPacket(rst, uPtr.CurrentConnection)


		DebugInfo(2, "num", num) //I used it for packet search
		num++

		//OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, GameUI_InvalidID)
		DebugInfo(2, "Achievement background unlock packet", pkt.ACHType, "from", client.RemoteAddr().String())
	case lock:
		OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, GameUI_InvalidID)
		DebugInfo(2, "Achievement background lock packet", pkt.ACHType, "from", client.RemoteAddr().String())
	default:
		DebugInfo(2, "Unknown switch case Achievement background packet", pkt.ACHType, "from", client.RemoteAddr().String())
	}

	DebugInfo(2, "SelectedImage", pkt.SelectedCase)
	DebugInfo(2, "Type", pkt.Type)
	DebugInfo(2, "Type2", pkt.Type2)
	DebugInfo(2, "Type3", pkt.Type3)

	DebugInfo(2, "Achievement background packet", pkt.ACHType, "from", client.RemoteAddr().String())
}

func SetBackground(BGType uint8, num uint8) []byte {  // not work...
	buf := make([]byte, 90)
	offset := 0
	WriteUint8(&buf, 5, &offset) // background
	WriteUint8(&buf, 1, &offset) // active
	WriteUint8(&buf, 1, &offset)
	WriteUint8(&buf, BGType, &offset)
	WriteUint8(&buf, BGType, &offset)
	WriteUint8(&buf, BGType, &offset)
	WriteUint8(&buf, BGType, &offset)
	WriteUint8(&buf, BGType, &offset)
	WriteUint8(&buf, BGType, &offset)
	WriteUint8(&buf, BGType, &offset)
	WriteUint8(&buf, BGType, &offset)
	WriteUint8(&buf, BGType, &offset)
	WriteUint8(&buf, BGType, &offset)
	WriteUint8(&buf, BGType, &offset)
	WriteUint8(&buf, BGType, &offset)
	return buf[:offset]
}
