package chat

import (
	"net"
	"strconv"
	"strings"

	. "github.com/KouKouChan/CSO2-Server/blademaster/core/message"
	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

func OnChatTeamMessage(p *InChatPacket, client net.Conn) {
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent TeamMessage but not in server !")
		return
	}
	if uPtr.GetUserRoomID() <= 0 {
		DebugInfo(2, "Error : User", string(uPtr.IngameName), "sent TeamMessage but not in room !")
		return
	}
	//找到对应房间
	rm := GetRoomFromID(uPtr.GetUserChannelServerID(), uPtr.GetUserChannelID(), uPtr.GetUserRoomID())
	if rm == nil || rm.Id <= 0 {
		DebugInfo(2, "Error : User", string(uPtr.IngameName), "sent TeamMessage but not in room !")
		return
	}
	if !uPtr.CurrentIsIngame {
		DebugInfo(2, "Error : User", string(uPtr.IngameName), "sent TeamMessage but not ingame !")
		return
	}
	//发送数据
	strs := strings.Fields(string(p.Message))
	if len(strs) >= 3 && strs[2] == "/users" {
		idx := 0
		rm.RoomMutex.Lock()
		defer rm.RoomMutex.Unlock()
		for _, v := range rm.Users {
			if v == nil {
				continue
			}
			if v.CurrentIsIngame {
				var rst []byte
				if rm.HostUserID == v.Userid {
					rst = BytesCombine([]byte("[Host] "), []byte(v.UserName))
				} else {
					rst = BytesCombine([]byte("["+strconv.Itoa(idx)+"] "), []byte(v.UserName))
				}
				idx++
				OnSendMessage(uPtr.CurrentSequence, uPtr.CurrentConnection, MessageNotice, rst)
			}
		}
		return
	}

	msg := BuildChatMessage(uPtr, p, ChatIngameTeam)
	rm.RoomMutex.Lock()
	defer rm.RoomMutex.Unlock()
	for _, v := range rm.Users {
		if v.CurrentIsIngame && v.GetUserTeam() == uPtr.GetUserTeam() {
			SendPacket(BytesCombine(BuildHeader(v.CurrentSequence, PacketTypeChat), msg), v.CurrentConnection)
		}
	}
	DebugInfo(1, "User", string(uPtr.IngameName), "say team message <", string(p.Message), "> in room id", rm.Id)
}
