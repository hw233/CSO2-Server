package shop

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/core/inventory"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/message"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

type rewardBundle struct {
	itemid uint32
	count  uint16
	day    uint16
}

func OnShopBuyItem(p *PacketData, client net.Conn) {
	//检索数据包
	var pkt InShopBuyItemPacket
	if !p.PraseShopBuyItemPacket(&pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a error buyitem packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to request buyitem but not in server !")
		return
	}
	//找到购买的物品
	idx, optIdx := pkt.ItemID/5, pkt.ItemID%5
	if idx >= 0 && idx < uint32(len(ShopItemList)) && optIdx >= 0 && optIdx < uint32(len(ShopItemList[idx].Opt)) {
		//检查物品是否是套装
		if isBanItem(ShopItemList[idx].ItemID) {
			OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, CSO2_UI_RANK_CONDITION_HIDE)
			DebugInfo(2, "User", uPtr.UserName, "try to buy item", pkt.ItemID, "but is unimpletemented")
			return
		}
		if isBoughtOnce(ShopItemList[idx].ItemID, uPtr) {
			OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, CSO2_BUY_FAIL_BUY_CONDITION_FAILED_BY_FAIL_ITEM)
			DebugInfo(2, "User", uPtr.UserName, "try to buy item", pkt.ItemID, "but is bought befor")
			return
		}
		//发送数据
		switch ShopItemList[idx].Currency {
		case 0: //credit
			if !uPtr.UseCredits(ShopItemList[idx].Opt[optIdx].Price) {
				OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, CSO2_BUY_FAIL_CN_0X11_NO_CASH)
				DebugInfo(2, "User", uPtr.UserName, "try to buy item", pkt.ItemID, "but not enough cash")
				return
			}
		case 1: //point
			if !uPtr.UsePoints(uint64(ShopItemList[idx].Opt[optIdx].Price)) {
				OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, CSO2_BUY_FAIL_NO_POINT)
				DebugInfo(2, "User", uPtr.UserName, "try to buy item", pkt.ItemID, "but not enough points")
				return
			}
		case 2: //mpoint
			if !uPtr.UseMPoints(ShopItemList[idx].Opt[optIdx].Price) {
				OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, CSO2_BUY_FAIL_NO_MILEAGE)
				DebugInfo(2, "User", uPtr.UserName, "try to buy item", pkt.ItemID, "but not enough mpoints")
				return
			}
		default:
			OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, CSO2_BUY_FAIL_CN_0X12_DB_ERROR)
			DebugInfo(2, "User", uPtr.UserName, "try to buy item", pkt.ItemID, "but unkown currency", ShopItemList[idx].Currency)
			return
		}
		ok, rewards := isBundle(ShopItemList[idx].ItemID, uPtr, idx, optIdx)
		if ok {
			OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, CSO2_BUY_SUCCEED)
			for _, v := range rewards {
				//发送得到的物品
				itemIdx := uPtr.AddItem(v.itemid, v.count, uint64(v.day))
				rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeInventory_Create),
					BuildInventoryInfoSingle(uPtr, 0, itemIdx))
				SendPacket(rst, uPtr.CurrentConnection)
			}
			uPtr.SetBoughtItem(ShopItemList[idx].ItemID)
		} else {
			//发送得到的物品
			itemIdx := uPtr.AddItem(ShopItemList[idx].ItemID, ShopItemList[idx].Opt[optIdx].Count, uint64(ShopItemList[idx].Opt[optIdx].Day))
			OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, CSO2_BUY_SUCCEED)
			rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeInventory_Create),
				BuildInventoryInfoSingle(uPtr, 0, itemIdx))
			SendPacket(rst, uPtr.CurrentConnection)
		}
		//UserInfo部分
		rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeUserInfo), BuildUserInfo(0xFFFFFFFF, NewUserInfo(uPtr), uPtr.Userid, true))
		SendPacket(rst, uPtr.CurrentConnection)
		DebugInfo(2, "User", uPtr.UserName, "bought item", ShopItemList[idx].ItemID, "from shop")
		return
	}
	//未找到物品
	OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageDialogBox, CSO2_BUY_FAIL_CN_0X12_DB_ERROR)
	DebugInfo(2, "User", uPtr.UserName, "try to buy item idx", pkt.ItemID, "but failed")
}

func isBundle(itemId uint32, user *User, idx, optIdx uint32) (bool, []rewardBundle) {
	switch itemId {
	case 100090: //冰雪牛仔礼包1
		return true, []rewardBundle{{10128, 1, 1}, {55015, 1, 0}}
	case 100091: //冰雪牛仔礼包2
		return true, []rewardBundle{{10128, 1, 2}, {55015, 10, 0}}
	case 100092: //冰雪牛仔礼包3
		return true, []rewardBundle{{10128, 1, 3}, {55015, 50, 0}}
	case 100093: //冰雪牛仔礼包4
		return true, []rewardBundle{{10128, 1, 4}, {55015, 100, 0}}
	case 100048: //黑色牛仔礼包1
		return true, []rewardBundle{{10106, 1, 1}, {55001, 1, 0}}
	case 100049: //黑色牛仔礼包2
		return true, []rewardBundle{{10106, 1, 2}, {55001, 10, 0}}
	case 100050: //黑色牛仔礼包3
		return true, []rewardBundle{{10106, 1, 3}, {55001, 50, 0}}
	case 100051: //黑色牛仔礼包4
		return true, []rewardBundle{{10106, 1, 4}, {55001, 100, 0}}
	case 100056: //橘红牛仔礼包1
		return true, []rewardBundle{{10108, 1, 1}, {55007, 1, 0}}
	case 100057: //橘红牛仔礼包2
		return true, []rewardBundle{{10108, 1, 2}, {55007, 10, 0}}
	case 100058: //橘红牛仔礼包3
		return true, []rewardBundle{{10108, 1, 3}, {55007, 50, 0}}
	case 100059: //橘红牛仔礼包4
		return true, []rewardBundle{{10108, 1, 4}, {55007, 100, 0}}
	case 100060: //圣诞牛仔礼包1
		return true, []rewardBundle{{10109, 1, 1}, {55003, 1, 0}}
	case 100061: //圣诞牛仔礼包2
		return true, []rewardBundle{{10109, 1, 2}, {55003, 10, 0}}
	case 100062: //圣诞牛仔礼包3
		return true, []rewardBundle{{10109, 1, 3}, {55003, 50, 0}}
	case 100063: //圣诞牛仔礼包4
		return true, []rewardBundle{{10109, 1, 4}, {55003, 100, 0}}
	case 100078: //冰橙牛仔礼包1
		return true, []rewardBundle{{10116, 1, 1}, {55012, 1, 0}}
	case 100079: //冰橙牛仔礼包2
		return true, []rewardBundle{{10116, 1, 2}, {55012, 10, 0}}
	case 100080: //冰橙牛仔礼包3
		return true, []rewardBundle{{10116, 1, 3}, {55012, 50, 0}}
	case 100081: //冰橙牛仔礼包4
		return true, []rewardBundle{{10116, 1, 4}, {55012, 100, 0}}
	case 200018: //marble 汉白玉
		return true, []rewardBundle{{5326, 1, 0}, {5327, 1, 0}, {5328, 1, 0}, {5329, 1, 0}}
	case 200017: //黄金TR礼包
		return true, []rewardBundle{{5322, 1, 0}, {5323, 1, 0}}
	case 200016: //黄金CT礼包
		return true, []rewardBundle{{5324, 1, 0}, {5325, 1, 0}}
	case 100077: //财运滚滚
		user.GetPoints(220000)
		return true, []rewardBundle{{5305, 1, 0}}
	case 100076: //金猪福运
		user.GetPoints(300000)
		return true, []rewardBundle{{20087, 1, 0}}
	case 100035: //榴弹礼包
		return true, []rewardBundle{{5233, 1, 0}, {5234, 1, 0}, {5235, 1, 0}, {5236, 1, 0}}
	case 100037: //阿宽/哈桑套装
		return true, []rewardBundle{{1030, 1, 0}, {1031, 1, 0}}
	case 100031: //炽焰扳手套装
		return true, []rewardBundle{{5221, 1, 0}, {50017, 20, 0}}
	case 200019: //萌犬套装
		return true, []rewardBundle{{5365, 1, 0}, {10131, 1, 0}, {20107, 1, 0}, {30027, 1, 0}}
	case 200013: //粉红奈奈礼包
		return true, []rewardBundle{{5307, 1, 0}, {5308, 1, 0}, {5309, 1, 0}, {5310, 1, 0}, {5311, 1, 0}}
	case 200015: //三合一礼包
		return true, []rewardBundle{{5315, 1, 0}, {5316, 1, 0}, {5317, 1, 0}}
	case 5000078: //UMP45
		return true, []rewardBundle{{5346, 1, ShopItemList[idx].Opt[optIdx].Day}}
	case 5000079: //KRISS SUPER
		return true, []rewardBundle{{5347, 1, ShopItemList[idx].Opt[optIdx].Day}}
	case 5000080: //MP7 RUBY
		return true, []rewardBundle{{5348, 1, ShopItemList[idx].Opt[optIdx].Day}}
	case 5000060: //莉莉娅
		return true, []rewardBundle{{1043, 1, ShopItemList[idx].Opt[optIdx].Day}}
	case 5000075: //粉红双肩包
		return true, []rewardBundle{{20104, 1, ShopItemList[idx].Opt[optIdx].Day}}
	case 5000076: //紫色玩偶兔
		return true, []rewardBundle{{20105, 1, ShopItemList[idx].Opt[optIdx].Day}}
	case 5000085: //皮草帽
		return true, []rewardBundle{{10129, 1, ShopItemList[idx].Opt[optIdx].Day}}
	case 5000086: //毛绒帽
		return true, []rewardBundle{{10130, 1, ShopItemList[idx].Opt[optIdx].Day}}
	case 100027: //紫色军刀100个
		return true, []rewardBundle{{55004, 100, 0}}
	case 100034: //西瓜军刀100个
		return true, []rewardBundle{{55005, 100, 0}}
	case 5000081: //m99 railgun
		return true, []rewardBundle{{150, 1, ShopItemList[idx].Opt[optIdx].Day}}
	case 5000064: //同级生 奈奈
		return true, []rewardBundle{{1053, 1, ShopItemList[idx].Opt[optIdx].Day}}
	case 150014: //黄金狩猎刀
		return true, []rewardBundle{{79, 1, ShopItemList[idx].Opt[optIdx].Day}}
	case 150038: //电击刀
		return true, []rewardBundle{{89, 1, ShopItemList[idx].Opt[optIdx].Day}}
	case 5000063: //骷髅手套
		return true, []rewardBundle{{30005, 1, ShopItemList[idx].Opt[optIdx].Day}}
	default:
		return false, nil
	}

}

func isBoughtOnce(itemId uint32, user *User) bool {
	if user == nil {
		return false
	}
	switch itemId {
	case 200018, 200017, 200016, 100031, 100037, 200013, 200015, 100076, 100077:
		if user.Inventory.OnlyOnceBundleItemId[itemId] {
			return true
		}
	default:
		return false
	}
	return false
}

func isBanItem(itemId uint32) bool {
	switch itemId {
	case 100031:
		return true
	default:
		return false
	}
}
