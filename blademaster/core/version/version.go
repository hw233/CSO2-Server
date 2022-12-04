package version

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnVersionPacket(seq *uint8, client net.Conn) {
	header := BuildHeader(seq, PacketTypeVersion)
	header[1] = 0 //seq to 0
	*seq = 0
	IsBadHash := make([]byte, 1)
	IsBadHash[0] = 0
	hash := []byte("6246015df9a7d1f7311f888e7e861f18")
	rst := BytesCombine(header, IsBadHash, hash)
	SendPacket(rst, client)
	DebugInfo(1, "Sent a version reply to", client.RemoteAddr().String())

	// clientStr := strings.Split(client.RemoteAddr().String(), ":")[0]
	// if IsLoginTenth(clientStr) {
	// 	OnSendMessage(seq, client, MessageDialogBox, GAME_LOGIN_TENTH_FAILED)
	// 	return
	// }
	// pkt := InLoginPacket{
	// 	NexonUsername: []byte("KouKouChan"),
	// 	GameUsername:  []byte("KouKouChan"),
	// 	PassWd:        []byte("12345"),
	// }
	// nu := string(pkt.NexonUsername)
	// u, result := GetUserByLogin(nu, pkt.PassWd)
	// if result == 2 {
	// 	nu, _ = LocalToUtf8(string(pkt.NexonUsername))
	// 	u, result = GetUserByLogin(nu, pkt.PassWd)
	// }
	// switch result {
	// case USER_PASSWD_ERROR:
	// 	DebugInfo(2, "Error : User", nu, "from", client.RemoteAddr().String(), "login failed with error password !")
	// 	if IsLoginTenth(clientStr) {
	// 		OnSendMessage(seq, client, MessageDialogBox, GAME_LOGIN_TENTH_FAILED)
	// 	} else {
	// 		CountFailLogin(clientStr)
	// 		if IsLoginTenth(clientStr) {
	// 			OnSendMessage(seq, client, MessageDialogBox, GAME_LOGIN_TENTH_FAILED)
	// 			CountTenMinutes(clientStr)
	// 		} else {
	// 			OnSendMessage(seq, client, MessageDialogBox, GAME_LOGIN_BAD_PASSWORD)
	// 		}
	// 	}
	// 	return
	// case USER_ALREADY_LOGIN:
	// 	DebugInfo(2, "Error : User", nu, "from", client.RemoteAddr().String(), "already logged in !")
	// 	OnSendMessage(seq, client, MessageDialogBox, GAME_LOGIN_ALREADY)
	// 	OnSendMessage(u.CurrentSequence, u.CurrentConnection, MessageDialogBox, GAME_LOGIN_EXIT_FORCE)
	// 	u.CurrentConnection.Close()
	// 	u.QuitChannel()
	// case USER_NOT_FOUND:
	// 	DebugInfo(2, "Error : User", nu, "from", client.RemoteAddr().String(), "not registered !")
	// 	if IsLoginTenth(clientStr) {
	// 		OnSendMessage(seq, client, MessageDialogBox, GAME_LOGIN_TENTH_FAILED)
	// 	} else {
	// 		CountFailLogin(clientStr)
	// 		if IsLoginTenth(clientStr) {
	// 			DebugInfo(2, "Error : User", nu, "from", client.RemoteAddr().String(), "login failed in tenth !")
	// 			OnSendMessage(seq, client, MessageDialogBox, GAME_LOGIN_TENTH_FAILED)
	// 			CountTenMinutes(clientStr)
	// 		} else {
	// 			OnSendMessage(seq, client, MessageDialogBox, GAME_LOGIN_BAD_USERNAME)
	// 		}
	// 	}
	// 	return
	// case USER_UNKOWN_ERROR:
	// 	DebugInfo(2, "Error : User", nu, "from", client.RemoteAddr().String(), "login but suffered a error !")
	// 	OnSendMessage(seq, client, MessageDialogBox, GAME_LOGIN_ERROR)
	// 	return
	// default:
	// }

	// //检查升级
	// // if u.CheckUpdate != 1 {
	// // 	DebugInfo(2, "Updating User", nu, "data ...")
	// // 	//当前版本升级
	// // 	for _, v := range DefaultInventoryItem {
	// // 		u.AddItemSingle(v.Id, 0)
	// // 	}
	// // 	u.Updated()
	// // 	DebugInfo(2, "Updating User", nu, "data done !")
	// // }

	// //设置数据
	// ClearCount(clientStr)
	// u.CurrentConnection = client
	// u.CurrentSequence = seq
	// u.CheckOutdatedItem()
	// // illegalItems := u.CheckIllegalItem()
	// // DebugInfo(2, "Found", len(illegalItems), "illegal items for user", nu)

	// //把用户加入用户管理器
	// if !UsersManager.AddUser(u) {
	// 	DebugInfo(2, "Error : User", nu, "from", client.RemoteAddr().String(), "login failed !")
	// 	return
	// }

	// //UserStart部分
	// rst = BytesCombine(BuildHeader(u.CurrentSequence, PacketTypeUserStart), user.BuildUserStart(u))
	// SendPacket(rst, u.CurrentConnection)
	// DebugInfo(1, "User", u.UserName, "from", client.RemoteAddr().String(), "logged in !")

	// //UserInfo部分
	// rst = BytesCombine(BuildHeader(u.CurrentSequence, PacketTypeUserInfo), BuildUserInfo(0xFFFFFFFF, NewUserInfo(u), u.Userid, true))
	// SendPacket(rst, u.CurrentConnection)

	// //Inventory部分
	// rst = BytesCombine(BuildHeader(u.CurrentSequence, PacketTypeInventory_Create),
	// 	BuildInventoryInfo(u))
	// SendPacket(rst, u.CurrentConnection)

	// //unlock
	// rst = BytesCombine(BuildHeader(u.CurrentSequence, PacketTypeUnlock), BuildUnlockReply(u))
	// SendPacket(rst, u.CurrentConnection)

	// //偏好装备
	// rst = BytesCombine(BuildHeader(u.CurrentSequence, PacketTypeFavorite), BuildCosmetics(&u.Inventory))
	// SendPacket(rst, u.CurrentConnection)
	// rst = BytesCombine(BuildHeader(u.CurrentSequence, PacketTypeFavorite), BuildLoadout(&u.Inventory))
	// SendPacket(rst, u.CurrentConnection)

	// //购买菜单
	// rst = BytesCombine(BuildHeader(u.CurrentSequence, PacketTypeOption), BuildBuyMenu(&u.Inventory))
	// SendPacket(rst, u.CurrentConnection)

	// //achievement

	// //friends

	// //ServerList部分
	// OnServerList(u.CurrentConnection)

	// //motd
	// OnSendMessage(u.CurrentSequence, u.CurrentConnection, MessageNotice, Locales.MOTD)
	// OnSendMessage(u.CurrentSequence, u.CurrentConnection, MessageNotice, []byte("本游戏免费，如果你是买到的本游戏，说明你被骗了！官网：github.com/KouKouChan/CSO2-Server"))
}
