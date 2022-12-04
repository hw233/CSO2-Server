package servermanager

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/database/redis"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/verbose"
	"github.com/garyburd/redigo/redis"
)

var (
	BFusername *BloomFilter
	BFgamename *BloomFilter
	DB         *sql.DB
	Dblock     sync.Mutex
	DBPath     string
	Redis      redis.Conn
)

//从数据库中读取用户数据
//如果是新用户则保存到数据库中
func GetUserFromDatabase(loginname string, passwd []byte) (*User, int) {
	u := GetNewUser()
	var bytes []byte
	if DB == nil {
		filepath := DBPath + loginname
		rst, _ := PathExists(filepath)
		if rst {

			Dblock.Lock()
			bytes, _ = ioutil.ReadFile(filepath)

			Dblock.Unlock()

		} else {
			return nil, USER_NOT_FOUND
		}
	} else {
		if !BFusername.Contains(loginname) {
			return nil, USER_NOT_FOUND
		}

		if redisIsExistsUser(loginname) {
			bytes, _ = redisGetUser(loginname)
			DebugInfo(1, "Get User", loginname+"'s data from redis")
		} else {
			sql := "SELECT * FROM UserData Where username = '" + loginname + "'"
			rows, err := DB.Query(sql)
			defer rows.Close()
			if err != nil {
				DebugInfo(1, "Searching User "+loginname+"'s data failed !")
				return nil, USER_UNKOWN_ERROR
			}

			if rows.Next() {
				rows.Scan(&u.Userid, &u.UserName, &u.IngameName, &u.Password, &u.UserMail, &bytes)
			} else {
				return nil, USER_NOT_FOUND
			}

			redisAddUser(loginname, bytes)
		}
	}
	err := json.Unmarshal(bytes, &u)
	if err != nil {
		DebugInfo(1, "Suffered a error while getting User", loginname+"'s data !", err)
		return nil, USER_UNKOWN_ERROR
	}

	//检查密码
	str := fmt.Sprintf("%x", md5.Sum([]byte(loginname+string(passwd))))
	for i := 0; i < 32; i++ {
		if str[i] != u.Password[i] {
			return nil, USER_PASSWD_ERROR
		}
	}

	u.SetID(GetNewUserID())

	DebugInfo(1, "User", u.UserName, "data found !")

	return &u, USER_LOGIN_SUCCESS
}

func AddUserToDB(u *User) error {
	if u == nil {
		return nil
	}
	data, _ := json.MarshalIndent(u, "", "     ")
	if DB == nil { //json
		filepath := DBPath + u.UserName
		Dblock.Lock()
		err := ioutil.WriteFile(filepath, data, 0644)

		Dblock.Unlock()
		if err != nil {
			return err
		}

		filepath = DBPath + u.IngameName + ".check"
		Dblock.Lock()
		err = ioutil.WriteFile(filepath, []byte(u.UserName), 0644)
		Dblock.Unlock()
		if err != nil {
			return err
		}
		return nil
	}
	//mysql
	stmt, _ := DB.Prepare(`INSERT INTO UserData (username,gamename,password,mail,data) VALUES (?, ?, ?, ?, ?)`)
	defer stmt.Close()

	_, err := stmt.Exec(u.UserName, u.IngameName, u.Password, u.UserMail, data)
	if err != nil {
		DebugInfo(1, "Insert User", u.UserName, "data failed !")
		return err
	}

	BFusername.Add(u.UserName)
	BFgamename.Add(u.IngameName)

	return nil
}

func DelOldNickNameFile(oldName string) error {
	//删除源文件
	if DB == nil {
		Dblock.Lock()
		err := os.Remove(DBPath + oldName + ".check")
		Dblock.Unlock()
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateUserToDB(u *User) error {
	if u == nil {
		return nil
	}
	data, _ := json.MarshalIndent(u, "", "     ")
	if DB == nil {
		filepath := DBPath + u.UserName
		filepathNickName := DBPath + u.IngameName + ".check"

		Dblock.Lock()
		err := ioutil.WriteFile(filepath, data, 0644)
		Dblock.Unlock()
		if err != nil {
			return err
		}

		Dblock.Lock()
		err = ioutil.WriteFile(filepathNickName, []byte(u.UserName), 0644)
		Dblock.Unlock()
		if err != nil {
			return err
		}
		return nil
	}
	stmt, _ := DB.Prepare("UPDATE UserData set gamename=?,password=?,mail=?,data=? where username=?")
	defer stmt.Close()

	_, err := stmt.Exec(u.IngameName, u.Password, u.UserMail, data, u.UserName)
	if err != nil {
		DebugInfo(1, "Update User", u.UserName, "data failed !")
		return err
	}

	redisAddUser(u.UserName, data)

	return nil
}

func IsExistsUser(username []byte) bool {
	if DB == nil {
		filepath := DBPath + string(username)
		rst, _ := PathExists(filepath)
		if rst {
			return true
		}
		return false
	}
	if !BFusername.Contains(string(username)) {
		return false
	}
	sql := "SELECT * FROM UserData Where username = '" + string(username) + "'"
	rows, err := DB.Query(sql)
	defer rows.Close()
	if err != nil {
		DebugInfo(1, "Searching User "+string(username)+"'s data failed !")
		return false
	}

	if rows.Next() {
		return true
	}
	return false

}

func IsExistsIngameName(name []byte) bool {
	if DB == nil {
		filepath := DBPath + string(name) + ".check"
		rst, _ := PathExists(filepath)
		if rst {
			return true
		}
		return false
	}
	if !BFgamename.Contains(string(name)) {
		return false
	}
	sql := "SELECT * FROM UserData Where gamename = '" + string(name) + "'"
	rows, err := DB.Query(sql)
	defer rows.Close()
	if err != nil {
		DebugInfo(1, "Searching User "+string(name)+"'s data failed !")
		return false
	}

	if rows.Next() {
		return true
	}
	return false
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func SaveAllUsers() bool {
	UsersManager.Lock.Lock()
	defer UsersManager.Lock.Unlock()
	for _, v := range UsersManager.Users {
		if v == nil {
			continue
		}

		if UpdateUserToDB(v) != nil {
			return false
		}
	}
	return true
}

func InitBloomFilter() bool {
	if DB != nil {
		fmt.Println("Initializing bloomfilter ...")
		BFusername = NewBloomFilter(2<<24, []uint{7, 11, 13, 31, 37, 61})
		BFgamename = NewBloomFilter(2<<24, []uint{7, 11, 13, 31, 37, 61})
		query, err := DB.Prepare("SELECT username,gamename FROM UserData")
		if err == nil {
			defer query.Close()
			rows, err := query.Query()
			if err != nil {
				DebugInfo(2, err)
				return false
			}
			defer rows.Close()
			var username string
			var gamename string
			for rows.Next() {
				rows.Scan(&username, &gamename)
				BFusername.Add(username)
				BFgamename.Add(gamename)
			}
		}
		return true
	}
	fmt.Println("Can't Initialize bloomfilter !")
	return false
}

func redisIsExistsUser(username string) bool {
	return RedisIsExist(Redis, "CSO2Server:Users:"+username)
}

func redisAddUser(username string, data []byte) bool {

	return RedisSetVWithTime(Redis, "CSO2Server:Users:"+username, data, "1200")
}

func redisGetUser(username string) ([]byte, error) {
	return RedisGetVBytes(Redis, "CSO2Server:Users:"+username)
}
