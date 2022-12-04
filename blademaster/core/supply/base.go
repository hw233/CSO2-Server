package supply

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

const (
	requestList = 0
	openbox     = 1
)

func OnSupplyRequest(p *PacketData, client net.Conn) {
	var pkt InSupplyPacket
	if p.PraseSupplyPacket(&pkt) {
		switch pkt.Type {
		case requestList:
			OnSupplyList(p, client)
		case openbox:
			OnSupplyOpenBox(p, client)
		default:
			DebugInfo(2, "Unknown supply packet", pkt.Type, "from", client.RemoteAddr().String())
		}
	} else {
		DebugInfo(2, "Error : Recived a illegal shop packet from", client.RemoteAddr().String())
	}
}
