package configure

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
)

const (
	WeaponListCSV  = "/CSO2-Server/assert/cstrike/scripts/item_list.csv"
	ExpLevelCSV    = "/CSO2-Server/assert/cstrike/scripts/exp_level.csv"
	UnlockCSV      = "/CSO2-Server/assert/cstrike/scripts/item_unlock.csv"
	BoxCSV         = "/CSO2-Server/assert/cstrike/scripts/supplyList.csv"
	VipCSV         = "/CSO2-Server/assert/cstrike/scripts/vip_info.csv"
	DefaultItemCSV = "/CSO2-Server/assert/cstrike/scripts/defaultItemList.csv"
	ShopItemCSV    = "/CSO2-Server/assert/cstrike/scripts/shop.csv"
)

func InitCSV(path string) {
	fmt.Println("Reading game data file ...")
	readWeaponList(path)
	readUnlockList(path)
	readBoxList(path)
	readDefaultItemList(path)
	if Conf.EnableShop == 1 {
		readShopItemList(path)
	}
}

func readWeaponList(path string) {
	//读取武器数据
	filepath := path + WeaponListCSV

	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if err == nil && len(record[1]) > 16 {
			id, err := strconv.Atoi(record[0])
			if err != nil {
				continue
			}
			itemtype, err := strconv.Atoi(record[12])
			if err != nil {
				continue
			}
			itemd := ItemData{
				uint32(id),
				record[1][16:],
				record[4],
				record[5],
				record[6],
				itemtype,
			}

			ItemList[itemd.ItemID] = itemd
		} else {
			continue
		}
	}

}

func readUnlockList(path string) {
	//读取武器解锁数据
	filepath := path + UnlockCSV

	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if err == nil {
			id, err := strconv.Atoi(record[0])
			if err != nil {
				continue
			}
			nextid, err := strconv.Atoi(record[1])
			if err != nil {
				continue
			}
			flag0, err := strconv.Atoi(record[2])
			if err != nil {
				continue
			}
			count0, err := strconv.Atoi(record[3])
			if err != nil {
				continue
			}
			flag1, err := strconv.Atoi(record[4])
			if err != nil {
				continue
			}
			count1, err := strconv.Atoi(record[5])
			if err != nil {
				continue
			}
			flag2, err := strconv.Atoi(record[6])
			if err != nil {
				continue
			}
			count2, err := strconv.Atoi(record[7])
			if err != nil {
				continue
			}
			flag3, err := strconv.Atoi(record[8])
			if err != nil {
				continue
			}
			count3, err := strconv.Atoi(record[9])
			if err != nil {
				continue
			}
			flag4, err := strconv.Atoi(record[10])
			if err != nil {
				continue
			}
			count4, err := strconv.Atoi(record[11])
			if err != nil {
				continue
			}
			cat, err := strconv.Atoi(record[12])
			if err != nil {
				continue
			}
			itemd := UnlockData{
				uint32(id),
				uint32(nextid),
				uint32(flag0),
				uint32(count0),
				uint32(flag1),
				uint32(count1),
				uint32(flag2),
				uint32(count2),
				uint32(flag3),
				uint32(count3),
				uint32(flag4),
				uint32(count4),
				uint32(cat),
			}

			UnlockList[itemd.NextItemID] = itemd
		} else {
			continue
		}
	}
}

func readBoxList(path string) {
	//读取箱子数据
	filepath := path + BoxCSV

	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if err == nil {
			boxid, err := strconv.Atoi(record[0])
			if err != nil {
				continue
			}
			boxname := record[1]
			itemid, err := strconv.Atoi(record[2])
			if err != nil {
				continue
			}
			itemname := record[3]
			value, err := strconv.Atoi(record[4])
			if err != nil {
				continue
			}
			keyid, err := strconv.Atoi(record[5])
			if err != nil {
				continue
			}
			day, err := strconv.Atoi(record[6])
			if err != nil {
				continue
			}
			//保存当前物品数据
			if value <= 0 {
				fmt.Println("Warning ! illeagal value", value, "for item", itemid, "in box", boxid)
				continue
			}
			item := BoxItem{
				uint32(itemid),
				itemname,
				value,
				uint16(day),
			}

			if v, ok := BoxList[uint32(boxid)]; ok {
				//如果该box数据已经存在
				v.Items = append(v.Items, item)
				v.TotalValue += value
				BoxList[uint32(boxid)] = v
			} else {
				BoxList[uint32(boxid)] = BoxData{
					uint32(boxid),
					boxname,
					[]BoxItem{item},
					value,
					uint32(keyid),
					0,
				}
				BoxIDs = append(BoxIDs, uint32(boxid))
			}
		} else {
			continue
		}
	}
}

func readDefaultItemList(path string) {
	//读取数据
	filepath := path + DefaultItemCSV

	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	DefaultInventoryItem = []UserInventoryItem{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if err == nil {
			itemid, err := strconv.Atoi(record[0])
			if err != nil {
				continue
			}
			count, err := strconv.Atoi(record[2])
			if err != nil {
				continue
			}
			item := UserInventoryItem{uint32(itemid), uint16(count), 1, 0}
			DefaultInventoryItem = append(DefaultInventoryItem, item)
		} else {
			continue
		}
	}
}

func readShopItemList(path string) {
	//读取数据
	filepath := path + ShopItemCSV

	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if err == nil {
			itemid, err := strconv.Atoi(record[0])
			if err != nil {
				continue
			}
			price, err := strconv.Atoi(record[1])
			if err != nil {
				continue
			}
			currency, err := strconv.Atoi(record[2])
			if err != nil {
				continue
			}
			count, err := strconv.Atoi(record[3])
			if err != nil {
				continue
			}
			day, err := strconv.Atoi(record[4])
			if err != nil {
				continue
			}

			found := false
			for k, _ := range ShopItemList {
				if ShopItemList[k].ItemID == uint32(itemid) {
					ShopItemList[k].Opt = append(ShopItemList[k].Opt, ItemOption{uint32(price), uint16(count), uint16(day)})
					found = true
					break
				}
			}
			if !found {
				ShopItemList = append(ShopItemList, ShopItem{uint32(itemid), uint8(currency), []ItemOption{{uint32(price), uint16(count), uint16(day)}}})
			}
		} else {
			continue
		}
	}
}
