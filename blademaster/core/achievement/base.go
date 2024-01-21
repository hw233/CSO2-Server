package achievement

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

const (
	campaign = 3
	boss     = 4
)

func OnAchievement(p *PacketData, client net.Conn) {
	var pkt InAchievementPacket
	if p.PraseInAchievementPacket(&pkt) {
		switch pkt.Type {
		case campaign:
			OnAchievementCampaign(p, client)
		default:
			DebugInfo(2, "Unknown achievement packet", pkt.Type, "from", client.RemoteAddr().String())
		}
	} else {
		DebugInfo(2, "Error : Recived a illegal achievement packet from", client.RemoteAddr().String())
	}
}

func BuildBackground() []byte { // lobby background feature
	buf := make([]byte, 24)
	offset := 0
	WriteUint8(&buf, 5, &offset) // background
	WriteUint8(&buf, 0, &offset) // call list idk
	WriteUint8(&buf, 1, &offset) // IDK
	WriteUint8(&buf, 1, &offset) // IDK
	return buf[:offset]
}
