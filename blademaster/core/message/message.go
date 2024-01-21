package message

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
)
var num uint8 = 2 // packet test (I'm search with numbers)
func OnSendMessageMegaphone(seq *uint8, client net.Conn, tp uint8, itemid uint32, playername []byte, message []byte) {
	rst := BuildHeader(seq, PacketTypeChat)
	rst = append(rst, tp)
	rst = BytesCombine(rst, BuildMessageeMegaphone(tp, itemid, playername, message))
	SendPacket(rst, client)
}


func BuildMessageeMegaphone(tp uint8, itemid uint32, playername []byte, message []byte) []byte {
	if tp == MessageAnnouncement {
		buf := make([]byte, 512)
		offset := 0
		chattype := uint8(1)
		switch itemid {
		case 2010:
			chattype = 1 // Global Message
			fmt.Println("1")
		case 2011:
			chattype = 2 // Server
			fmt.Println("2")
		case 2012:
			chattype = 3
			fmt.Println("3") // Channel
		case 2015:
			chattype = 4
			fmt.Println("4") // All Servers
		default:
			fmt.Println("Invalid value")
		}
		WriteUint8(&buf, 1, &offset)
		WriteUint8(&buf, chattype, &offset) // Server Chat Type
		WriteUint8(&buf, 0, &offset)
		WriteUint8(&buf, 1, &offset)
		WriteString(&buf, playername, &offset) // Oyuncu İsmi
		WriteString(&buf, message, &offset)    // Mesaj
		WriteUint8(&buf, 1, &offset)
		num++
		return buf[:offset]
	}
	return BuildLongString(playername)
}


func OnSendMessage(seq *uint8, client net.Conn, tp uint8, msg []byte) {
	rst := BuildHeader(seq, PacketTypeChat)
	rst = append(rst, tp)
	rst = BytesCombine(rst, BuildMessage(msg, tp))
	SendPacket(rst, client)
}

func BuildMessage(msg []byte, tp uint8) []byte {
	if tp == MessageCongratulate {
		buf := make([]byte, 1)
		buf[0] = 0
		return BytesCombine(buf, BuildString(msg))
	}
	return BuildLongString(msg)
}
