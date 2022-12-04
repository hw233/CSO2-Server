package automatch

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnAutoMatch(p *PacketData, client net.Conn) {
	DebugInfo(2, "Unknown AutoMatch packet from", client.RemoteAddr().String())
}
