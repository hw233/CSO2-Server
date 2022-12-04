package host

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/core/inventory"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnHostItemUsing(p *PacketData, client net.Conn) {
	//检查数据包
	var pkt InHostItemUsingPacket
	if !p.PraseInHostItemUsingPacket(&pkt) {
		DebugInfo(2, "Error : Cannot prase a ItemUsing packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : A host send ItemUsing but not in server!")
		return
	}
	dest := GetUserFromID(pkt.UserID)
	if dest == nil ||
		dest.Userid <= 0 {
		DebugInfo(2, "Error : A host send ItemUsing but dest user is null!")
		return
	}
	//找到玩家的房间
	rm := GetRoomFromID(uPtr.GetUserChannelServerID(),
		uPtr.GetUserChannelID(),
		uPtr.GetUserRoomID())
	if rm == nil ||
		rm.Id <= 0 {
		DebugInfo(2, "Error : User", uPtr.UserName, "try to send ItemUsing but in a null room !")
		return
	}
	//发送用户背包数据
	itemIdx, ok := dest.DecreaseItem(pkt.ItemID)
	if !ok {
		DebugInfo(2, "User", uPtr.UserName, "use item", pkt.ItemID, "in match failed")
		return
	}
	rst := BytesCombine(BuildHeader(dest.CurrentSequence, PacketTypeInventory_Create),
		BuildInventoryInfoSingle(dest, 0, itemIdx))
	SendPacket(rst, dest.CurrentConnection)
	//发送房主数据包
	rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeHost), BuildItemUsing(pkt.UserID, pkt.ItemID, dest.GetItemCount(itemIdx)))
	SendPacket(rst, uPtr.CurrentConnection)
	DebugInfo(2, "Send User", dest.UserName, "ItemUsed packet to host", uPtr.UserName)
}

func BuildItemUsing(uid uint32, itemid uint32, num int) []byte {
	buf := make([]byte, 20)
	offset := 0
	WriteUint8(&buf, ItemUsing, &offset)
	WriteUint32(&buf, uid, &offset)
	WriteUint32(&buf, itemid, &offset)
	if itemid == 2019 {
		WriteUint8(&buf, 1, &offset)
		WriteUint8(&buf, 1, &offset)
		WriteUint8(&buf, 1, &offset)
		WriteUint8(&buf, 1, &offset)
	} else {
		WriteUint8(&buf, uint8(num), &offset)
	}
	return buf[:offset]
}
