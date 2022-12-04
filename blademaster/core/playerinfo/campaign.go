package playerinfo

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/core/achievement"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/inventory"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

const (
	started  = 0
	finished = 1
)

func OnSetCampaign(p *PacketData, client net.Conn) {
	var pkt InSetCampaignPacket
	if !p.PraseSetCampaignPacket(&pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a illegal SetCampaign packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to SetCampaign but not in server !")
		return
	}
	//修改数据
	switch pkt.PacketType {
	case started:
		DebugInfo(1, "User", uPtr.UserName, "Started Campaign")
	case finished:
		if isMissionCampaignIdValid(pkt.CampaignId) {
			if uPtr.CheckCampaign(pkt.CampaignId) {
				DebugInfo(1, "User", uPtr.UserName, "already finished Campaign ", pkt.CampaignId)
				return
			}
			uPtr.UpdateCampaign(pkt.CampaignId)
			//发送数据包
			rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeUserInfo), BuildUserInfo(0x1000, NewUserInfo(uPtr), uPtr.Userid, true))
			SendPacket(rst, uPtr.CurrentConnection)
			switch pkt.CampaignId {
			case Campaign_1:
				for _, v := range RewardCapmgaign1 {
					idx := uPtr.AddItem(v.ItemId, 1, 0)
					rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeInventory_Create),
						BuildInventoryInfoSingle(uPtr, 0, idx))
					SendPacket(rst, uPtr.CurrentConnection)
				}
				uPtr.GetExp(3000)
				rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeUserInfo), BuildUserInfo(0XFFFFFFFF, NewUserInfo(uPtr), uPtr.Userid, true))
				SendPacket(rst, uPtr.CurrentConnection)
				//OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, GAME_USER_NEW_ITEM_RESTART)
			case Campaign_2:
				for _, v := range RewardCapmgaign2 {
					idx := uPtr.AddItem(v.ItemId, 1, 0)
					rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeInventory_Create),
						BuildInventoryInfoSingle(uPtr, 0, idx))
					SendPacket(rst, uPtr.CurrentConnection)
				}
				uPtr.GetPoints(5000)
				rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeUserInfo), BuildUserInfo(0XFFFFFFFF, NewUserInfo(uPtr), uPtr.Userid, true))
				SendPacket(rst, uPtr.CurrentConnection)
				//OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, GAME_USER_NEW_ITEM_RESTART)
			case Campaign_3:
				for _, v := range RewardCapmgaign3 {
					idx := uPtr.AddItem(v.ItemId, 1, 0)
					rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeInventory_Create),
						BuildInventoryInfoSingle(uPtr, 0, idx))
					SendPacket(rst, uPtr.CurrentConnection)
				}
				uPtr.GetExp(8000)
				rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeUserInfo), BuildUserInfo(0XFFFFFFFF, NewUserInfo(uPtr), uPtr.Userid, true))
				SendPacket(rst, uPtr.CurrentConnection)
				//OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, GAME_USER_NEW_ITEM_RESTART)
			case Campaign_4:
				for _, v := range RewardCapmgaign4 {
					idx := uPtr.AddItem(v.ItemId, 1, 0)
					rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeInventory_Create),
						BuildInventoryInfoSingle(uPtr, 0, idx))
					SendPacket(rst, uPtr.CurrentConnection)
				}
				uPtr.GetExp(10000)
				rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeUserInfo), BuildUserInfo(0XFFFFFFFF, NewUserInfo(uPtr), uPtr.Userid, true))
				SendPacket(rst, uPtr.CurrentConnection)
				//OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, GAME_USER_NEW_ITEM_RESTART)
			case Campaign_5:
				for _, v := range RewardCapmgaign5 {
					idx := uPtr.AddItem(v.ItemId, 1, 0)
					rst = BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeInventory_Create),
						BuildInventoryInfoSingle(uPtr, 0, idx))
					SendPacket(rst, uPtr.CurrentConnection)
				}
			}
			DebugInfo(1, "User", uPtr.UserName, "finished Campaign ", pkt.CampaignId)
		} else {
			DebugInfo(2, "User", uPtr.UserName, "sent a invalid SetCampaign packet id", pkt.CampaignId)
		}
	default:
		DebugInfo(2, "User", uPtr.UserName, "sent a unkown SetCampaign packet type", pkt.PacketType)
	}
}

func isMissionCampaignIdValid(id uint16) bool {
	return id == Campaign_0 ||
		id == Campaign_1 ||
		id == Campaign_2 ||
		id == Campaign_3 ||
		id == Campaign_4 ||
		id == Campaign_5 ||
		id == Campaign_6
}
