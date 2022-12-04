package inventory

import (
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/configure"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

var (
	FullInventoryItem  []UserInventoryItem
	FullInventoryReply []byte

	numPerDay = float64(24*3600) / 256
)

func BuildInventoryInfo(u *User) []byte {
	if Conf.UnlockAllWeapons == 1 {
		return FullInventoryReply
	}
	buf := make([]byte, 5+u.Inventory.NumOfItem*19)
	offset, offset_front := 2, 0
	num := 0

	for k, v := range u.Inventory.Items {
		if v.Count <= 0 {
			continue
		}
		WriteUint16(&buf, uint16(k), &offset)
		WriteUint8(&buf, 1, &offset)
		WriteUint32(&buf, v.Id, &offset)
		WriteUint16(&buf, v.Count, &offset)
		WriteUint8(&buf, 1, &offset)
		WriteUint8(&buf, 0, &offset)
		WriteUint64(&buf, v.Time, &offset)

		num++
	}
	WriteUint16(&buf, uint16(num), &offset_front) //把实际的数量写进前部
	return buf[:offset]
}

func BuildInventoryInfoSingle(u *User, itemid uint32, idx int) []byte {
	buf := make([]byte, 128)
	offset := 0
	WriteUint16(&buf, 1, &offset)
	if idx < 0 || idx >= len(u.Inventory.Items) {
		for k, v := range u.Inventory.Items {
			if v.Id != itemid {
				continue
			}
			WriteUint16(&buf, uint16(k), &offset)
			if v.Count == 0 {
				WriteUint8(&buf, 0, &offset) //existed
			} else {
				WriteUint8(&buf, 1, &offset)
			}
			WriteUint32(&buf, v.Id, &offset)
			WriteUint16(&buf, v.Count, &offset)
			WriteUint8(&buf, 1, &offset)
			WriteUint8(&buf, 0, &offset)
			WriteUint64(&buf, v.Time, &offset)
			break
		}
	} else {
		WriteUint16(&buf, uint16(idx), &offset)
		if u.Inventory.Items[idx].Count == 0 {
			WriteUint8(&buf, 0, &offset) //existed
		} else {
			WriteUint8(&buf, 1, &offset)
		}
		WriteUint32(&buf, u.Inventory.Items[idx].Id, &offset)
		WriteUint16(&buf, u.Inventory.Items[idx].Count, &offset)
		WriteUint8(&buf, 1, &offset)
		WriteUint8(&buf, 0, &offset)
		WriteUint64(&buf, u.Inventory.Items[idx].Time, &offset)
	}
	return buf[:offset]
}

func BuildFullInventoryInfo() []byte {
	buf := make([]byte, 5+uint16(len(FullInventoryItem))*19)
	offset := 0
	WriteUint16(&buf, uint16(len(FullInventoryItem)), &offset)
	for k, v := range FullInventoryItem {
		WriteUint16(&buf, uint16(k), &offset)
		WriteUint8(&buf, 1, &offset)
		WriteUint32(&buf, v.Id, &offset)
		WriteUint16(&buf, v.Count, &offset)
		WriteUint8(&buf, 1, &offset)
		WriteUint8(&buf, 0, &offset)
		WriteUint64(&buf, v.Time, &offset)
	}
	return buf[:offset]
}

func BuildDefaultInventoryInfo() []byte {
	buf := make([]byte, 5+len(DefaultInventoryItem)*19)
	offset := 0
	WriteUint16(&buf, 25, &offset)
	for k, v := range DefaultInventoryItem {
		WriteUint16(&buf, uint16(k), &offset)
		WriteUint8(&buf, 1, &offset)
		WriteUint32(&buf, v.Id, &offset)
		WriteUint16(&buf, v.Count, &offset)
		WriteUint8(&buf, 1, &offset)
		WriteUint8(&buf, 0, &offset)
		WriteUint64(&buf, v.Time, &offset)

	}
	return buf[:offset]
}

func BuildUnlockReply(u *User) []byte {
	if u == nil {
		return []byte{}
	}
	buf := make([]byte, 4096)
	offset := 0
	WriteUint8(&buf, 1, &offset)                            //type ?
	WriteUint16(&buf, uint16(len(UnlockFullList)), &offset) //num of weapons
	for k, v := range UnlockFullList {
		WriteUint32(&buf, v.Itemid, &offset)
		WriteUint32(&buf, uint32(k), &offset)
		WriteUint8(&buf, v.CostType, &offset)
		WriteUint32(&buf, v.Price, &offset)
	}

	WriteUint16(&buf, 1, &offset) //num of weapons

	WriteUint32(&buf, 2, &offset) //前置
	WriteUint32(&buf, 1, &offset) //当前
	WriteUint32(&buf, 2, &offset) //杀敌数
	WriteUint16(&buf, 1, &offset)
	WriteUint16(&buf, 1, &offset)
	WriteUint16(&buf, 1, &offset)

	WriteUint16(&buf, 1, &offset) //unk

	WriteUint32(&buf, 1, &offset)

	return buf[:offset]
}

// func BuildWeaponKillNum(u *User) []byte {
// 	if u == nil {
// 		return []byte{}
// 	}
// 	buf := make([]byte, 4096)
// 	offset := 0
// 	count := 0
// 	WriteUint16(&buf, uint16(count), &offset)
// 	for _, v := range UnlockFullList {
// 		if _, ok := u.WeaponKills[v.Itemid]; ok && u.WeaponKills[v.Itemid] > 0 {
// 			count++
// 			WriteUint32(&buf, 2, &offset) //前置
// 			WriteUint32(&buf, 1, &offset) //当前
// 			WriteUint32(&buf, 2, &offset) //杀敌数
// 			WriteUint16(&buf, 1, &offset)
// 			WriteUint16(&buf, 1, &offset)
// 			WriteUint16(&buf, 1, &offset)
// 		}
// 	}
// 	return buf[:offset]

// }

func CreateFullUserInventory() UserInventory {
	Inv := UserInventory{
		0,
		//createDeafaultInventoryItem(),
		CreateFullInventoryItem(),
		map[uint32]bool{},
		1047,
		1048,
		0,
		0,
		0,
		0,
		0,
		42001,
		CreateFullUserBuyMenu(),
		CreateFullLoadout(),
	}
	Inv.NumOfItem = uint16(len(Inv.Items))
	return Inv
}

func CreateFullInventoryItem() []UserInventoryItem {
	items := []UserInventoryItem{}
	// for _, v := range ItemList {
	// 	if v.Category == "weapon" || v.Category == "class" {
	// 		items = append(items, UserInventoryItem{v.ItemID, 1})
	// 	}
	// }
	var i uint32
	//用户角色
	for i = 1001; i <= 1058; i++ {
		items = append(items, UserInventoryItem{i, 1, 1, 0})
	}
	//添加默认武器
	number := []uint32{2, 3, 4, 6, 8, 13, 14, 15, 18, 19, 21, 23, 27, 34, 36, 37, 80, 128, 101, 49009, 49004}
	for _, v := range number {
		items = append(items, UserInventoryItem{v, 1, 1, 0})
	}
	//解锁武器
	for i = 1; i <= 33; i++ {
		if IsIllegal(i) {
			continue
		}
		items = append(items, UserInventoryItem{i, 1, 1, 0})
	}
	for i = 44; i <= 163; i++ {
		if IsIllegal(i) {
			continue
		}
		items = append(items, UserInventoryItem{i, 1, 1, 0})
	}
	//隔离区技能
	items = append(items, UserInventoryItem{2019, 1, 1, 0})
	items = append(items, UserInventoryItem{2020, 1, 1, 0})
	for i = 2021; i <= 2023; i++ {
		items = append(items, UserInventoryItem{i, 1, 1, 0})
	}
	//武器皮肤
	for i = 5042; i <= 5370; i++ {
		if IsIllegal(i) {
			continue
		}
		items = append(items, UserInventoryItem{i, 1, 1, 0})
	}
	items = append(items, UserInventoryItem{5997, 1, 1, 0})
	//帽子
	for i = 10001; i <= 10133; i++ {
		items = append(items, UserInventoryItem{i, 1, 1, 0})
	}
	//背包
	for i = 20001; i <= 20107; i++ {
		items = append(items, UserInventoryItem{i, 1, 1, 0})
	}
	//手套
	for i = 30001; i <= 30027; i++ {
		items = append(items, UserInventoryItem{i, 1, 1, 0})
	}
	//脚部特效
	for i = 40001; i <= 40025; i++ {
		items = append(items, UserInventoryItem{i, 1, 1, 0})
	}
	//喷漆
	for i = 42001; i <= 42020; i++ {
		items = append(items, UserInventoryItem{i, 1, 1, 0})
	}
	//道具
	for i = 49001; i <= 49010; i++ {
		items = append(items, UserInventoryItem{i, 1, 1, 0})
	}
	items = append(items, UserInventoryItem{49999, 1, 1, 0})
	//角色卡片
	for i = 60001; i <= 60004; i++ {
		items = append(items, UserInventoryItem{i, 1, 1, 0})
	}
	return items
}
