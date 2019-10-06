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

// CheckCategoryNewID func
func CheckCategoryNewID(conn redis.Conn) (string, error) {
	categoryExis, errExist := redis.Int(conn.Do("exists", "category:new:id"))
	if errExist != nil {
		return "", errExist
	}
	if categoryExis == 0 {
		errSetCategory := conn.Send("set", "category:new:id", 1)
		if errSetCategory != nil {
			return "", errSetCategory
		}
		categoryID, errGetCategory := redis.String(conn.Do("get", "category:new:id"))
		if errGetCategory != nil {
			return "", errGetCategory
		}
		return categoryID, nil
	}
	categoryID, _ := redis.Int(conn.Do("incr", "category:new:id"))
	return strconv.Itoa(categoryID), nil
}

// GetParamCategoryFromBody func
func GetParamCategoryFromBody(c *gin.Context) models.Category {
	var Category models.Category
	body, _ := ioutil.ReadAll(c.Request.Body)
	var jsonObject map[string]interface{}
	json.Unmarshal([]byte(string(body)), &jsonObject)
	Category.ConvertParamsToCategory(jsonObject)
	return Category
}

// AddCategory func
func AddCategory(conn redis.Conn, Category models.Category) error {
	_, errInsert := conn.Do("hmset", "category:object:"+strconv.Itoa(Category.SID),
		"sid", Category.SID,
		"title", Category.Title,
		"description", Category.Description,
		"parentid", Category.ParentID,
		"status", Category.Status,
		"createdby", Category.Createdby,
	)
	if errInsert != nil {
		return errInsert
	}
	return nil
}

// SortCategoryBySID func
func SortCategoryBySID(Categories []*models.Category) []*models.Category {
	sort.Slice(Categories, func(i, j int) bool {
		switch strings.Compare(strconv.Itoa(Categories[i].SID), strconv.Itoa(Categories[j].SID)) {
		case -1:
			return true
		case 1:
			return false
		}
		return Categories[i].SID > Categories[j].SID
	})
	return Categories
}

// GetAllCategories func
func GetAllCategories(conn redis.Conn) ([]*models.Category, error) {
	listCategories, err := redis.Strings(conn.Do("keys", "category:object:*"))
	if err != nil {
		return nil, err
	}
	conn.Send("multi")
	for _, categoryKey := range listCategories {
		errGetCategory := conn.Send("hgetall", categoryKey)
		if errGetCategory != nil {
			conn.Send("discard")
			return nil, errGetCategory
		}
	}
	arrCategory, errExec := redis.Values(conn.Do("exec"))
	if errExec != nil {
		return nil, errExec
	}
	var Categories = make([]*models.Category, 0)

	for _, item := range arrCategory {
		var Category models.Category
		errCategory := redis.ScanStruct(item.([]interface{}), &Category)
		if errCategory != nil {
			return nil, errCategory
		}
		Categories = append(Categories, &Category)
	}
	return Categories, nil
}

// GetCategoryBySID func
func GetCategoryBySID(conn redis.Conn, CategoryID string) (models.Category, error) {
	var Category models.Category
	result, errGetCategory := redis.Values(conn.Do("hgetall", "category:object:"+CategoryID))
	if errGetCategory != nil {
		return Category, errGetCategory
	}
	errCategory := redis.ScanStruct(result, &Category)
	if errCategory != nil {
		return Category, errCategory
	}
	return Category, nil
}

// UpdateCategoryBySID func
func UpdateCategoryBySID(conn redis.Conn, currentCategory models.Category, paramsCategory models.Category) (models.Category, error) {
	if currentCategory.Title != paramsCategory.Title && paramsCategory.Title != "" {
		_, ErrUpdate := conn.Do(
			"hset",
			"category:object:"+strconv.Itoa(currentCategory.SID),
			"title", paramsCategory.Title)
		if ErrUpdate == nil {
			UpdateIndexForCategory(
				conn, "title",
				currentCategory.Title,
				paramsCategory.Title,
				strconv.Itoa(currentCategory.SID),
			)
		}
	}
	if currentCategory.Description != paramsCategory.Description {
		conn.Send(
			"hset",
			"category:object:"+strconv.Itoa(currentCategory.SID),
			"description", paramsCategory.Description)
	}
	if currentCategory.ParentID != paramsCategory.ParentID {
		_, ErrUpdate := conn.Do(
			"hset",
			"category:object:"+strconv.Itoa(currentCategory.SID),
			"parentid", paramsCategory.ParentID)
		if ErrUpdate == nil {
			UpdateIndexForCategory(
				conn, "parentid",
				strconv.Itoa(currentCategory.ParentID),
				strconv.Itoa(paramsCategory.ParentID),
				strconv.Itoa(currentCategory.SID),
			)
		}
	}
	if currentCategory.Status != paramsCategory.Status {
		ErrUpdate := conn.Send(
			"hset",
			"category:object:"+strconv.Itoa(currentCategory.SID),
			"status", paramsCategory.Status)
		if ErrUpdate == nil {
			UpdateIndexForCategory(
				conn, "status",
				strconv.Itoa(currentCategory.Status),
				strconv.Itoa(paramsCategory.Status),
				strconv.Itoa(currentCategory.SID),
			)
		}
	}
	if currentCategory.Createdby != paramsCategory.Createdby {
		_, ErrUpdate := conn.Do(
			"hset",
			"category:object:"+strconv.Itoa(currentCategory.SID),
			"createdby", paramsCategory.Createdby)
		if ErrUpdate == nil {
			UpdateIndexForCategory(
				conn, "createdby",
				strconv.Itoa(currentCategory.Createdby),
				strconv.Itoa(paramsCategory.Createdby),
				strconv.Itoa(currentCategory.SID),
			)
		}
	}
	currentCategory, _ = GetCategoryBySID(conn, strconv.Itoa(currentCategory.SID))
	return currentCategory, nil
}

// DeleteCategoryBySID func
func DeleteCategoryBySID(conn redis.Conn, Category models.Category) error {
	errDelete := conn.Send("del", "category:object:"+strconv.Itoa(Category.SID))
	if errDelete != nil {
		return errDelete
	}
	DeleteIndexForCategory(conn, "title", Category.Title, strconv.Itoa(Category.SID))
	DeleteIndexForCategory(conn, "parentid", strconv.Itoa(Category.ParentID), strconv.Itoa(Category.SID))
	DeleteIndexForCategory(conn, "status", strconv.Itoa(Category.Status), strconv.Itoa(Category.SID))
	DeleteIndexForCategory(conn, "createdby", strconv.Itoa(Category.Createdby), strconv.Itoa(Category.SID))
	return nil
}

// SearchCategoryByField func
func SearchCategoryByField(conn redis.Conn, params models.SearchCategory) ([]models.Category, error) {
	var (
		Categories []models.Category
		arrIndex   map[string]string
		err        error
	)
	result, errSearch := redis.Values(conn.Do("hscan", "category:index:"+params.Keyword, "0", "match", "*"+params.Value+"*"))
	if errSearch != nil {
		return nil, errSearch
	}
	arrIndex, errArr := redis.StringMap(result[1], err)
	if errArr != nil {
		return nil, errArr
	}
	// Loop each item found
	for _, strItem := range arrIndex {
		arrSID := strings.Split(strItem, ",")
		if len(arrSID) < 1 {
			continue
		} else if len(arrSID) == 1 {
			// if length array is one item
			var Category models.Category
			responseCategory, errCategory := redis.Values(conn.Do("hgetall", "category:object:"+arrSID[0]))
			if errCategory != nil {
				continue
			}
			errCategory = redis.ScanStruct(responseCategory, &Category)
			if errCategory != nil {
				continue
			}
			Categories = append(Categories, Category)
		} else {
			// if length array is many items
			for _, sid := range arrSID {
				var Category models.Category
				responseCategory, errCategory := redis.Values(conn.Do("hgetall", "category:object:"+sid))
				if errCategory != nil {
					continue
				}
				errCategory = redis.ScanStruct(responseCategory, &Category)
				if errCategory != nil {
					continue
				}
				Categories = append(Categories, Category)
			}
		}
	}
	return Categories, nil
}

//======================function process index for Category===================================

// GetCurrentIndexForCategory func
func GetCurrentIndexForCategory(conn redis.Conn, key string, field string) string {
	categoryExis, errExist := redis.Int(conn.Do("hexists", key, field))
	if errExist != nil {
		return ""
	}
	if categoryExis == 0 {
		return ""
	}
	valueIndex, _ := redis.String(conn.Do("hget", key, field))
	return valueIndex
}

// CreateIndexForCategory func
func CreateIndexForCategory(conn redis.Conn, key string, field string, value string) error {
	categoryExis := GetCurrentIndexForCategory(conn, key, field)
	if categoryExis == "" {
		_, errSetIndex := conn.Do("hset", key, field, value)
		if errSetIndex != nil {
			return errSetIndex
		}
		return nil
	}
	result := categoryExis + "," + value
	_, errSetIndex := conn.Do("hset", key, field, result)
	if errSetIndex != nil {
		return errSetIndex
	}
	return nil
}

// AddNewIndexForCategory func
func AddNewIndexForCategory(conn redis.Conn, Category models.Category) error {
	errTitle := CreateIndexForCategory(conn, "category:index:title", Category.Title, strconv.Itoa(Category.SID))
	if errTitle != nil {
		conn.Do("discard")
		return errTitle
	}
	errParentID := CreateIndexForCategory(conn, "category:index:parentid", strconv.Itoa(Category.ParentID), strconv.Itoa(Category.SID))
	if errParentID != nil {
		conn.Do("discard")
		return errParentID
	}
	errStatus := CreateIndexForCategory(conn, "category:index:status", strconv.Itoa(Category.Status), strconv.Itoa(Category.SID))
	if errStatus != nil {
		conn.Do("discard")
		return errStatus
	}
	errCreatedby := CreateIndexForCategory(conn, "category:index:createdby", strconv.Itoa(Category.Createdby), strconv.Itoa(Category.SID))
	if errCreatedby != nil {
		conn.Do("discard")
		return errCreatedby
	}
	return nil
}

// UpdateIndexForCategory func
func UpdateIndexForCategory(conn redis.Conn, key string, currentField string, newField string, sid string) error {
	currentIndex := GetCurrentIndexForCategory(conn, "category:index:"+key, currentField)
	newIndex := GetCurrentIndexForCategory(conn, "category:index:"+key, newField)
	// check current index category
	if currentIndex != "" {
		arrValueIndex := strings.Split(currentIndex, ",")
		if len(arrValueIndex) == 1 {
			conn.Send("hdel", "category:index:"+key, currentField)
		} else {
			posItem := libs.SearchItemInArray(arrValueIndex, sid)
			if posItem > -1 {
				arrNewValueIndex := libs.RemoveItemInArray(arrValueIndex, posItem)
				strValueIndex := strings.Join(arrNewValueIndex, ",")
				conn.Send("hset", "category:index:"+key, currentField, strValueIndex)
			}
		}
	}
	// check new index category
	if newIndex == "" {
		errSetIndex := conn.Send("hset", "category:index:"+key, newField, sid)
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
	errSetIndex := conn.Send("hset", "category:index:"+key, newField, strValueIndex)
	if errSetIndex != nil {
		return errSetIndex
	}
	return nil
}

// DeleteIndexForCategory func
func DeleteIndexForCategory(conn redis.Conn, key string, field string, sid string) string {
	currentIndex := GetCurrentIndexForCategory(conn, "category:index:"+key, field)
	if currentIndex == "" {
		return "Index not found !"
	}
	arrValueIndex := strings.Split(currentIndex, ",")
	if len(arrValueIndex) == 1 {
		conn.Do("hdel", "category:index:"+key, field)
	} else {
		posItem := libs.SearchItemInArray(arrValueIndex, sid)
		arrNewValueIndex := libs.RemoveItemInArray(arrValueIndex, posItem)
		strValueIndex := strings.Join(arrNewValueIndex, ",")
		conn.Send("hset", "category:index:"+key, field, strValueIndex)
	}
	return ""
}
