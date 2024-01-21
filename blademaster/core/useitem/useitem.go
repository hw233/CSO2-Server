package useitem

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/KouKouChan/CSO2-Server/blademaster/core/inventory"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/inventory"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/message"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

const (
	lotto_base       = 1 //银币
	lotto_max        = 6000
	lotto_event_base = 1 //铜币
	lotto_event_max  = 8900
	lotto_gold_base  = 1 //金币
	lotto_gold_max   = 12000
)


// Map to store users currently using the megaphone
var megaphoneUsers = make(map[uint32]bool)

// Queue to keep track of users waiting to use the megaphone
var megaphoneQueue []uint32

func OnItemUse(p *PacketData, client net.Conn) {
	//检索数据包
	var pkt InItemUsePacket
	if !p.PraseItemUsePacket(&pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a error pointlottouse packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to request use pointlotto but not in server !")
		return
	}
	//发送数据
	itemID := uPtr.GetItemIDBySeq(pkt.ItemSeq)
	switch itemID {
	case 2001: //改名卡
		if IsExistsIngameName(pkt.String) {
			OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, CSO2_POPUP_NICKNAME_ALREADY_EXIST)
			DebugInfo(2, "User", uPtr.UserName, "try change nickname to", string(pkt.String), "but this name already exists")
			return
		}

		if err := DelOldNickNameFile(uPtr.IngameName); err != nil {
			OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, CSO2_UI_RAMDOMBOX_ALERT_000)
			DebugInfo(2, "User", uPtr.UserName, "try change nickname to", string(pkt.String), "failed", err)
			return
		}

		idx, ok := uPtr.DecreaseItem(itemID)
		if !ok {
			OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, CSO2_UI_RAMDOMBOX_ALERT_000)
			DebugInfo(2, "User", uPtr.UserName, "use item", itemID, "failed")
			return
		}

		rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeInventory_Create),
			BuildInventoryInfoSingle(uPtr, 0, idx))
		SendPacket(rst, uPtr.CurrentConnection)

		uPtr.SetUserName(uPtr.UserName, string(pkt.String))

		rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeUseItem),
			buildNickNameChange())
		SendPacket(rst, uPtr.CurrentConnection)

		DebugInfo(2, "User", uPtr.UserName, "changed nickname to", string(pkt.String))
	case 2008, 2013, 2014: //银币 id 2008，铜币，金币
		idx, ok := uPtr.DecreaseItem(itemID)
		if !ok {
			OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, CSO2_UI_RAMDOMBOX_ALERT_000)
			DebugInfo(2, "User", uPtr.UserName, "use item", itemID, "failed")
			return
		}

		rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeInventory_Create),
			BuildInventoryInfoSingle(uPtr, 0, idx))
		SendPacket(rst, uPtr.CurrentConnection)

		point := UsePointLotto(itemID)
		uPtr.GetPoints(point)

		rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeUseItem),
			buildUsePoint(uint32(point)))
		SendPacket(rst, uPtr.CurrentConnection)

		DebugInfo(2, "User", uPtr.UserName, "got point", point, "by using pointlotto", itemID)
	/*case 2010: //频道喇叭
		//查找玩家当前频道
		chlsrv := GetChannelServerWithID(uPtr.GetUserChannelServerID())
		if chlsrv == nil || chlsrv.ServerIndex <= 0 {
			OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, CSO2_POPUP_MEGAPHONE_USE_FAIL_INVALID_ITEM)
			DebugInfo(2, "User", uPtr.UserName, "use item", itemID, "failed,not in channel server")
			return
		}
		chl := GetChannelServerWithID(uPtr.GetUserChannelID())
		if chl == nil || chl.ServerIndex <= 0 {
			OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, CSO2_POPUP_MEGAPHONE_USE_FAIL_INVALID_ITEM)
			DebugInfo(2, "User", uPtr.UserName, "use item", itemID, "failed,not in channel")
			return
		}

		idx, ok := uPtr.DecreaseItem(itemID)
		if !ok {
			OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, CSO2_UI_RAMDOMBOX_ALERT_000)
			DebugInfo(2, "User", uPtr.UserName, "use item", itemID, "failed")
			return
		}

		rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeInventory_Create),
			BuildInventoryInfoSingle(uPtr, 0, idx))
		SendPacket(rst, uPtr.CurrentConnection)

		// Check if user is already in queue
		for _, user := range megaphoneQueue {
			if user == uPtr.Userid {
				rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeUseItem), buildMegaphoneControl(1))
				SendPacket(rst, client)
				return
			}
		}

		// Add user to megaphone queue
		megaphoneQueue = append(megaphoneQueue, uPtr.Userid)

		// Send megaphone control packet with countdown timer
		rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeUseItem), buildMegaphoneControl(uint64(len(megaphoneQueue)-1)))
		SendPacket(rst, client)

		// Wait for user's turn
		for {
			if megaphoneQueue[0] == uPtr.Userid {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}

		const (
			Channel     = 2010 // A
			Server      = 2011 // A
			Global      = 2012 // A
			GlobalEvent = 2015 // A
		)

		// Send megaphone message
		msg := ""
		switch itemID {
		case 2010:
			msg = "1" // Global Message
		case 2011:
			msg = "2" // Server
		case 2012:
			msg = "3" // Channel
		case 2015:
			msg = "4"
		default:
			fmt.Println("Invalid value")
		}

		msg = fmt.Sprintf("%s: %s", uPtr.IngameName, string(pkt.String))

		// Sunucu kanal seçimi yapılmışsa, kanalda bulunan tüm oyunculara gönderilecek
		for _, v := range UsersManager.Users {
			if v != nil && v.GetUserChannelServerID() == v.GetUserChannelID() {
				OnSendMessageMegaphone(v.CurrentSequence, v.CurrentConnection, MessageAnnouncement, itemID, []byte(uPtr.IngameName), pkt.String)
			}
		}
		fmt.Println("Message sent:", msg)

		// Remove user from megaphone queue
		megaphoneQueue = megaphoneQueue[1:]

		// Update user inventory
		rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeInventory_Create), inventory.BuildInventoryInfoSingle(uPtr, 0, idx))
		SendPacket(rst, client)

		//发送消息

		//rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeUseItem),
		//	buildMegaphone(pkt.String))
		//SendPacket(rst, uPtr.CurrentConnection)

		DebugInfo(2, "User", uPtr.UserName, "say <", string(pkt.String), "> with item", itemID)
	case 2011: //服务器喇叭
	case 2012: //全体喇叭
	*/
	default:
		DebugInfo(2, "User", uPtr.UserName, "try using item but itemid is", itemID)
		return
	}

	//UserInfo部分
	rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeUserInfo), BuildUserInfo(0XFFFFFFFF, NewUserInfo(uPtr), uPtr.Userid, true))
	SendPacket(rst, uPtr.CurrentConnection)
}

func buildUsePoint(point uint32) []byte {
	buf := make([]byte, 25)
	offset := 0
	WriteUint8(&buf, useitem, &offset)
	WriteUint8(&buf, 5, &offset)
	WriteUint8(&buf, 1, &offset)      //unk00
	WriteUint32(&buf, 0, &offset)     //unk01
	WriteUint32(&buf, point, &offset) //mpoint
	return buf[:offset]
}

func buildNickNameChange() []byte {
	buf := make([]byte, 25)
	offset := 0
	WriteUint8(&buf, useitem, &offset)
	WriteUint8(&buf, 0, &offset)
	WriteUint8(&buf, 1, &offset)
	return buf[:offset]
}

func UsePointLotto(itemid uint32) uint64 {
	rand.Seed(time.Now().UnixNano())
	switch itemid {
	case 2008: //银币 id 2008
		return uint64(lotto_base + rand.Intn(lotto_max))
	case 2013: //铜币
		return uint64(lotto_event_base + rand.Intn(lotto_event_max))
	case 2014: //金币
		return uint64(lotto_gold_base + rand.Intn(lotto_gold_max))
	default:
		return 0
	}
}

func buildMegaphone(str []byte) []byte {
	temp := make([]byte, 256)
	offset := 0
	WriteUint8(&temp, useitem, &offset)
	WriteUint8(&temp, 6, &offset)
	WriteUint8(&temp, 1, &offset)
	WriteString(&temp, str, &offset)
	return temp[:offset]
}
