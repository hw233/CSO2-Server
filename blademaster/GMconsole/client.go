package GMconsole

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

func ToConsoleHost(username, password, addr, port string) {
	fmt.Println("")
	fmt.Println("Command:")
	fmt.Println("	userlist       			  		:show how many users in server")
	fmt.Println("	kick [username]       		  		:kick user")
	fmt.Println("	save       			  		:save all online user's data to database")
	fmt.Println("	additem [username] [itemid] [count]      	:give user item")
	fmt.Println("	removeitem [username] [itemid]    		:remove a item for user")
	fmt.Println("	delroom [serverindex] [channelindex] [roomid]   :delete room by force")
	fmt.Println("	vip [username]       		  		:set user vip")
	fmt.Println("	gm [username]       		  		:set user gm")
	fmt.Println("	exit       			  		:quit the console")
	fmt.Println("")

	server, err := net.Dial("tcp", addr+":"+port)
	if err != nil {
		fmt.Println("Connect to server failed !")
		return
	}
	fmt.Println("Connect to server success !")
	defer server.Close()
	fmt.Println("Login ...")
	loginbuf := []byte(GMLogin + " " + username + " " + password)
	GMSendPacket(&loginbuf, server)

	dataPacket, rb := readReply(server)
	if !rb {
		fmt.Println("Recv login packet failed !")
		return
	}

	switch string(dataPacket.data) {
	case GMLoginOk:
		fmt.Println(GMLoginOk)
	case GMLoginFailed:
		fmt.Println(GMLoginFailed, "please be sure your have a correct username and password !")
		return
	default:
		fmt.Println("Unkown reply !")
		return
	}

	var instream string
	for {
		fmt.Printf("]")
		instream = ScanLine()
		//log.Println(instream)
		cmds := strings.Fields(instream)
		//log.Println(cmds)

		switch cmds[0] {
		case GMReqUserList:
			rst := []byte(GMReqUserList)
			GMSendPacket(&rst, server)

			dataPacket, rb = readReply(server)
			if !rb {
				fmt.Println("Recv packet failed !")
				return
			}

			var pkt OutUserList
			err = json.Unmarshal(dataPacket.data, &pkt)
			if err != nil {
				fmt.Println("Prase userlist packet failed !")
				return
			}

			for k, v := range pkt.UserNames {
				fmt.Println("["+strconv.Itoa(k+1)+"]", v)
			}
			fmt.Println("UserNum:", pkt.UserNum)

			continue
		case GMKickUser:
			rst := []byte(instream)
			GMSendPacket(&rst, server)
		case GMadditem:
			rst := []byte(instream)
			GMSendPacket(&rst, server)
		case GMremoveitem:
			rst := []byte(instream)
			GMSendPacket(&rst, server)
		case GMdelroom:
			rst := []byte(instream)
			GMSendPacket(&rst, server)
		case GMsave:
			rst := []byte(GMsave)
			GMSendPacket(&rst, server)
		case GMBeVIP:
			rst := []byte(instream)
			GMSendPacket(&rst, server)
		case GMbeGM:
			rst := []byte(instream)
			GMSendPacket(&rst, server)
		case GMexit:
			server.Close()
			return
		default:
			fmt.Println("Error command !")
			continue
		}

		dataPacket, rb = readReply(server)
		if !rb {
			fmt.Println("Recv packet failed !")
			return
		}

		fmt.Println(dataPacket.Req)

	}
}

func readReply(con net.Conn) (GMpacket, bool) {
	//读取3字节数据包头部
	headBytes, err := GMReadHead(con)
	if !err {
		fmt.Println("Recv packet failed !")
		return GMpacket{}, false
	}
	if headBytes[0] != GMSignature {
		fmt.Println("Recived a illegal head from", con.RemoteAddr().String())
		return GMpacket{}, false
	}
	offset := 1
	Len := ReadUint16(headBytes, &offset)
	//读取数据部分
	bytes, err := GMReadData(con, Len)
	if !err {
		fmt.Println("Recv packet data failed !")
		return GMpacket{}, false
	}
	dataPacket := GMpacket{
		bytes,
		Len,
		string(bytes),
	}
	return dataPacket, true
}
