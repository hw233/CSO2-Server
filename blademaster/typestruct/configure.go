package typestruct

type ItemData struct {
	ItemID      uint32
	Name        string
	Class       string
	Category    string
	BuyCategory string
	ItemType    int
}

type UnlockData struct {
	ItemID         uint32
	NextItemID     uint32
	ConditionFlag0 uint32
	Count0         uint32
	ConditionFlag1 uint32
	Count1         uint32
	ConditionFlag2 uint32
	Count2         uint32
	ConditionFlag3 uint32
	Count3         uint32
	ConditionFlag4 uint32
	Count4         uint32
	Category       uint32
}

type BoxData struct {
	BoxID      uint32
	BoxName    string
	Items      []BoxItem
	TotalValue int
	KeyId      uint32
	OptId      uint32
}

type BoxItem struct {
	ItemID   uint32
	ItemName string
	Value    int
	Day      uint16
}

type ShopItem struct {
	ItemID   uint32
	Currency uint8
	Opt      []ItemOption
}

type ItemOption struct {
	Price uint32
	Count uint16
	Day   uint16
}

var (
	ItemList     = make(map[uint32]ItemData)
	UnlockList   = make(map[uint32]UnlockData)
	BoxList      = make(map[uint32]BoxData)
	BoxIDs       = []uint32{}
	ShopItemList = []ShopItem{}
)
