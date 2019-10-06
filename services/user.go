package services

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	libs "github.com/TIG/api-redis/helpers"
	"github.com/TIG/api-redis/models"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

// AddUser func
func AddUser(conn redis.Conn, User models.User) error {
	_, errInsert := conn.Do("hmset", "user:object:"+strconv.Itoa(User.SID),
		"sid", User.SID,
		"username", User.UserName,
		"password", User.Passord,
		"fullname", User.Fullname,
		"role", User.RoleID,
	)
	if errInsert != nil {
		return errInsert
	}
	return nil
}

// CheckUserNewID func
func CheckUserNewID(conn redis.Conn) (string, error) {
	userExis, errExist := redis.Int(conn.Do("exists", "user:new:id"))
	if errExist != nil {
		return "", errExist
	}
	if userExis == 0 {
		_, errSetUser := conn.Do("set", "user:new:id", 1)
		if errSetUser != nil {
			return "", errSetUser
		}
		userID, errGetUser := redis.String(conn.Do("get", "user:new:id"))
		if errGetUser != nil {
			return "", errGetUser
		}
		return userID, nil
	}
	userID, _ := redis.Int(conn.Do("incr", "user:new:id"))
	return strconv.Itoa(userID), nil
}

// SortUserBySID func
func SortUserBySID(Users []*models.User) []*models.User {
	sort.Slice(Users, func(i, j int) bool {
		switch strings.Compare(strconv.Itoa(Users[i].SID), strconv.Itoa(Users[j].SID)) {
		case -1:
			return true
		case 1:
			return false
		}
		return Users[i].SID > Users[j].SID
	})
	return Users
}

// GetAllUsers func
func GetAllUsers(conn redis.Conn) ([]*models.User, error) {
	listUsers, err := redis.Strings(conn.Do("keys", "user:object:*"))
	if err != nil {
		return nil, err
	}
	conn.Send("multi")
	for _, userKey := range listUsers {
		errGetUser := conn.Send("hgetall", userKey)
		if errGetUser != nil {
			conn.Send("discard")
			return nil, errGetUser
		}
	}
	arrUsers, errExec := redis.Values(conn.Do("exec"))
	if errExec != nil {
		return nil, errExec
	}
	var Users = make([]*models.User, 0)

	// get role name for user and append list
	for _, item := range arrUsers {
		var User models.User
		errUser := redis.ScanStruct(item.([]interface{}), &User)
		if errUser != nil {
			return nil, errUser
		}
		roleName, errRoleName := redis.String(conn.Do("hget", "user:table:role", User.RoleID))
		if errRoleName == nil {
			User.RoleName = roleName
		}
		Users = append(Users, &User)
	}
	return Users, nil
}

// GetUserBySID func
func GetUserBySID(conn redis.Conn, UserID string) (models.User, error) {
	var User models.User
	result, errGetUser := redis.Values(conn.Do("hgetall", "user:object:"+UserID))
	if errGetUser != nil {
		return User, errGetUser
	}
	errUser := redis.ScanStruct(result, &User)
	if errUser != nil {
		return User, errUser
	}
	User.RoleName, _ = redis.String(conn.Do("hget", "user:table:role", User.RoleID))
	return User, nil
}

// GetParamFromBody func
func GetParamFromBody(c *gin.Context) models.User {
	var User models.User
	body, _ := ioutil.ReadAll(c.Request.Body)
	var jsonObject map[string]interface{}
	json.Unmarshal([]byte(string(body)), &jsonObject)
	User.ConvertParamsToUser(jsonObject)
	return User
}

// UpdateUserDatabase func
func UpdateUserDatabase(conn redis.Conn, currentUser models.User, ParamsUser models.User) (models.User, error) {
	if currentUser.Fullname != ParamsUser.Fullname && ParamsUser.Fullname != "" {
		_, ErrUpdate := conn.Do(
			"hset",
			"user:object:"+strconv.Itoa(currentUser.SID),
			"fullname", ParamsUser.Fullname)
		if ErrUpdate == nil {
			UpdateIndexForUser(
				conn, "fullname",
				currentUser.Fullname,
				ParamsUser.Fullname,
				strconv.Itoa(currentUser.SID),
			)
		}
	}
	if currentUser.RoleID != ParamsUser.RoleID {
		_, ErrUpdate := conn.Do(
			"hset",
			"user:object:"+strconv.Itoa(currentUser.SID),
			"role", ParamsUser.RoleID)
		if ErrUpdate == nil {
			UpdateIndexForUser(
				conn, "role",
				strconv.Itoa(currentUser.RoleID),
				strconv.Itoa(ParamsUser.RoleID),
				strconv.Itoa(currentUser.SID),
			)
		}
	}
	currentUser, _ = GetUserBySID(conn, strconv.Itoa(currentUser.SID))
	return currentUser, nil
}

// DeleteUserBySID func
func DeleteUserBySID(conn redis.Conn, User models.User) error {
	errDelete := conn.Send("del", "user:object:"+strconv.Itoa(User.SID))
	if errDelete != nil {
		return errDelete
	}
	DeleteIndexForUser(conn, "username", User.UserName, strconv.Itoa(User.SID))
	DeleteIndexForUser(conn, "password", User.Passord, strconv.Itoa(User.SID))
	DeleteIndexForUser(conn, "fullname", User.Fullname, strconv.Itoa(User.SID))
	DeleteIndexForUser(conn, "role", strconv.Itoa(User.RoleID), strconv.Itoa(User.SID))
	return nil
}

//======================function process index for User===================================

// GetCurrentIndexForUser func
func GetCurrentIndexForUser(conn redis.Conn, key string, field string) string {
	userExis, errExist := redis.Int(conn.Do("hexists", key, field))
	if errExist != nil {
		return ""
	}
	if userExis == 0 {
		return ""
	}
	valueIndex, _ := redis.String(conn.Do("hget", key, field))
	return valueIndex
}

// CreateIndexForUser func
func CreateIndexForUser(conn redis.Conn, key string, field string, value string) error {
	userExis := GetCurrentIndexForUser(conn, key, field)
	if userExis == "" {
		_, errSetIndex := conn.Do("hset", key, field, value)
		if errSetIndex != nil {
			return errSetIndex
		}
		return nil
	}
	result := userExis + "," + value
	_, errSetIndex := conn.Do("hset", key, field, result)
	if errSetIndex != nil {
		return errSetIndex
	}
	return nil
}

// AddNewIndexForUser func
func AddNewIndexForUser(conn redis.Conn, User models.User) error {
	errUsername := CreateIndexForUser(conn, "user:index:username", User.UserName, strconv.Itoa(User.SID))
	if errUsername != nil {
		conn.Do("discard")
		return errUsername
	}
	errPassword := CreateIndexForUser(conn, "user:index:password", User.Passord, strconv.Itoa(User.SID))
	if errPassword != nil {
		conn.Do("discard")
		return errPassword
	}
	errFullname := CreateIndexForUser(conn, "user:index:fullname", User.Fullname, strconv.Itoa(User.SID))
	if errFullname != nil {
		conn.Do("discard")
		return errFullname
	}
	errRole := CreateIndexForUser(conn, "user:index:role", strconv.Itoa(User.RoleID), strconv.Itoa(User.SID))
	if errRole != nil {
		conn.Do("discard")
		return errRole
	}
	return nil
}

// UpdateIndexForUser func
func UpdateIndexForUser(conn redis.Conn, key string, currentField string, newField string, sid string) error {
	currentIndex := GetCurrentIndexForUser(conn, "user:index:"+key, currentField)
	newIndex := GetCurrentIndexForUser(conn, "user:index:"+key, newField)
	// check current index User
	if currentIndex != "" {
		arrValueIndex := strings.Split(currentIndex, ",")
		if len(arrValueIndex) == 1 {
			conn.Send("hdel", "user:index:"+key, currentField)
		} else {
			posItem := libs.SearchItemInArray(arrValueIndex, sid)
			if posItem > -1 {
				arrNewValueIndex := libs.RemoveItemInArray(arrValueIndex, posItem)
				strValueIndex := strings.Join(arrNewValueIndex, ",")
				conn.Send("hset", "user:index:"+key, currentField, strValueIndex)
			}
		}
	}
	// check new index User
	if newIndex == "" {
		errSetIndex := conn.Send("hset", "user:index:"+key, newField, sid)
		if errSetIndex != nil {
			return errSetIndex
		}
		return nil
	}
	arrValueIndex := strings.Split(newIndex, ",")
	posItem := libs.SearchItemInArray(arrValueIndex, sid)
	if posItem > -1 {
		return nil
	}
	strValueIndex := newIndex + "," + sid
	errSetIndex := conn.Send("hset", "user:index:"+key, newField, strValueIndex)
	if errSetIndex != nil {
		return errSetIndex
	}
	return nil
}

// DeleteIndexForUser func
func DeleteIndexForUser(conn redis.Conn, key string, field string, sid string) string {
	currentIndex := GetCurrentIndexForUser(conn, "user:index:"+key, field)
	if currentIndex == "" {
		return "Index not found !"
	}
	arrValueIndex := strings.Split(currentIndex, ",")
	if len(arrValueIndex) == 1 {
		conn.Do("hdel", "user:index:"+key, field)
	} else {
		posItem := libs.SearchItemInArray(arrValueIndex, sid)
		arrNewValueIndex := libs.RemoveItemInArray(arrValueIndex, posItem)
		strValueIndex := strings.Join(arrNewValueIndex, ",")
		conn.Send("hset", "user:index:"+key, field, strValueIndex)
	}
	return ""
}
