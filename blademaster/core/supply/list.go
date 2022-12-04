package supply

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

var (
	supplyboxlist []byte
)

func OnSupplyList(p *PacketData, client net.Conn) {
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to request supply list but not in server !")
		return
	}
	//发送数据
	rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeSupply), BuildSupplyList())
	SendPacket(rst, uPtr.CurrentConnection)
	DebugInfo(2, "Send a supply box list to User", uPtr.UserName)

}
func BuildSupplyList() []byte {
	tmp := make([]byte, 5)
	offset := 0
	WriteUint8(&tmp, requestList, &offset)
	WriteUint8(&tmp, 1, &offset)
	return BytesCombine(tmp[:offset], supplyboxlist)
}

func InitBoxReply() {
	buf := make([]byte, 1)
	offset := 0
	WriteUint8(&buf, uint8(len(BoxList)), &offset)
	for _, v := range BoxList {
		tmp := make([]byte, len(v.Items)*16+20)
		offset = 0
		WriteUint32(&tmp, v.BoxID, &offset)
		WriteUint32(&tmp, v.OptId, &offset) //nextOptIndex
		WriteUint8(&tmp, uint8(len(v.Items)), &offset)
		for _, item := range v.Items {
			WriteUint32(&tmp, item.ItemID, &offset)
			WriteUint16(&tmp, 0, &offset)
			WriteUint16(&tmp, item.Day, &offset)
			WriteUint64(&tmp, 0, &offset)
		}
		if v.KeyId == 0 {
			WriteUint8(&tmp, 0, &offset)
		} else {
			WriteUint8(&tmp, 1, &offset)
			WriteUint32(&tmp, v.KeyId, &offset)
		}
		buf = BytesCombine(buf, tmp[:offset])
	}
	supplyboxlist = buf
}
