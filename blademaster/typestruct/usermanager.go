package typestruct

import (
	"sync"

	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

//UManager 全局用户管理
type UManager struct {
	UserNum int
	Users   map[uint32]*User
	Lock    sync.Mutex
}

var (
	usersManagerLock sync.Mutex
	//UsersManager 全局用户管理
	UsersManager = UManager{
		0,
		map[uint32]*User{},
		usersManagerLock,
	}
	MaxUsers uint32 = 0
)

const (
	//MAXUSERNUM , 8096 is enough
	MAXUSERNUM = 8096
)

func (dest *UManager) AddUser(src *User) bool {
	if dest == nil || src == nil {
		return false
	}
	if src.Userid <= 0 {
		DebugInfo(2, "Usermanager Error : ID of User", (*src).UserName, "is ", src.Userid)
		return false
	}
	if dest.UserNum > MAXUSERNUM {
		DebugInfo(2, "Usermanager Error : Online users is too more to add user !")
		return false
	}
	if _, ok := dest.Users[src.Userid]; ok {
		DebugInfo(2, "Usermanager Error : User is already in !")
		return false
	}
	dest.Lock.Lock()
	defer dest.Lock.Unlock()
	dest.UserNum++
	dest.Users[src.Userid] = src
	return true
}

func (dest *UManager) DelUser(src *User) bool {
	if dest == nil {
		return false
	}
	if src.Userid <= 0 {
		DebugInfo(2, "Usermanager Error : ID of User", (*src).UserName, "is illegal !")
		return false
	}
	if dest.UserNum <= 0 {
		DebugInfo(2, "Usermanager Error : There is no online user !")
		return false
	}
	if _, ok := dest.Users[src.Userid]; ok {
		dest.Lock.Lock()
		defer dest.Lock.Unlock()
		delete(dest.Users, src.Userid)
		dest.UserNum--
		return true
	}
	return false
}

func (dest *UManager) GetChannelUsers(chlsrvid, chlid uint8) []byte {
	if dest.UserNum <= 0 {
		DebugInfo(2, "Usermanager Error : There is no online user !")
		return []byte{0, 0}
	}
	buf := make([]byte, 2)
	offset, num := 0, 0
	dest.Lock.Lock()
	defer dest.Lock.Unlock()
	for _, u := range dest.Users {
		if u.GetUserChannelServerID() == chlsrvid &&
			u.GetUserChannelID() == chlid {
			num++

			tmpbuf := make([]byte, 128)
			WriteUint32(&tmpbuf, u.Userid, &offset)
			WriteString(&tmpbuf, []byte(u.IngameName), &offset)

			buf = BytesCombine(buf, tmpbuf[:offset], BuildUserInfo(0xFFFFFFFF, NewUserInfo(u), 0, false))

			offset = 0
		}
	}
	WriteUint16(&buf, uint16(num), &offset)
	return buf[:offset]
}
