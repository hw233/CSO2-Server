package GMconsole

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
)

type GMpacket struct {
	data []byte
	Len  uint16
	Req  string
}

type GMlogin struct {
	GMname     string
	GMpassword string
}

type OutUserList struct {
	UserNum   int
	UserNames []string
}

const (
	GMSignature = 0x66
	GMHeaderLen = 3

	GMReqUserList       = "userlist"
	GMLogin             = "login"
	GMKickUser          = "kick"
	GMadditem           = "additem"
	GMremoveitem        = "removeitem"
	GMdelroom           = "delroom"
	GMsave              = "save"
	GMBeVIP             = "vip"
	GMbeGM              = "gm"
	GMexit              = "exit"
	GMLoginOk           = "LoginSuccess"
	GMLoginFailed       = "LoginFailed"
	GMKickSuccess       = "KickSuccess"
	GMKickFailed        = "KickFailed"
	GMAdditemFailed     = "AddFailed"
	GMAdditemSuccess    = "AddSuccess"
	GMRemoveitemFailed  = "RemoveFailed"
	GMRemoveitemSuccess = "RemoveSuccess"
	GMDelRoomFailed     = "DelRoomFailed"
	GMDelRoomSuccess    = "DelRoomSuccess"
	GMSaveSuccess       = "SaveSuccess"
	GMSaveFailed        = "SaveFailed"
	GMBeVIPSuccess      = "VIPSuccess"
	GMBeVIPFailed       = "VIPFailed"
	GMBeGMSuccess       = "GMSuccess"
	GMBeGMFailed        = "GMFailed"
)

func GMReadHead(client net.Conn) ([]byte, bool) {
	head, curlen := make([]byte, GMHeaderLen), 0
	for {
		n, err := client.Read(head[curlen:])
		if err != nil {
			return head, false
		}
		curlen += n
		if curlen >= GMHeaderLen {
			break
		}
	}
	return head, true
}

func GMReadData(client net.Conn, len uint16) ([]byte, bool) {
	data, curlen := make([]byte, len), 0
	for {
		n, err := client.Read(data[curlen:])
		if err != nil {
			return data, false
		}
		curlen += n
		if curlen >= int(len) {
			break
		}
	}
	return data, true
}
func GMSendPacket(data *[]byte, client net.Conn) {
	head := make([]byte, 3)

	head[0] = GMSignature

	headerL := uint16(len(*data))
	head[1] = uint8(headerL)
	head[2] = uint8(headerL >> 8)

	client.Write(BytesCombine(head, *data))
}
