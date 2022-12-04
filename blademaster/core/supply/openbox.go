package supply

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

func OnSupplyOpenBox(p *PacketData, client net.Conn) {
	//检索数据包
	var pkt InOpenBoxPacket
	if !p.PraseOpenBoxPacket(&pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a error OpenBox packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to request openbox but not in server !")
		return
	}
	//检查军刀是否够
	if !checkBoxKnife(pkt.BoxID, uPtr) {

		OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, CSO2_UI_RANDOMBOX_KEY_ALERT_BUY_01)
		DebugInfo(2, "User", uPtr.UserName, "open box", BoxList[pkt.BoxID].BoxName, "failed with not enough knife")
		return
	}
	//发送数据
	if !uPtr.UseBox(pkt.BoxID, BoxList[pkt.BoxID].KeyId) {

		OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, CSO2_UI_RAMDOMBOX_ALERT_000)
		DebugInfo(2, "User", uPtr.UserName, "open box", BoxList[pkt.BoxID].BoxName, "key id", BoxList[pkt.BoxID].KeyId, "failed")
		return
	}

	rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeInventory_Create),
		BuildInventoryInfoSingle(uPtr, pkt.BoxID, -1))
	SendPacket(rst, uPtr.CurrentConnection)

	rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeInventory_Create),
		BuildInventoryInfoSingle(uPtr, BoxList[pkt.BoxID].KeyId, -1))
	SendPacket(rst, uPtr.CurrentConnection)

	itemid, day := GetBoxItem(pkt.BoxID)
	idx := uPtr.AddItem(itemid, 1, uint64(day))

	rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeSupply),
		BuildSupplyOpenBox(itemid, 0, day))
	SendPacket(rst, uPtr.CurrentConnection)

	rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeInventory_Create),
		BuildInventoryInfoSingle(uPtr, 0, idx))
	SendPacket(rst, uPtr.CurrentConnection)

	DebugInfo(2, "User", uPtr.UserName, "got item", ItemList[itemid].Name, "id", itemid, "by openning box", BoxList[pkt.BoxID].BoxName)

}

func BuildSupplyOpenBox(itemid uint32, count, day uint16) []byte {
	buf := make([]byte, 25)
	offset := 0
	WriteUint8(&buf, openbox, &offset)
	WriteUint8(&buf, 1, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint32(&buf, itemid, &offset) //itemid
	WriteUint16(&buf, count, &offset)  //item count
	WriteUint16(&buf, day, &offset)    //day
	WriteUint32(&buf, 0, &offset)      //mpoint
	return buf[:offset]
}

func GetBoxItem(boxid uint32) (uint32, uint16) {
	if v, ok := BoxList[boxid]; ok {
		rand.Seed(time.Now().UnixNano())
		radV := rand.Intn(v.TotalValue)
		for _, item := range v.Items {
			radV -= item.Value
			if radV < 0 {
				return item.ItemID, item.Day
			}
		}
	}
	DebugInfo(2, "Error : can't open box", boxid)
	return 0, 0
}

func checkBoxKnife(itemid uint32, u *User) bool {
	_, ok := BoxList[itemid]
	if u == nil || !ok {
		return false
	}
	if ok && BoxList[itemid].KeyId == 0 {
		return true
	}
	u.UserMutex.Lock()
	defer u.UserMutex.Unlock()
	for _, v := range u.Inventory.Items {
		if v.Id == BoxList[itemid].KeyId && v.Count > 0 {
			return true
		}
	}
	return false
}
