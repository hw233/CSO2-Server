package disassemble

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

const (
	disassemble = 1
)

func OnDisassemble(p *PacketData, client net.Conn) {
	var pkt InDisassemblePacket
	if p.PraseDisassemblePacket(&pkt) {
		switch pkt.Type {
		case disassemble:
			OnDisassembleItem(p, client)
		default:
			DebugInfo(2, "Unknown disassemble packet", pkt.Type, "from", client.RemoteAddr().String())
		}
	} else {
		DebugInfo(2, "Error : Recived a illegal disassemble packet from", client.RemoteAddr().String())
	}
}
