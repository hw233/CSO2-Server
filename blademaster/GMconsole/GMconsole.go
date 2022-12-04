package GMconsole

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"

	. "github.com/KouKouChan/CSO2-Server/blademaster/core/inventory"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/message"
	. "github.com/KouKouChan/CSO2-Server/blademaster/core/room"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/configure"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

type GMInfo struct {
	GMip       string
	GMport     string
	GMname     string
	GMpassword string
}

var (
	clients map[string]bool
)

func InitGMconsole() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Console server suffered a fault !")
			fmt.Println("error:", err)
			fmt.Println("Fault end!")
		}
	}()
	clients = map[string]bool{}

	//init TCP
	GMserver, err := net.Listen("tcp", fmt.Sprintf(":%d", Conf.GMport))
	if err != nil {
		fmt.Println("Init console tcp socket error !\n")
		panic(err)
	}
	defer GMserver.Close()

	fmt.Println("Console is running at", "[AnyAdapter]:"+strconv.Itoa(int(Conf.GMport)))
	for {
		client, err := GMserver.Accept()
		if err != nil {
			DebugInfo(2, "Console server accept data error !\n")
			continue
		}
		DebugInfo(2, "Console server accept a new connection request at", client.RemoteAddr().String())
		go RecvGMmsg(client)
	}
}

func RecvGMmsg(client net.Conn) {
	defer client.Close()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("GM Client", client.RemoteAddr().String(), "suffered a fault !")
			fmt.Println(err)
			fmt.Println("Fault end!")
			delete(clients, client.RemoteAddr().String())
		}
	}()

	for {
		//读取3字节数据包头部
		headBytes, err := GMReadHead(client)
		if !err {
			goto end
		}
		if headBytes[0] != GMSignature {
			DebugInfo(2, "Recived a illegal GM head from", client.RemoteAddr().String())
			continue
		}
		offset := 1
		Len := ReadUint16(headBytes, &offset)
		//读取数据部分
		bytes, err := GMReadData(client, Len)
		if !err {
			goto end
		}
		dataPacket := GMpacket{
			bytes,
			Len,
			string(bytes),
		}
		req := strings.Fields(string(dataPacket.Req))

		switch req[0] {
		case GMLogin:
			login(client, req)
		case GMReqUserList:
			userlist(client, req)
		case GMKickUser:
			kickUser(client, req)
		case GMadditem:
			additem(client, req)
		case GMremoveitem:
			removeitem(client, req)
		case GMdelroom:
			delroom(client, req)
		case GMsave:
			save(client, req)
		case GMBeVIP:
			vipUser(client, req)
		case GMbeGM:
			gmUser(client, req)
		default:
		}
	}
end:
	DebugInfo(1, "Console client", client.RemoteAddr().String(), "closed the connection")
	delete(clients, client.RemoteAddr().String())
	return
}

func login(client net.Conn, req []string) {
	if len(req) < 3 {
		DebugInfo(2, "Console from", client.RemoteAddr().String(), "sent a illegal login packet")
		rst := []byte(GMLoginFailed)
		GMSendPacket(&rst, client)
		return
	}
	if req[1] == Conf.GMusername && req[2] == Conf.GMpassword {
		clients[client.RemoteAddr().String()] = true
		rst := []byte(GMLoginOk)
		GMSendPacket(&rst, client)
		DebugInfo(1, "Console from", client.RemoteAddr().String(), "logged in")
	} else {
		rst := []byte(GMLoginFailed)
		GMSendPacket(&rst, client)
		DebugInfo(1, "Console from", client.RemoteAddr().String(), "login failed")
	}
}

func userlist(client net.Conn, req []string) {
	if _, ok := clients[client.RemoteAddr().String()]; !ok || !clients[client.RemoteAddr().String()] {
		DebugInfo(2, "Console from", client.RemoteAddr().String(), "sent a userlist req but not logged in")
		return
	}
	rst := OutUserList{
		len(UsersManager.Users),
		[]string{},
	}
	for _, v := range UsersManager.Users {
		if v == nil {
			rst.UserNum--
			continue
		}
		rst.UserNames = append(rst.UserNames, v.UserName)
	}
	jsondata, _ := json.Marshal(rst)
	GMSendPacket(&jsondata, client)
	DebugInfo(1, "Console from", client.RemoteAddr().String(), "request a userlist")

}

func kickUser(client net.Conn, req []string) {
	if len(req) < 2 {
		DebugInfo(2, "Console from", client.RemoteAddr().String(), "sent a illegal kick packet")
		rst := []byte(GMKickFailed)
		GMSendPacket(&rst, client)
		return
	}
	if _, ok := clients[client.RemoteAddr().String()]; !ok || !clients[client.RemoteAddr().String()] {
		DebugInfo(2, "Console from", client.RemoteAddr().String(), "sent a kick user req but not logged in")
		return
	}
	for _, v := range UsersManager.Users {
		if v == nil {
			continue
		}
		if v.UserName == req[1] {
			OnSendMessage(v.CurrentSequence, v.CurrentConnection, MessageDialogBox, GAME_SERVER_ERROR)
			OnLeaveRoom(v.CurrentConnection, true)
			DelUserWithConn(v.CurrentConnection)
			v.CurrentConnection.Close()

			rst := []byte(GMKickSuccess)
			GMSendPacket(&rst, client)
			DebugInfo(1, "Console from", client.RemoteAddr().String(), "kicked player", v.UserName)
			return
		}
	}
	rst := []byte(GMKickFailed)
	GMSendPacket(&rst, client)
}

func additem(client net.Conn, req []string) {
	if len(req) < 4 {
		DebugInfo(2, "Console from", client.RemoteAddr().String(), "sent a illegal additem packet")
		rst := []byte(GMAdditemFailed)
		GMSendPacket(&rst, client)
		return
	}
	if _, ok := clients[client.RemoteAddr().String()]; !ok || !clients[client.RemoteAddr().String()] {
		DebugInfo(2, "Console from", client.RemoteAddr().String(), "sent a additem req but not logged in")
		return
	}

	id, err := strconv.Atoi(req[2])
	if err != nil {
		rst := []byte(GMAdditemFailed)
		GMSendPacket(&rst, client)
	}

	count, err := strconv.Atoi(req[3])
	if err != nil {
		rst := []byte(GMAdditemFailed)
		GMSendPacket(&rst, client)
	}

	for _, v := range UsersManager.Users {
		if v == nil {
			continue
		}
		if v.UserName == req[1] {

			idx := v.AddItem(uint32(id), uint16(count), 0)

			rst := BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeInventory_Create),
				BuildInventoryInfoSingle(v, uint32(id), idx))
			SendPacket(rst, v.CurrentConnection)
			OnSendMessage(v.CurrentSequence, v.CurrentConnection, MessageNotice, GAME_USER_NEW_ITEM)

			rst = []byte(GMAdditemSuccess)
			GMSendPacket(&rst, client)

			DebugInfo(1, "Console from", client.RemoteAddr().String(), "add item", id, "to User", v.UserName)
			return
		}
	}
	filepath := DBPath + req[1]
	rb, _ := PathExists(filepath)
	if rb {
		u := GetNewUser()

		Dblock.Lock()
		dataEncoded, _ := ioutil.ReadFile(filepath)
		Dblock.Unlock()

		err := json.Unmarshal(dataEncoded, &u)
		if err != nil {
			rst := []byte(GMAdditemFailed)
			GMSendPacket(&rst, client)

			DebugInfo(1, "Console from", client.RemoteAddr().String(), "add item", id, "to User", u.UserName, "failed")
			return

		}

		u.AddItem(uint32(id), uint16(count), 0)
		err = UpdateUserToDB(&u)

		if err == nil {
			rst := []byte(GMAdditemSuccess)
			GMSendPacket(&rst, client)

			DebugInfo(1, "Console from", client.RemoteAddr().String(), "add item", id, "to User", u.UserName, "success")
			return

		}
	}
	rst := []byte(GMAdditemFailed)
	GMSendPacket(&rst, client)
	DebugInfo(1, "Console from", client.RemoteAddr().String(), "add item", id, "to User", req[1], "not found")
}

func removeitem(client net.Conn, req []string) {
	if len(req) < 3 {
		DebugInfo(2, "Console from", client.RemoteAddr().String(), "sent a illegal removeitem packet")
		rst := []byte(GMRemoveitemFailed)
		GMSendPacket(&rst, client)
		return
	}
	if _, ok := clients[client.RemoteAddr().String()]; !ok || !clients[client.RemoteAddr().String()] {
		DebugInfo(2, "Console from", client.RemoteAddr().String(), "sent a removeitem req but not logged in")
		return
	}

	id, err := strconv.Atoi(req[2])
	if err != nil {
		rst := []byte(GMRemoveitemFailed)
		GMSendPacket(&rst, client)
	}

	for _, v := range UsersManager.Users {
		if v == nil {
			continue
		}
		if v.UserName == req[1] {

			v.RemoveItem(uint32(id))

			rst := BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeUserInfo), BuildUserInfo(0xFFFFFFFF, NewUserInfo(v), v.Userid, true))
			SendPacket(rst, v.CurrentConnection)

			rst = []byte(GMRemoveitemSuccess)
			GMSendPacket(&rst, client)

			DebugInfo(1, "Console from", client.RemoteAddr().String(), "removed item", id, "to User", v.UserName)
			return
		}
	}
	filepath := DBPath + req[1]
	rb, _ := PathExists(filepath)
	if rb {
		u := GetNewUser()

		Dblock.Lock()
		dataEncoded, _ := ioutil.ReadFile(filepath)
		Dblock.Unlock()

		err := json.Unmarshal(dataEncoded, &u)
		if err != nil {
			rst := []byte(GMRemoveitemFailed)
			GMSendPacket(&rst, client)

			DebugInfo(1, "Console from", client.RemoteAddr().String(), "remove item", id, "to User", u.UserName, "failed")
			return

		}

		u.RemoveItem(uint32(id))
		err = UpdateUserToDB(&u)

		if err == nil {
			rst := []byte(GMRemoveitemSuccess)
			GMSendPacket(&rst, client)

			DebugInfo(1, "Console from", client.RemoteAddr().String(), "removed item", id, "to User", u.UserName, "success")
			return

		}
	}
	rst := []byte(GMRemoveitemFailed)
	GMSendPacket(&rst, client)
	DebugInfo(1, "Console from", client.RemoteAddr().String(), "remove item", id, "to User", req[1], "not found")
}

func delroom(client net.Conn, req []string) {
	if len(req) < 4 {
		DebugInfo(2, "Console from", client.RemoteAddr().String(), "sent a illegal delroom packet")
		rst := []byte(GMDelRoomFailed)
		GMSendPacket(&rst, client)
		return
	}
	if _, ok := clients[client.RemoteAddr().String()]; !ok || !clients[client.RemoteAddr().String()] {
		DebugInfo(2, "Console from", client.RemoteAddr().String(), "sent a delroom req but not logged in")
		return
	}

	chlsrvid, err := strconv.Atoi(req[1])
	if err != nil {
		rst := []byte(GMDelRoomFailed)
		GMSendPacket(&rst, client)
	}

	chlid, err := strconv.Atoi(req[2])
	if err != nil {
		rst := []byte(GMDelRoomFailed)
		GMSendPacket(&rst, client)
	}

	roomid, err := strconv.Atoi(req[3])
	if err != nil {
		rst := []byte(GMDelRoomFailed)
		GMSendPacket(&rst, client)
	}

	DelChannelRoom(uint16(roomid),
		uint8(chlid),
		uint8(chlsrvid))

	rst := []byte(GMDelRoomSuccess)
	GMSendPacket(&rst, client)
	DebugInfo(1, "Console from", client.RemoteAddr().String(), "deleted room ID=", roomid)
}

func save(client net.Conn, req []string) {
	if _, ok := clients[client.RemoteAddr().String()]; !ok || !clients[client.RemoteAddr().String()] {
		DebugInfo(2, "Console from", client.RemoteAddr().String(), "sent a additem req but not logged in")
		return
	}
	if !SaveAllUsers() {
		rst := []byte(GMSaveFailed)
		GMSendPacket(&rst, client)
		return
	}
	DebugInfo(1, "Console from", client.RemoteAddr().String(), "request to save all data")

	rst := []byte(GMSaveSuccess)
	GMSendPacket(&rst, client)
}

func vipUser(client net.Conn, req []string) {
	if len(req) < 2 {
		DebugInfo(2, "Console from", client.RemoteAddr().String(), "sent a illegal vip packet")
		rst := []byte(GMBeVIPFailed)
		GMSendPacket(&rst, client)
		return
	}
	if _, ok := clients[client.RemoteAddr().String()]; !ok || !clients[client.RemoteAddr().String()] {
		DebugInfo(2, "Console from", client.RemoteAddr().String(), "sent a vip user req but not logged in")
		return
	}
	for _, v := range UsersManager.Users {
		if v == nil {
			continue
		}
		if v.UserName == req[1] {
			v.SetVIP()

			rst := BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeUserInfo), BuildUserInfo(0xFFFFFFFF, NewUserInfo(v), v.Userid, true))
			SendPacket(rst, v.CurrentConnection)

			rst = []byte(GMBeVIPSuccess)
			GMSendPacket(&rst, client)
			DebugInfo(1, "Console from", client.RemoteAddr().String(), "set player", v.UserName, "vip success")
			return
		}
	}

	filepath := DBPath + req[1]
	rb, _ := PathExists(filepath)
	if rb {
		u := GetNewUser()

		Dblock.Lock()
		dataEncoded, _ := ioutil.ReadFile(filepath)
		Dblock.Unlock()

		err := json.Unmarshal(dataEncoded, &u)
		if err != nil {
			rst := []byte(GMBeVIPFailed)
			GMSendPacket(&rst, client)

			DebugInfo(1, "Console from", client.RemoteAddr().String(), "set player", req[1], "vip failed")
			return

		}

		u.SetVIP()

		err = UpdateUserToDB(&u)

		if err == nil {
			rst := []byte(GMBeVIPSuccess)
			GMSendPacket(&rst, client)
			DebugInfo(1, "Console from", client.RemoteAddr().String(), "set player", req[1], "vip success")
			return

		}
	}

	rst := []byte(GMBeVIPFailed)
	GMSendPacket(&rst, client)
	DebugInfo(1, "Console from", client.RemoteAddr().String(), "set player", req[1], "vip failed")
}

func gmUser(client net.Conn, req []string) {
	if len(req) < 2 {
		DebugInfo(2, "Console from", client.RemoteAddr().String(), "sent a illegal gm packet")
		rst := []byte(GMBeGMFailed)
		GMSendPacket(&rst, client)
		return
	}
	if _, ok := clients[client.RemoteAddr().String()]; !ok || !clients[client.RemoteAddr().String()] {
		DebugInfo(2, "Console from", client.RemoteAddr().String(), "sent a gm user req but not logged in")
		return
	}
	for _, v := range UsersManager.Users {
		if v == nil {
			continue
		}
		if v.UserName == req[1] {
			v.SetGM()
			rst := BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeUserInfo), BuildUserInfo(0xFFFFFFFF, NewUserInfo(v), v.Userid, true))
			SendPacket(rst, v.CurrentConnection)

			rst = []byte(GMBeGMSuccess)
			GMSendPacket(&rst, client)
			DebugInfo(1, "Console from", client.RemoteAddr().String(), "set player", v.UserName, "gm success")
			return
		}
	}

	filepath := DBPath + req[1]
	rb, _ := PathExists(filepath)
	if rb {
		u := GetNewUser()

		Dblock.Lock()
		dataEncoded, _ := ioutil.ReadFile(filepath)
		Dblock.Unlock()

		err := json.Unmarshal(dataEncoded, &u)
		if err != nil {
			rst := []byte(GMBeGMFailed)
			GMSendPacket(&rst, client)

			DebugInfo(1, "Console from", client.RemoteAddr().String(), "set player", req[1], "gm failed")
			return

		}

		u.SetGM()

		err = UpdateUserToDB(&u)

		if err == nil {
			rst := []byte(GMBeGMSuccess)
			GMSendPacket(&rst, client)
			DebugInfo(1, "Console from", client.RemoteAddr().String(), "set player", req[1], "gm success")
			return

		}
	}

	rst := []byte(GMBeGMFailed)
	GMSendPacket(&rst, client)
	DebugInfo(1, "Console from", client.RemoteAddr().String(), "set player", req[1], "gm failed")
}
