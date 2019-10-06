package models

import (
	"strconv"

	libs "github.com/TIG/api-redis/helpers"
)

// User struct
type User struct {
	SID      int    `redis:"sid"`
	UserName string `redis:"username"`
	Passord  string `redis:"password"`
	Fullname string `redis:"fullname"`
	RoleID   int    `redis:"role"`
	RoleName string
}

// ConvertParamsToUser func
func (user *User) ConvertParamsToUser(params map[string]interface{}) {
	var (
		res   interface{}
		value string
	)
	value, res = libs.PassValueFromJSONToObject("sid", params)
	if res != nil {
		user.SID, _ = strconv.Atoi(value)
	}
	value, res = libs.PassValueFromJSONToObject("username", params)
	if res != nil {
		user.UserName = value
	}
	value, res = libs.PassValueFromJSONToObject("password", params)
	if res != nil {
		user.Passord = value
	}
	value, res = libs.PassValueFromJSONToObject("fullname", params)
	if res != nil {
		user.Fullname = value
	}
	value, res = libs.PassValueFromJSONToObject("role", params)
	if res != nil {
		user.RoleID, _ = strconv.Atoi(value)
	}
	return
}
