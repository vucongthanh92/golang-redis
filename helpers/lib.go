package libs

import (
	"fmt"
	"strings"

	"github.com/TIG/api-redis/database"
	"github.com/gomodule/redigo/redis"
)

// CheckTableUserRole func
func CheckTableUserRole() {
	conn := database.Connection()
	defer conn.Close()
	roleExis, errExis := redis.Int(conn.Do("exists", "user:table:role"))
	if errExis != nil {
		fmt.Println("Create table role error: ", errExis)
	} else {
		if roleExis == 0 {
			conn.Send("hmset", "user:table:role",
				"1", "admin",
				"2", "manager",
				"3", "user",
			)
		}
	}
}

// PassValueFromJSONToObject func
func PassValueFromJSONToObject(FieldName string, jsonObject map[string]interface{}) (string, interface{}) {
	valuePostFromObject := jsonObject[FieldName]
	valuePostFrom := fmt.Sprintf("%v", valuePostFromObject)
	if valuePostFromObject == nil || valuePostFrom == "" {
		valuePostFromObjectLower := jsonObject[strings.ToLower(FieldName)]
		valuePostLower := fmt.Sprintf("%v", valuePostFromObjectLower)
		return valuePostLower, valuePostFromObjectLower
	}
	return valuePostFrom, valuePostFromObject
}

// RemoveItemInArray func
func RemoveItemInArray(arr []string, index int) []string {
	return append(arr[:index], arr[index+1:]...)
}

// SearchItemInArray func
func SearchItemInArray(arr []string, strQuery string) int {
	var position = -1
	for index, value := range arr {
		if strQuery == value {
			position = index
		}
	}
	return position
}
