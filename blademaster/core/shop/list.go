package shop

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/configure"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

var (
	ShopReply []byte
)

func OnShopList(p *PacketData, client net.Conn) {
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to request shop list but not in server !")
		return
	}
	//发送数据
	rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeShop), BuildShopList())
	SendPacket(rst, uPtr.CurrentConnection)
	DebugInfo(2, "Send shop list to User", uPtr.UserName)

}
func BuildShopList() []byte {
	if Conf.EnableShop == 0 {
		return []byte{0, 0, 0}
	}
	return ShopReply
}

func InitShopReply() {
	buf := make([]byte, 3)
	offset, optIdx := 0, 0
	WriteUint8(&buf, outshoplist, &offset)
	WriteUint16(&buf, uint16(len(ShopItemList)), &offset)
	for k, _ := range ShopItemList {
		optIdx = k * 5
		tmp := make([]byte, 512)
		offset = 0

		saveBoxOptId(ShopItemList[k].ItemID, uint32(optIdx))

		WriteUint32(&tmp, ShopItemList[k].ItemID, &offset)
		WriteUint8(&tmp, ShopItemList[k].Currency, &offset)
		WriteUint8(&tmp, uint8(len(ShopItemList[k].Opt)), &offset) //numopt
		for _, opt := range ShopItemList[k].Opt {
			WriteUint32(&tmp, uint32(optIdx), &offset) //optidx
			WriteUint16(&tmp, opt.Day, &offset)        //quantity
			WriteUint64(&tmp, 0, &offset)              //continue~day
			WriteUint8(&tmp, 0, &offset)
			WriteUint16(&tmp, opt.Count, &offset)
			WriteUint32(&tmp, opt.Price, &offset)
			WriteUint32(&tmp, opt.Price, &offset)
			WriteUint8(&tmp, 0, &offset) //discount
			WriteUint32(&tmp, 0, &offset)
			WriteUint32(&tmp, 0, &offset)
			WriteUint8(&tmp, 0, &offset) //flags
			WriteUint8(&tmp, 0, &offset)
			WriteUint8(&tmp, 1, &offset)
			WriteUint8(&tmp, 0, &offset)
			WriteUint32(&tmp, 0, &offset)
			WriteUint8(&tmp, 0, &offset) //bmacket
			WriteUint8(&tmp, 0, &offset) //efficiency
			optIdx++
			if optIdx/5 == 0 { //避免溢出
				break
			}
		}

		buf = BytesCombine(buf, tmp[:offset])
	}
	ShopReply = buf
}

func saveBoxOptId(boxid, optid uint32) {
	for k, v := range BoxList {
		if BoxList[k].BoxID == boxid {
			v.OptId = optid
			BoxList[k] = v
			break
		}
	}
}
