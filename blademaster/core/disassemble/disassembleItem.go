package disassemble

import (
	"math/rand"
	"net"
	"time"

	. "github.com/KouKouChan/CSO2-Server/blademaster/core/inventory"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/message"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

var (
	lottoIDs = []uint32{2008, 2013, 2014}
)

func OnDisassembleWeapon(p *PacketData, client net.Conn) {
	//检索数据包
	var pkt InDisassembleWeaponPacket
	if !p.PraseDisassembleWeaponPacket(&pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a error disassemble weapon packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent disassemble weapon packet but not in server !")
		return
	}
	//减少物品
	idx, ok := uPtr.DecreaseItem(pkt.ItemID)
	if !ok {
		OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, CSO2_UI_RAMDOMBOX_ALERT_000)
		DebugInfo(2, "User", uPtr.UserName, "disassemble item", pkt.ItemID, "failed")
		return
	}

	lottoID := getRandomLotto()

	rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeDisassemble),
		buildDisassembleWeapon(lottoID))
	SendPacket(rst, uPtr.CurrentConnection)

	rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeInventory_Create),
		BuildInventoryInfoSingle(uPtr, 0, idx))
	SendPacket(rst, uPtr.CurrentConnection)

	itemIdx := uPtr.AddItem(lottoID, 1, 0)

	rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeInventory_Create),
		BuildInventoryInfoSingle(uPtr, 0, itemIdx))
	SendPacket(rst, uPtr.CurrentConnection)

	DebugInfo(2, "User", uPtr.UserName, "disassembled item", pkt.ItemID, "success")
}

func getRandomLotto() uint32 {
	rand.Seed(time.Now().UnixNano())
	return lottoIDs[rand.Intn(3)]
}

func buildDisassembleWeapon(itemid uint32) []byte {
	buf := make([]byte, 64)
	offset := 0
	WriteUint8(&buf, disassemble, &offset)
	WriteUint8(&buf, disassembleItem, &offset)
	WriteUint8(&buf, 1, &offset)
	WriteUint8(&buf, 1, &offset)
	WriteUint32(&buf, itemid, &offset)
	WriteUint32(&buf, 1, &offset)
	return buf[:offset]
}
