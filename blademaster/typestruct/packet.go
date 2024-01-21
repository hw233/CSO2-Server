package typestruct

import (
	"log"
	"math"
	"net"
	"unsafe"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type (
	//PacketHeader ,header of packet , 4 bytes len
	PacketHeader struct {
		Data     []byte
		Sequence uint8
		Length   uint16
	}
	//PacketData ,data part of packet
	PacketData struct {
		Data      []byte
		Sequence  uint8
		Length    uint16
		Id        uint8
		CurOffset int
	}

	InLoginPacket struct {
		LenOfNexonUsername uint8
		NexonUsername      []byte //假定nexonUsername是唯一
		LenOfGameUsername  uint8
		GameUsername       []byte
		Unknown01          uint8
		LenOfPassWd        uint16
		PassWd             []byte
		//HddHwid 	   	   16 bytes
		HddHwid []byte
		//netCafeID 	   4 bytes
		NetCafeID          uint32
		Unknown02          uint32
		UserSn             uint64
		LenOfUnknownString uint16
		UnknownString03    []byte
		Unknown04          uint8
		IsLeague           uint8
		LenOfString        uint8
		String             []byte
	}

	//房间请求
	InRoomPaket struct {
		InRoomType uint8
	}

	InShopPacket struct {
		InShopType uint8
	}
	InShopBuyItemPacket struct {
		ItemID uint32
	}
	InSupplyPacket struct {
		Type uint8
	}
	InPointLottoPacket struct {
		Type uint8
	}
	InOpenBoxPacket struct {
		BoxID uint32
		Unk00 uint32
	}
	InItemUsePacket struct {
		unk00       uint8
		ItemSeq     uint16
		ItemType    uint8
		lenOfString uint8
		String      []byte
	}
	InTryItemUsePacket struct {
		unk00        uint8
		ItemSeq      uint16
		ItemType     uint8
		lenOfNewName uint8
		NewName      []byte
	}
	//InRoomListRequestPacket 房间列表请求，用于请求频道
	InRoomListRequestPacket struct {
		ChannelServerIndex uint8
		ChannelIndex       uint8
	}
	InReportPacket struct {
		Type uint8
	}
	InReportSearchUserPacket struct {
		LenOfName uint8
		Name      []byte
	}
	InReportMsgPacket struct {
		LenOfName uint8
		Name      []byte
		Type      uint8
		LenOfMsg  uint16
		Msg       []byte
		Unk00     uint16
	}
	InFavoritePacket struct {
		PacketType uint8
	}

	InFavoriteSetCosmetics struct {
		Slot   uint8
		ItemId uint32
	}

	InHostPacket struct {
		InHostType uint8
	}
	InNotifyPacket struct {
		InNotifyType uint8
	}
	InNotifyListPacket struct {
		InType uint8
	}

	InQuickPacket struct {
		InQuickType uint8
	}

	InQuickList struct {
		GameModID   uint8
		IsEnableBot uint8
	}

	InChatPacket struct {
		Type           uint8
		DestinationLen uint8
		Destination    []byte
		MessageLen     uint16
		Message        []byte
	}

	InAchievementPacket struct {
		Type uint8
	}

	InEventPacket struct {
		Type uint8
	}

	InAchievementCampaignPacket struct {
		CampaignId uint16
	}

	InDisassemblePacket struct {
		Type uint8
	}

	InDisassembleItemPacket struct {
		SubType uint8
	}
	InDisassembleWeaponPacket struct {
		Unk00  uint32
		Unk01  uint8
		ItemID uint32
		Unk02  uint32
		Unk03  uint32
	}

	//InNewRoomPacket 新建房间时传进来的数据包
	InNewRoomPacket struct {
		LenOfName  uint8
		RoomName   []byte
		Unk00      uint16
		Unk01      uint8
		GameModeID uint8
		MapID      uint8
		WinLimit   uint8
		KillLimit  uint16
		Unk02      uint8
		Unk03      uint8
		Unk04      uint8
		LenOfUnk05 uint8
		Unk05      []byte
		Unk06      uint8
		Unk07      uint8
		Unk08      uint8
		Unk09      uint8
		Unk10      uint8
		Unk11      uint32
	}
	InUpSettingReq struct {
		Flags              uint64
		LenOfRoomName      uint8
		RoomName           []byte
		Unk00              uint8
		Unk01              uint8
		Unk02              uint32
		Unk03              uint32
		LenOfpasswd        uint8
		Passwd             []byte
		Unk10              uint16
		ForceCamera        uint8
		GameModeID         uint8
		MapID              uint8
		Unk13              uint8
		MaxPlayers         uint8
		WinLimit           uint8
		KillLimit          uint16
		Unk17              uint8
		Unk18              uint8
		WeaponRestrictions uint8
		Status             uint8
		Unk21              uint8
		MapCycleType       uint8
		Unk23              uint8
		Unk24              uint8
		Unk25              uint8
		NumOfMultiMaps     uint8
		MultiMaps          []uint8
		TeamBalanceType    uint8
		Unk29              uint8
		Unk30              uint8
		Unk31              uint8
		Unk32              uint8
		Unk33              uint8
		BotEnabled         uint8
		BotDifficulty      uint8
		NumCtBots          uint8
		NumTrBots          uint8
		Unk35              uint8
		Unk36              uint8
		Unk37              uint8
		Unk38              uint8
		Unk39              uint8
		IsIngame           uint8
		StartMoney         uint16
		ChangeTeams        uint8
		Unk43              uint8
		HltvEnabled        uint8
		Unk45              uint8
		RespawnTime        uint8
	}

	InRoomCountdownPacket struct {
		CountdownType uint8
		Count         uint8
	}

	InJoinRoomPacket struct {
		RoomId        uint16
		LenOfPassWord uint8
		PassWord      []byte
	}

	InFeedbackPacket struct {
		unk00     uint16
		ErrorCode uint16
		unk01     uint32
	}

	InChangeTeamPacket struct {
		NewTeam uint8
	}

	InFavoriteSetLoadout struct {
		Loadout    uint8
		WeaponSlot uint8
		ItemId     uint32
	}

	InPlayerInfoPacket struct {
		InfoType uint8
	}

	InSetSignaturePacket struct {
		Len       uint8
		Signature []byte
	}

	InSetAvatarPacket struct {
		AvatarId uint16
	}

	InSetTitlePacket struct {
		TitleId uint16
	}

	InSetCampaignPacket struct {
		PacketType uint8
		CampaignId uint16
	}

	InOptionPacket struct {
		OptionPacketType uint8
	}

	InOptionBuyMenu struct {
		MenuLength uint16
		Unk00      uint8
		Buymenu    UserBuyMenu
	}

	InKillPacket struct {
		Unk00      uint8 //一直是0
		KillerID   uint32
		Unk01      uint32 //一直是0
		Unk02      uint8
		KillType   uint8  //貌似是击杀方式？
		KillNum    uint16 //杀敌数,生化模式3倍
		PlayerTeam uint8  //待定
	}

	InWeaponPointPacket struct {
		KillType             uint32
		KillerID             uint32
		KillerWeaponID       uint32
		KillerTeam           uint8
		KillerClientType     uint8
		KillerCharacterType  uint8
		KillerCharacterClass uint32

		DeadType             uint32
		VictimID             uint32
		VictimWeaponID       uint32
		VictimTeam           uint8
		VictimClientType     uint8
		VictimCharacterType  uint8
		VictimCharacterClass uint32

		KillerX uint32
		KillerY uint32
		KillerZ uint32

		VictimX uint32
		VictimY uint32
		VictimZ uint32

		ArraySize uint16 //待定
		Array     []InWeaponPointArray
	}

	InWeaponPointArray struct {
		unk00 uint8
		unk01 uint32
		unk02 uint32
		unk03 uint8
	}

	InDeathPacket struct {
		DeadID     uint32
		Unk00      uint32 //一直是0
		Unk01      uint8  //貌似是死亡方式？
		DeathNum   uint16 //死亡数,生化模式3倍
		PlayerTeam uint8  //待定
	}

	InRevivedPacket struct {
		UserID uint32
		X      uint32 //待定，但是极像坐标
		Y      uint32
		Z      uint32
		Unk00  uint8
	}

	InAssistPacket struct {
		KillerID     uint32
		Unk00        uint8  //可能是辅助击杀人数？
		AssisterID   uint32 //貌似是击杀方式？
		Unk01        uint16
		Unk02        uint16
		Unk03        uint16
		AssisterTeam uint8 //待定,也可能是杀手的队伍
	}
	InHostSetBuyMenu struct {
		Userid uint32
	}

	InHostTeamChangingPacket struct {
		UserId uint32
		Unk00  uint8
		//unk01   uint8
		NewTeam uint8
	}

	InHostSetLoadoutPacket struct {
		UserID uint32
	}

	InGameScorePacket struct {
		WinnerTeam uint8
		TrScore    uint8
		CtScore    uint8
		PacketType uint8 //maybe
		HostID     uint32
		Unk00      uint32
	}

	InHostSetInventoryPacket struct {
		UserID uint32
	}

	InHostItemUsingPacket struct {
		UserID uint32
		ItemID uint32
	}

	InHostBuyItemPacket struct {
		UserID         uint32
		NumItemsBought uint8
		Items          []uint32
	}

	//未知，用于请求频道
	OutLobbyJoinRoom struct {
		Unk00 uint8
		Unk01 uint8
		Unk02 uint8
	}
	//OutAchievementCampaign from l-leite
	OutAchievementCampaign struct {
		Unk00        uint16
		Unk01        uint32
		RewardTitle  uint16
		RewardIcon   uint16
		RewardPoints uint32
		RewardXp     uint32
		NumOfItems   uint8
		Items        []OutAchievementCampaignItems
		Unk02        uint16
	}
	//OutAchievementCampaignItems from l-leite
	OutAchievementCampaignItems struct {
		ItemId      uint32
		Ammount     uint16
		TimeLimited uint16 // 1 and 0
	}

	InAchievementBackgroundPacket struct {
		ACHType      uint8
		SelectedCase uint8
		Type         uint8
		Type2        uint8
		Type3        uint8
	}
)

const (
	//packet's first main type
	PacketTypeVersion          = 0
	PacketTypeReply            = 1
	PacketTypeLogin            = 3
	PacketTypeServerList       = 5
	PacketTypeCharacter        = 6
	PacketTypeRequestRoomList  = 7
	PacketTypeRequestChannels  = 10
	PacketTypeRoom             = 65
	PacketTypeChat             = 67
	PacketTypeHost             = 68
	PacketTypePlayerInfo       = 69
	PacketTypeUdp              = 70
	PacketTypeShop             = 72
	PacketTypeBan              = 74
	PacketTypeOption           = 76
	PacketTypeFavorite         = 77
	PacketTypeUseItem          = 78
	PacketTypeQuickJoin        = 80
	PacketTypeReport           = 83
	PacketTypeSignature        = 85
	PacketTypeQuickStart       = 86
	PacketTypeAutomatch        = 88
	PacketTypeFriend           = 89
	PacketTypeUnlock           = 90
	PacketTypeNotify           = 91
	PacketTypeEvent            = 92
	PacketTypeGZ               = 95
	PacketTypeAchievement      = 96
	PacketTypeSupply           = 102
	PacketTypeDisassemble      = 104
	PacketTypeConfigInfo       = 106
	PacketTypeLobby            = 107
	PacketTypeUserStart        = 150
	PacketTypeRoomList         = 151
	PacketTypeInventory_Add    = 152
	PacketTypeInventory_Create = 154
	PacketTypeUserInfo         = 157
	//beacuse there is only 1 byte of sequence in packet , so max number is 0xff
	MINSEQUENCE = 0
	MAXSEQUENCE = math.MaxUint8
	//server will read 4 bytes of header
	HeaderLen = 4

	SUBMENU_ITEM_NUM = 9
)

func (p *PacketHeader) PraseHeadPacket() {
	// if p.Data[0] != PacketTypeSignature {
	// 	p.IsGoodPacket = false
	// 	return
	// }
	// p.IsGoodPacket = true
	offset := 0
	p.Sequence = ReadUint8(p.Data, &offset)
	p.Length = ReadUint16(p.Data, &offset)
}

func (dataPacket *PacketData) PraseLoginPacket(p *InLoginPacket) bool {
	if dataPacket.Length < 50 {
		return false
	}
	p.LenOfNexonUsername = ReadUint8(dataPacket.Data, &dataPacket.CurOffset)
	p.NexonUsername = ReadString(dataPacket.Data, &dataPacket.CurOffset, int(p.LenOfNexonUsername))
	p.LenOfGameUsername = ReadUint8(dataPacket.Data, &dataPacket.CurOffset)
	p.GameUsername = ReadString(dataPacket.Data, &dataPacket.CurOffset, int(p.LenOfGameUsername))
	p.Unknown01 = ReadUint8(dataPacket.Data, &dataPacket.CurOffset)
	p.LenOfPassWd = ReadUint16(dataPacket.Data, &dataPacket.CurOffset)
	p.PassWd = ReadString(dataPacket.Data, &dataPacket.CurOffset, int(p.LenOfPassWd))
	p.HddHwid = ReadString(dataPacket.Data, &dataPacket.CurOffset, 16)
	p.NetCafeID = ReadUint32BE(dataPacket.Data, &dataPacket.CurOffset)
	p.Unknown02 = ReadUint32(dataPacket.Data, &dataPacket.CurOffset)
	p.UserSn = ReadUint64(dataPacket.Data, &dataPacket.CurOffset)
	p.LenOfUnknownString = ReadUint16(dataPacket.Data, &dataPacket.CurOffset)
	p.UnknownString03 = ReadString(dataPacket.Data, &dataPacket.CurOffset, int(p.LenOfUnknownString))
	p.Unknown04 = ReadUint8(dataPacket.Data, &dataPacket.CurOffset)
	p.IsLeague = ReadUint8(dataPacket.Data, &dataPacket.CurOffset)
	//p.LenOfString = ReadUint8(dataPacket.Data, &dataPacket.CurOffset)
	//p.String = ReadString(dataPacket.Data, &dataPacket.CurOffset, int(p.LenOfString))
	//...
	return true
}


func (p *PacketData) PraseInAchievementBackgroundPacket(dest *InAchievementBackgroundPacket) bool {
	//id + type = 2 bytes
	if dest == nil ||
		p.Length < 6 {
		return false
	}
	dest.ACHType = ReadUint8(p.Data, &p.CurOffset)
	dest.SelectedCase = ReadUint8(p.Data, &p.CurOffset) // value of selected background image ex: 1-2-3 v.b
	dest.Type = ReadUint8(p.Data, &p.CurOffset)
	dest.Type2 = ReadUint8(p.Data, &p.CurOffset)
	dest.Type3 = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseRoomPacket(dest *InRoomPaket) bool {
	// id + type = 2 bytes
	if p.Length < 2 ||
		dest == nil {
		return false
	}
	dest.InRoomType = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseQuickPacket(dest *InQuickPacket) bool {
	// id + type = 2 bytes
	if p.Length < 2 ||
		dest == nil {
		return false
	}
	dest.InQuickType = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseInQuickListPacket(dest *InQuickList) bool {
	// id + type + gameModID + IsEnableBot = 4 bytes
	if p.Length < 4 ||
		dest == nil {
		return false
	}
	dest.GameModID = ReadUint8(p.Data, &p.CurOffset)
	dest.IsEnableBot = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseFavoritePacket(dest *InFavoritePacket) bool {
	// id + type = 2 bytes
	if p.Length < 2 ||
		dest == nil {
		return false
	}
	dest.PacketType = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseShopPacket(dest *InShopPacket) bool {
	// id + type = 2 bytes
	if p.Length < 2 ||
		dest == nil {
		return false
	}
	dest.InShopType = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseShopBuyItemPacket(dest *InShopBuyItemPacket) bool {
	// id + type + id = 6 bytes
	if p.Length < 6 ||
		dest == nil {
		return false
	}
	dest.ItemID = ReadUint32(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseSupplyPacket(dest *InSupplyPacket) bool {
	// id + type = 2 bytes
	if p.Length < 2 ||
		dest == nil {
		return false
	}
	dest.Type = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PrasePointLottoPacket(dest *InPointLottoPacket) bool {
	// id + type = 2 bytes
	if p.Length < 2 ||
		dest == nil {
		return false
	}
	dest.Type = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseOpenBoxPacket(dest *InOpenBoxPacket) bool {
	// id + type + box + unk = 10 bytes
	if p.Length < 10 ||
		dest == nil {
		return false
	}
	dest.BoxID = ReadUint32(p.Data, &p.CurOffset)
	dest.Unk00 = ReadUint32(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseItemUsePacket(dest *InItemUsePacket) bool {
	// id + type + unk + seq + unk = 6 bytes
	if p.Length < 6 ||
		dest == nil {
		return false
	}
	dest.unk00 = ReadUint8(p.Data, &p.CurOffset)
	dest.ItemSeq = ReadUint16(p.Data, &p.CurOffset)
	dest.ItemType = ReadUint8(p.Data, &p.CurOffset)
	switch dest.ItemType {
	case 2:
		dest.lenOfString = ReadUint8(p.Data, &p.CurOffset)
		dest.String = ReadString(p.Data, &p.CurOffset, int(dest.lenOfString))
	default:
	}
	return true
}

func (p *PacketData) PraseTryItemUsePacket(dest *InTryItemUsePacket) bool {
	// id + type + unk + seq + type = 6 bytes
	if p.Length < 6 ||
		dest == nil {
		return false
	}
	dest.unk00 = ReadUint8(p.Data, &p.CurOffset)
	dest.ItemSeq = ReadUint16(p.Data, &p.CurOffset)
	dest.ItemType = ReadUint8(p.Data, &p.CurOffset)
	switch dest.ItemType {
	case 2:
		dest.lenOfNewName = ReadUint8(p.Data, &p.CurOffset)
		dest.NewName = ReadString(p.Data, &p.CurOffset, int(dest.lenOfNewName))
	default:
	}
	return true
}

func (p *PacketData) PraseFavoriteSetCosmeticsPacket(dest *InFavoriteSetCosmetics) bool {
	// id + type + slot + itemId = 7 bytes
	if p.Length < 7 {
		return false
	}
	dest.Slot = ReadUint8(p.Data, &p.CurOffset)
	dest.ItemId = ReadUint32(p.Data, &p.CurOffset)
	return true
}
func (p *PacketData) PraseReportPacket(dest *InReportPacket) bool {
	// id + type = 2 bytes
	if p.Length < 2 ||
		dest == nil {
		return false
	}
	dest.Type = ReadUint8(p.Data, &p.CurOffset)
	return true
}
func (p *PacketData) PraseReportMsgPacket(dest *InReportMsgPacket) bool {
	// id + type + type + len + type + unk = 9 bytes
	if p.Length < 9 ||
		dest == nil {
		return false
	}
	dest.LenOfName = ReadUint8(p.Data, &p.CurOffset)
	dest.Name = ReadString(p.Data, &p.CurOffset, int(dest.LenOfName))
	dest.Type = ReadUint8(p.Data, &p.CurOffset)
	dest.LenOfMsg = ReadUint16(p.Data, &p.CurOffset)
	dest.Msg = ReadString(p.Data, &p.CurOffset, int(dest.LenOfMsg))
	dest.Unk00 = ReadUint16(p.Data, &p.CurOffset)
	return true
}
func (p *PacketData) PraseReportSearchUserPacket(dest *InReportSearchUserPacket) bool {
	// id + type + type + len = 4 bytes
	if p.Length < 4 ||
		dest == nil {
		return false
	}
	dest.LenOfName = ReadUint8(p.Data, &p.CurOffset)
	dest.Name = ReadString(p.Data, &p.CurOffset, int(dest.LenOfName))
	return true
}
func (p *PacketData) PraseChannelRequest(dest *InRoomListRequestPacket) bool {
	// id + channelServerIndex + channelIndex = 3 bytes
	if p.Length < 3 {
		return false
	}
	dest.ChannelServerIndex = ReadUint8(p.Data, &p.CurOffset)
	dest.ChannelIndex = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseNewRoomQuest(dest *InNewRoomPacket) bool {
	// id + roomtype + newroompacket = 24 bytes
	if p.Length < 24 {
		return false
	}
	dest.LenOfName = ReadUint8(p.Data, &p.CurOffset)
	dest.RoomName = ReadString(p.Data, &p.CurOffset, int(dest.LenOfName))
	dest.Unk00 = ReadUint16(p.Data, &p.CurOffset)
	dest.Unk01 = ReadUint8(p.Data, &p.CurOffset)
	dest.GameModeID = ReadUint8(p.Data, &p.CurOffset)
	dest.MapID = ReadUint8(p.Data, &p.CurOffset)
	dest.WinLimit = ReadUint8(p.Data, &p.CurOffset)
	dest.KillLimit = ReadUint16(p.Data, &p.CurOffset)
	dest.Unk02 = ReadUint8(p.Data, &p.CurOffset)
	dest.Unk03 = ReadUint8(p.Data, &p.CurOffset)
	dest.Unk04 = ReadUint8(p.Data, &p.CurOffset)
	dest.LenOfUnk05 = ReadUint8(p.Data, &p.CurOffset)
	dest.Unk05 = ReadString(p.Data, &p.CurOffset, int(dest.LenOfUnk05))
	dest.Unk06 = ReadUint8(p.Data, &p.CurOffset)
	dest.Unk07 = ReadUint8(p.Data, &p.CurOffset)
	dest.Unk08 = ReadUint8(p.Data, &p.CurOffset)
	dest.Unk09 = ReadUint8(p.Data, &p.CurOffset)
	dest.Unk10 = ReadUint8(p.Data, &p.CurOffset)
	dest.Unk11 = ReadUint32(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseUpdateRoomPacket(dest *InUpSettingReq) bool {
	//读取flag，标记要读的有哪些数据
	flags := ReadUint64(p.Data, &p.CurOffset)
	dest.Flags = flags
	lowFlag := *(*uint32)(unsafe.Pointer(&flags))
	//右移32比特位
	flags = flags >> 32
	highFlag := *(*uint32)(unsafe.Pointer(&flags))
	if lowFlag&0x1 != 0 {
		dest.LenOfRoomName = ReadUint8(p.Data, &p.CurOffset)
		dest.RoomName = ReadString(p.Data, &p.CurOffset, int(dest.LenOfRoomName))
	}
	if lowFlag&0x2 != 0 {
		dest.Unk00 = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x4 != 0 {
		dest.Unk01 = ReadUint8(p.Data, &p.CurOffset)
		dest.Unk02 = ReadUint32(p.Data, &p.CurOffset)
		dest.Unk03 = ReadUint32(p.Data, &p.CurOffset)
	}
	if lowFlag&0x8 != 0 {
		dest.LenOfpasswd = ReadUint8(p.Data, &p.CurOffset)
		dest.Passwd = ReadString(p.Data, &p.CurOffset, int(dest.LenOfpasswd))
	}
	if lowFlag&0x10 != 0 {
		dest.Unk10 = ReadUint16(p.Data, &p.CurOffset)
	}
	if lowFlag&0x20 != 0 {
		dest.ForceCamera = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x40 != 0 {
		dest.GameModeID = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x80 != 0 {
		dest.MapID = ReadUint8(p.Data, &p.CurOffset)
		dest.Unk13 = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x100 != 0 {
		dest.MaxPlayers = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x200 != 0 {
		dest.WinLimit = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x400 != 0 {
		dest.KillLimit = ReadUint16(p.Data, &p.CurOffset)
	}
	if lowFlag&0x800 != 0 {
		dest.Unk17 = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x1000 != 0 {
		dest.Unk18 = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x2000 != 0 {
		dest.WeaponRestrictions = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x4000 != 0 {
		dest.Status = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x8000 != 0 {
		dest.Unk21 = ReadUint8(p.Data, &p.CurOffset)
		dest.MapCycleType = ReadUint8(p.Data, &p.CurOffset)
		dest.Unk23 = ReadUint8(p.Data, &p.CurOffset)
		dest.Unk24 = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x10000 != 0 {
		dest.Unk25 = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x20000 != 0 {
		dest.NumOfMultiMaps = ReadUint8(p.Data, &p.CurOffset)
		dest.MultiMaps = make([]uint8, dest.NumOfMultiMaps)
		for i := 0; i < int(dest.NumOfMultiMaps); i++ {
			dest.MultiMaps[i] = ReadUint8(p.Data, &p.CurOffset)
		}
	}
	if lowFlag&0x40000 != 0 {
		dest.TeamBalanceType = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x80000 != 0 {
		dest.Unk29 = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x100000 != 0 {
		dest.Unk30 = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x200000 != 0 {
		dest.Unk31 = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x400000 != 0 {
		dest.Unk32 = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x800000 != 0 {
		dest.Unk33 = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x1000000 != 0 {
		dest.BotEnabled = ReadUint8(p.Data, &p.CurOffset)
		if dest.BotEnabled != 0 {
			dest.BotDifficulty = ReadUint8(p.Data, &p.CurOffset)
			dest.NumCtBots = ReadUint8(p.Data, &p.CurOffset)
			dest.NumTrBots = ReadUint8(p.Data, &p.CurOffset)
		}
	}
	if lowFlag&0x2000000 != 0 {
		dest.Unk35 = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x4000000 != 0 {
		dest.Unk36 = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x8000000 != 0 {
		dest.Unk37 = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x10000000 != 0 {
		dest.Unk38 = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x20000000 != 0 {
		dest.Unk39 = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x40000000 != 0 {
		dest.IsIngame = ReadUint8(p.Data, &p.CurOffset)
	}
	if lowFlag&0x80000000 != 0 {
		dest.StartMoney = ReadUint16(p.Data, &p.CurOffset)
	}
	if highFlag&0x1 != 0 {
		dest.ChangeTeams = ReadUint8(p.Data, &p.CurOffset)
	}
	if highFlag&0x2 != 0 {
		dest.Unk43 = ReadUint8(p.Data, &p.CurOffset)
	}
	if highFlag&0x4 != 0 {
		dest.HltvEnabled = ReadUint8(p.Data, &p.CurOffset)
	}
	if highFlag&0x8 != 0 {
		dest.Unk45 = ReadUint8(p.Data, &p.CurOffset)
	}
	if highFlag&0x10 != 0 {
		dest.RespawnTime = ReadUint8(p.Data, &p.CurOffset)
	}
	return true
}

func (p *PacketData) PraseRoomCountdownPacket(dest *InRoomCountdownPacket) bool {
	//id + count + CountdownType + count = 4 bytes
	if p.Length < 3 ||
		dest == nil {
		return false
	}
	dest.CountdownType = ReadUint8(p.Data, &p.CurOffset)
	if dest.CountdownType == InProgress {
		if p.Length < 4 {
			return false
		}
		dest.Count = ReadUint8(p.Data, &p.CurOffset)
	}
	return true
}

func (p *PacketData) PraseRoomFeedbackPacket(dest *InFeedbackPacket) bool {
	//id + type + unk00 + code + unk01 = 10 bytes
	if p.Length < 10 ||
		dest == nil {
		return false
	}
	dest.unk00 = ReadUint16(p.Data, &p.CurOffset)
	dest.ErrorCode = ReadUint16(p.Data, &p.CurOffset)
	dest.unk01 = ReadUint32(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseJoinRoomPacket(dest *InJoinRoomPacket) bool {
	//id + join + roomId + lenOfPassWord = 5 bytes
	if p.Length < 5 {
		return false
	}
	dest.RoomId = ReadUint16(p.Data, &p.CurOffset)
	dest.LenOfPassWord = ReadUint8(p.Data, &p.CurOffset)
	dest.PassWord = ReadString(p.Data, &p.CurOffset, int(dest.LenOfPassWord))
	return true
}

func (p *PacketData) PraseChangeTeamPacket(dest *InChangeTeamPacket) bool {
	//id + change + destteam = 3 bytes
	if p.Length < 3 {
		return false
	}
	dest.NewTeam = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseFavoriteSetLoadoutPacket(dest *InFavoriteSetLoadout) bool {
	//id + loadout + Loadout +  ItemId + WeaponSlot = 8 bytes
	if p.Length < 8 {
		return false
	}
	dest.Loadout = ReadUint8(p.Data, &p.CurOffset)
	dest.WeaponSlot = ReadUint8(p.Data, &p.CurOffset)
	dest.ItemId = ReadUint32(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PrasePlayerInfoPacket(dest *InPlayerInfoPacket) bool {
	//id + type = 2 bytes
	if p.Length < 2 {
		return false
	}
	dest.InfoType = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseSetSignaturePacket(dest *InSetSignaturePacket) bool {
	//id + type + le = 3 bytes
	if p.Length < 3 {
		return false
	}
	dest.Len = ReadUint8(p.Data, &p.CurOffset)
	dest.Signature = ReadString(p.Data, &p.CurOffset, int(dest.Len))
	return true
}

func (p *PacketData) PraseSetAvatarPacket(dest *InSetAvatarPacket) bool {
	//id + type + avatar = 4 bytes
	if p.Length < 4 {
		return false
	}
	(*dest).AvatarId = ReadUint16(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseSetTitlePacket(dest *InSetTitlePacket) bool {
	//id + type + avatar = 4 bytes
	if p.Length < 4 {
		return false
	}
	dest.TitleId = ReadUint16(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseSetCampaignPacket(dest *InSetCampaignPacket) bool {
	//id + type + PacketType  = 3 bytes
	if p.Length < 3 {
		return false
	}
	dest.PacketType = ReadUint8(p.Data, &p.CurOffset)
	if dest.PacketType == 1 { //完成
		dest.CampaignId = ReadUint16(p.Data, &p.CurOffset)
	}
	return true
}

func (p *PacketData) PraseOptionPacket(dest *InOptionPacket) bool {
	//id + type = 2 bytes
	if p.Length < 2 {
		return false
	}
	dest.OptionPacketType = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseSaveBuyMenu(dest *InOptionBuyMenu) bool {
	//id + type + length + unk00 = 5 bytes
	if p.Length < 5 {
		return false
	}
	dest.MenuLength = ReadUint16(p.Data, &p.CurOffset)
	dest.Unk00 = ReadUint8(p.Data, &p.CurOffset)
	dest.Buymenu.Pistols = ReadSubMenu(p.Data, &p.CurOffset)
	dest.Buymenu.Shotguns = ReadSubMenu(p.Data, &p.CurOffset)
	dest.Buymenu.Smgs = ReadSubMenu(p.Data, &p.CurOffset)
	dest.Buymenu.Rifles = ReadSubMenu(p.Data, &p.CurOffset)
	dest.Buymenu.Snipers = ReadSubMenu(p.Data, &p.CurOffset)
	dest.Buymenu.Machineguns = ReadSubMenu(p.Data, &p.CurOffset)
	dest.Buymenu.Melees = ReadSubMenu(p.Data, &p.CurOffset)
	dest.Buymenu.Equipment = ReadSubMenu(p.Data, &p.CurOffset)
	return true
}

func ReadSubMenu(b []byte, offset *int) []uint32 {
	len := ReadUint8(b, offset)
	if len != SUBMENU_ITEM_NUM {
		log.Println("Length of submenu is illegal !")
	}
	var submenu []uint32
	for i := 0; i < SUBMENU_ITEM_NUM; i++ {
		ReadUint8(b, offset)
		submenu = append(submenu, ReadUint32(b, offset))
	}
	return submenu
}

func (p *PacketData) PraseHostPacket(dest *InHostPacket) bool {
	//id + type = 2 bytes
	if p.Length < 2 {
		return false
	}
	dest.InHostType = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseNotifyPacket(dest *InNotifyPacket) bool {
	//id + type = 2 bytes
	if p.Length < 2 {
		return false
	}
	dest.InNotifyType = ReadUint8(p.Data, &p.CurOffset)
	return true
}
func (p *PacketData) PraseNotifyListPacket(dest *InNotifyListPacket) bool {
	//id + type = 2 bytes
	if p.Length < 2 {
		return false
	}
	dest.InType = ReadUint8(p.Data, &p.CurOffset)
	return true
}
func (p *PacketData) PraseInKillPacket(dest *InKillPacket) bool {
	//id + type + ... = 16 bytes
	if p.Length < 16 ||
		dest == nil {
		return false
	}
	dest.Unk00 = ReadUint8(p.Data, &p.CurOffset)
	dest.KillerID = ReadUint32(p.Data, &p.CurOffset)
	dest.Unk01 = ReadUint32(p.Data, &p.CurOffset)
	dest.Unk02 = ReadUint8(p.Data, &p.CurOffset)
	dest.KillType = ReadUint8(p.Data, &p.CurOffset)
	dest.KillNum = ReadUint16(p.Data, &p.CurOffset)
	dest.PlayerTeam = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseInWeaponPointPacket(dest *InWeaponPointPacket) bool {
	//id + type + ... = 66 bytes
	if p.Length < 66 ||
		dest == nil {
		return false
	}
	dest.KillType = ReadUint32(p.Data, &p.CurOffset)
	dest.KillerID = ReadUint32(p.Data, &p.CurOffset)
	dest.KillerWeaponID = ReadUint32(p.Data, &p.CurOffset)
	dest.KillerTeam = ReadUint8(p.Data, &p.CurOffset)
	dest.KillerClientType = ReadUint8(p.Data, &p.CurOffset)
	dest.KillerCharacterType = ReadUint8(p.Data, &p.CurOffset)
	dest.KillerCharacterClass = ReadUint32(p.Data, &p.CurOffset)

	dest.DeadType = ReadUint32(p.Data, &p.CurOffset)
	dest.VictimID = ReadUint32(p.Data, &p.CurOffset)
	dest.VictimWeaponID = ReadUint32(p.Data, &p.CurOffset)
	dest.VictimTeam = ReadUint8(p.Data, &p.CurOffset)
	dest.VictimClientType = ReadUint8(p.Data, &p.CurOffset)
	dest.VictimCharacterType = ReadUint8(p.Data, &p.CurOffset)
	dest.VictimCharacterClass = ReadUint32(p.Data, &p.CurOffset)

	dest.KillerX = ReadUint32(p.Data, &p.CurOffset)
	dest.KillerY = ReadUint32(p.Data, &p.CurOffset)
	dest.KillerZ = ReadUint32(p.Data, &p.CurOffset)
	dest.VictimX = ReadUint32(p.Data, &p.CurOffset)
	dest.VictimY = ReadUint32(p.Data, &p.CurOffset)
	dest.VictimZ = ReadUint32(p.Data, &p.CurOffset)

	dest.ArraySize = ReadUint16(p.Data, &p.CurOffset)
	dest.Array = make([]InWeaponPointArray, dest.ArraySize)

	for i := 0; i < int(dest.ArraySize); i++ {
		var tmp InWeaponPointArray
		tmp.unk00 = ReadUint8(p.Data, &p.CurOffset)
		tmp.unk01 = ReadUint32(p.Data, &p.CurOffset)
		tmp.unk02 = ReadUint32(p.Data, &p.CurOffset)
		tmp.unk03 = ReadUint8(p.Data, &p.CurOffset)
		dest.Array = append(dest.Array, tmp)
	}
	return true
}

func (p *PacketData) PraseInRevivedPacket(dest *InRevivedPacket) bool {
	//id + type + ... = 19 bytes
	if p.Length < 19 ||
		dest == nil {
		return false
	}
	dest.UserID = ReadUint32(p.Data, &p.CurOffset)
	dest.X = ReadUint32(p.Data, &p.CurOffset)
	dest.Y = ReadUint32(p.Data, &p.CurOffset)
	dest.Z = ReadUint32(p.Data, &p.CurOffset)
	dest.Unk00 = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseInDeathPacket(dest *InDeathPacket) bool {
	//id + type + ... = 14 bytes
	if p.Length < 14 ||
		dest == nil {
		return false
	}
	dest.DeadID = ReadUint32(p.Data, &p.CurOffset)
	dest.Unk00 = ReadUint32(p.Data, &p.CurOffset)
	dest.Unk01 = ReadUint8(p.Data, &p.CurOffset)
	dest.DeathNum = ReadUint16(p.Data, &p.CurOffset)
	dest.PlayerTeam = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseInAssistPacket(dest *InAssistPacket) bool {
	//id + type + ... = 18 bytes
	if p.Length < 18 ||
		dest == nil {
		return false
	}
	dest.KillerID = ReadUint32(p.Data, &p.CurOffset)
	dest.Unk00 = ReadUint8(p.Data, &p.CurOffset)
	dest.AssisterID = ReadUint32(p.Data, &p.CurOffset)
	dest.Unk01 = ReadUint16(p.Data, &p.CurOffset)
	dest.Unk02 = ReadUint16(p.Data, &p.CurOffset)
	dest.Unk03 = ReadUint16(p.Data, &p.CurOffset)
	dest.AssisterTeam = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseSetBuyMenuPacket(dest *InHostSetBuyMenu) bool {
	//id + type + userid = 6 bytes
	if dest == nil ||
		p.Length < 6 {
		return false
	}
	dest.Userid = ReadUint32(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseInTeamChangingPacket(dest *InHostTeamChangingPacket) bool {
	//id + type + userid + team + unk = 8 bytes
	if p.Length < 8 ||
		dest == nil {
		return false
	}
	dest.UserId = ReadUint32(p.Data, &p.CurOffset)
	dest.NewTeam = ReadUint8(p.Data, &p.CurOffset)
	dest.Unk00 = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseSetUserLoadoutPacket(dest *InHostSetLoadoutPacket) bool {
	//id + type + userid = 6 bytes
	if dest == nil ||
		p.Length < 6 {
		return false
	}
	dest.UserID = ReadUint32(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseInGameScorePacket(dest *InGameScorePacket) bool {
	//id + type +... = 6 bytes
	if p.Length < 6 ||
		dest == nil {
		return false
	}
	dest.WinnerTeam = ReadUint8(p.Data, &p.CurOffset)
	dest.TrScore = ReadUint8(p.Data, &p.CurOffset)
	dest.CtScore = ReadUint8(p.Data, &p.CurOffset)
	dest.PacketType = ReadUint8(p.Data, &p.CurOffset)
	if dest.PacketType != 0 {
		dest.HostID = ReadUint32(p.Data, &p.CurOffset)
		dest.Unk00 = ReadUint32(p.Data, &p.CurOffset)
	}
	return true
}

func (p *PacketData) PraseSetUserInventoryPacket(dest *InHostSetInventoryPacket) bool {
	//id + type + userid = 6 bytes
	if dest == nil ||
		p.Length < 6 {
		return false
	}
	dest.UserID = ReadUint32(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseInChatPacket(dest *InChatPacket) bool {
	//id + type + len = 4 bytes
	if dest == nil ||
		p.Length < 5 {
		return false
	}
	dest.Type = ReadUint8(p.Data, &p.CurOffset)
	if dest.Type == ChatDirectMessage {
		dest.DestinationLen = ReadUint8(p.Data, &p.CurOffset)
		dest.Destination = ReadString(p.Data, &p.CurOffset, int(dest.DestinationLen))
	}
	dest.MessageLen = ReadUint16(p.Data, &p.CurOffset)
	dest.Message = ReadString(p.Data, &p.CurOffset, int(dest.MessageLen))
	return true
}

func (p *PacketData) PraseInAchievementPacket(dest *InAchievementPacket) bool {
	//id + type = 2 bytes
	if dest == nil ||
		p.Length < 2 {
		return false
	}
	dest.Type = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseInEventPacket(dest *InEventPacket) bool {
	//id + type = 2 bytes
	if dest == nil ||
		p.Length < 2 {
		return false
	}
	dest.Type = ReadUint8(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseDisassemblePacket(dest *InDisassemblePacket) bool {
	//id + type = 2 bytes
	if dest == nil ||
		p.Length < 2 {
		return false
	}
	dest.Type = ReadUint8(p.Data, &p.CurOffset)
	return true
}
func (p *PacketData) PraseDisassembleItemPacket(dest *InDisassembleItemPacket) bool {
	//id + type + subtype = 3 bytes
	if dest == nil ||
		p.Length < 3 {
		return false
	}
	dest.SubType = ReadUint8(p.Data, &p.CurOffset)
	return true
}
func (p *PacketData) PraseDisassembleWeaponPacket(dest *InDisassembleWeaponPacket) bool {
	//id + type + subtype + 17 = 20 bytes
	if dest == nil ||
		p.Length < 20 {
		return false
	}
	dest.Unk00 = ReadUint32(p.Data, &p.CurOffset)
	dest.Unk01 = ReadUint8(p.Data, &p.CurOffset)
	dest.ItemID = ReadUint32(p.Data, &p.CurOffset)
	dest.Unk02 = ReadUint32(p.Data, &p.CurOffset)
	dest.Unk03 = ReadUint32(p.Data, &p.CurOffset)
	return true
}
func (p *PacketData) PraseInAchievementCampaignPacket(dest *InAchievementCampaignPacket) bool {
	//id + type + campaign = 4 bytes
	if dest == nil ||
		p.Length < 4 {
		return false
	}
	dest.CampaignId = ReadUint16(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseInHostItemUsingPacket(dest *InHostItemUsingPacket) bool {
	//id + type + userid + itemid = 10 bytes
	if dest == nil ||
		p.Length < 10 {
		return false
	}
	dest.UserID = ReadUint32(p.Data, &p.CurOffset)
	dest.ItemID = ReadUint32(p.Data, &p.CurOffset)
	return true
}

func (p *PacketData) PraseInHostBuyItemPacket(dest *InHostBuyItemPacket) bool {
	//id + type + userid + itemid = 7 bytes
	if dest == nil ||
		p.Length < 7 {
		return false
	}
	dest.UserID = ReadUint32(p.Data, &p.CurOffset)
	dest.NumItemsBought = ReadUint8(p.Data, &p.CurOffset)
	for i := 0; i < int(dest.NumItemsBought); i++ {
		dest.Items = append(dest.Items, ReadUint32(p.Data, &p.CurOffset))
	}
	return true
}

func (p InRoomCountdownPacket) ShouldCountdown() bool {
	return p.CountdownType == InProgress
}

//GetNextSeq 获取下一次的seq数据包序号
func GetNextSeq(seq *uint8) uint8 {
	if *seq == MAXSEQUENCE {
		*seq = MINSEQUENCE
	} else {
		(*seq)++
	}
	return *seq
}

//BuildHeader 建立数据包通用头部
func BuildHeader(seq *uint8, id uint8) []byte {
	header := make([]byte, 5)
	header[0] = PacketTypeSignature
	header[1] = GetNextSeq(seq)
	header[2] = 0
	header[3] = 0
	header[4] = id
	return header
}

//WriteLen 写入数据长度到数据包通用头部
func WriteLen(data *[]byte) {
	headerL := uint16(len(*data)) - HeaderLen
	(*data)[2] = uint8(headerL)
	(*data)[3] = uint8(headerL >> 8)
}

//NewNullString 新建空的字符串
func NewNullString() []byte {
	return []byte{0x00, 0x00, 0x00, 0x00}
}

//SendPacket 发送数据包
func SendPacket(data []byte, client net.Conn) {
	WriteLen(&data)
	client.Write(data)
}
