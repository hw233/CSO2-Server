package servermanager

import (
	"net"
	"sync"
	"time"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

const (
	USER_LOGIN_SUCCESS = 0
	USER_ALREADY_LOGIN = 1
	USER_NOT_FOUND     = 2
	USER_PASSWD_ERROR  = 3
	USER_UNKOWN_ERROR  = 4
)

var (
	LoginCounter     = map[string]int{}
	LoginCounterLock sync.Mutex
)

func DelUserWithConn(con net.Conn) bool {
	if UsersManager.UserNum == 0 {
		DebugInfo(2, "UsersManager Error : There is no online user !")
		return false
	}
	for k, v := range UsersManager.Users {
		if v.CurrentConnection == con {
			v.CheckOutdatedItem()
			CheckErr(UpdateUserToDB(v))
			UsersManager.Lock.Lock()
			defer UsersManager.Lock.Unlock()
			delete(UsersManager.Users, k)
			UsersManager.UserNum--
			return true
		}
	}
	return false
}

func GetNewUserID() uint32 {
	if UsersManager.UserNum > MAXUSERNUM {
		DebugInfo(2, "Online users is too much , unable to get a new id !")
		//ID=0 是非法的
		return 0
	}
	//如果map中不存在该ID，则返回ID
	UsersManager.Lock.Lock()
	defer UsersManager.Lock.Unlock()
	for i := 1; i <= MAXUSERNUM; i++ {
		if _, ok := UsersManager.Users[uint32(i)]; !ok {
			return uint32(i)
		}
	}
	return 0
}

//getUserByLogin 假定nexonUsername是唯一
func GetUserByLogin(account string, passwd []byte) (*User, int) {
	//查看是否有已经登陆的同名用户
	for _, v := range UsersManager.Users {
		if v.UserName == account {
			return v, USER_ALREADY_LOGIN
		}
	}
	//查看数据库是否有该用户
	return GetUserFromDatabase(account, passwd)
}

//通过连接获取用户
func GetUserFromConnection(client net.Conn) *User {
	if UsersManager.UserNum <= 0 {
		return nil
	}
	for _, v := range UsersManager.Users {
		if v.CurrentConnection == client {
			return v
		}
	}
	return nil
}

//通过ID获取用户
func GetUserFromID(id uint32) *User {
	if UsersManager.UserNum <= 0 {
		return nil
	}
	if v, ok := UsersManager.Users[id]; ok {
		return v
	}
	return nil
}

//通过name获取用户
func GetUserFromIngameName(name []byte) *User {
	if UsersManager.UserNum <= 0 {
		return nil
	}
	for _, v := range UsersManager.Users {
		if CompareBytes([]byte(v.IngameName), name) {
			return v
		}
	}
	return nil
}

func CountFailLogin(client string) {
	LoginCounterLock.Lock()
	defer LoginCounterLock.Unlock()
	if _, ok := LoginCounter[client]; ok {
		LoginCounter[client]++
	} else {
		LoginCounter[client] = 1
	}
}

func IsLoginTenth(client string) bool {
	LoginCounterLock.Lock()
	defer LoginCounterLock.Unlock()
	if _, ok := LoginCounter[client]; ok && LoginCounter[client] >= 10 {
		return true
	}
	return false
}

func CountTenMinutes(client string) {
	timer := time.NewTimer(10 * time.Minute)
	<-timer.C

	LoginCounterLock.Lock()
	defer LoginCounterLock.Unlock()

	delete(LoginCounter, client)
}

func ClearCount(client string) {
	LoginCounterLock.Lock()
	defer LoginCounterLock.Unlock()

	delete(LoginCounter, client)
}
