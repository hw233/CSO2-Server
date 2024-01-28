package quick

import (
	"net"

	//. "github.com/KouKouChan/CSO2-Server/blademaster/core/room"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)


var roomgamemodeid uint16
var botistatus uint8


func OnQuickList(p *PacketData, client net.Conn) {
	//检索数据包
	var pkt InQuickList
	if !p.PraseInQuickListPacket(&pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a error QuickList packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to request QuickList but not in server !")
		return
	}

	botistatus = pkt.IsEnableBot          // User Bot Preference is based on the logic that if 0 and 1 are 0, there will be rooms without bots, and if 1, there will be rooms with bots.
	roomgamemodeid = uint16(pkt.GameModID) // Mode ID Selected in the List
	

	rst2 := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeQuickList), BuildQuickListTest(pkt))
	SendPacket(rst2, uPtr.CurrentConnection)
	DebugInfo(2, "Sent a QuickList response to user", uPtr.UserName)
}


func BuildQuickListTest(pkt InQuickList) []byte {
	buf := make([]byte, 512)
	offset := 0
	WriteUint8(&buf, QuickList, &offset)

	// Filter rooms via Room Manager
	var filteredRooms []*Room
	//RoomsManager.Lock.Lock()
	for _, room := range RoomsManager.Rooms {
		if room.Setting.GameModeID == pkt.GameModID && room.Setting.AreBotsEnabled == pkt.IsEnableBot {
			filteredRooms = append(filteredRooms, room)
		}
	}
	//RoomsManager.Lock.Unlock()

	WriteUint8(&buf, uint8(len(filteredRooms)), &offset) // Filtrelenmiş oda sayısı

	for _, room := range filteredRooms {
		WriteUint8(&buf, uint8(room.Id), &offset) // RoomID
		WriteUint8(&buf, botistatus, &offset)
		WriteUint8(&buf, botistatus, &offset)              // RoomID
		WriteUint8(&buf, room.Setting.GameModeID, &offset) // GameModeID
		WriteString(&buf, room.Setting.RoomName, &offset)  // RoomName
		WriteUint8(&buf, room.Setting.MapID, &offset)      // MapID
		WriteUint8(&buf, botistatus, &offset)
		WriteUint8(&buf, room.Setting.MaxPlayers, &offset) // Max Players
		WriteUint8(&buf, room.NumPlayers, &offset)         // Min Players
		WriteUint8(&buf, botistatus, &offset)
		WriteUint8(&buf, room.Setting.IsIngame, &offset)
	}
	return buf[:offset]
}


func BuildQuickListOld(pkt InQuickList) []byte {
	buf := make([]byte, 2)
	offset := 0
	WriteUint8(&buf, QuickList, &offset)
	WriteUint8(&buf, 0, &offset) //num of room

	return buf[:offset]
}
