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

// CheckProductNewID func
func CheckProductNewID(conn redis.Conn) (string, error) {
	productExis, errExist := redis.Int(conn.Do("exists", "product:new:id"))
	if errExist != nil {
		return "", errExist
	}
	if productExis == 0 {
		errSetProduct := conn.Send("set", "product:new:id", 1)
		if errSetProduct != nil {
			return "", errSetProduct
		}
		productSID, errGetProduct := redis.String(conn.Do("get", "product:new:id"))
		if errGetProduct != nil {
			return "", errGetProduct
		}
		return productSID, nil
	}
	productSID, _ := redis.Int(conn.Do("incr", "product:new:id"))
	return strconv.Itoa(productSID), nil
}

// GetParamProductFromBody func
func GetParamProductFromBody(c *gin.Context) models.Product {
	var Product models.Product
	body, _ := ioutil.ReadAll(c.Request.Body)
	var jsonObject map[string]interface{}
	json.Unmarshal([]byte(string(body)), &jsonObject)
	Product.ConvertParamsToProduct(jsonObject)
	return Product
}

// AddProduct func
func AddProduct(conn redis.Conn, Product models.Product) error {
	_, errInsert := conn.Do("hmset", "product:object:"+strconv.Itoa(Product.SID),
		"sid", Product.SID,
		"title", Product.Title,
		"description", Product.Description,
		"categoryid", Product.CategoryID,
		"images", Product.Images,
		"price", Product.Price,
		"quantity", Product.Quantity,
		"status", Product.Status,
		"createdby", Product.Createdby,
	)
	if errInsert != nil {
		return errInsert
	}
	return nil
}

// GetAllProducts func
func GetAllProducts(conn redis.Conn) ([]*models.Product, error) {
	listproducts, err := redis.Strings(conn.Do("keys", "product:object:*"))
	if err != nil {
		return nil, err
	}
	conn.Send("multi")
	for _, productKey := range listproducts {
		errGetProduct := conn.Send("hgetall", productKey)
		if errGetProduct != nil {
			conn.Send("discard")
			return nil, errGetProduct
		}
	}
	arrProduct, errExec := redis.Values(conn.Do("exec"))
	if errExec != nil {
		return nil, errExec
	}
	var Products = make([]*models.Product, 0)

	for _, item := range arrProduct {
		var Product models.Product
		errProduct := redis.ScanStruct(item.([]interface{}), &Product)
		if errProduct != nil {
			return nil, errProduct
		}
		Products = append(Products, &Product)
	}
	return Products, nil
}

// GetProductBySID func
func GetProductBySID(conn redis.Conn, ProductID string) (models.Product, error) {
	var Product models.Product
	result, errGetProduct := redis.Values(conn.Do("hgetall", "product:object:"+ProductID))
	if errGetProduct != nil {
		return Product, errGetProduct
	}
	errProduct := redis.ScanStruct(result, &Product)
	if errProduct != nil {
		return Product, errProduct
	}
	return Product, nil
}

// UpdateProductBySID func
func UpdateProductBySID(conn redis.Conn, currentProduct models.Product, paramsProduct models.Product) (models.Product, error) {
	if currentProduct.Title != paramsProduct.Title && paramsProduct.Title != "" {
		_, ErrUpdate := conn.Do(
			"hset",
			"product:object:"+strconv.Itoa(currentProduct.SID),
			"title", paramsProduct.Title)
		if ErrUpdate == nil {
			UpdateIndexForProduct(
				conn, "title",
				currentProduct.Title,
				paramsProduct.Title,
				strconv.Itoa(currentProduct.SID),
			)
		}
	}
	if currentProduct.Description != paramsProduct.Description {
		conn.Send(
			"hset",
			"product:object:"+strconv.Itoa(currentProduct.SID),
			"description", paramsProduct.Description)
	}
	if currentProduct.CategoryID != paramsProduct.CategoryID {
		_, ErrUpdate := conn.Do(
			"hset",
			"product:object:"+strconv.Itoa(currentProduct.SID),
			"categoryid", paramsProduct.CategoryID)
		if ErrUpdate == nil {
			UpdateIndexForProduct(
				conn, "categoryid",
				strconv.Itoa(currentProduct.CategoryID),
				strconv.Itoa(paramsProduct.CategoryID),
				strconv.Itoa(currentProduct.SID),
			)
		}
	}
	if currentProduct.Images != paramsProduct.Images {
		conn.Send(
			"hset",
			"product:object:"+strconv.Itoa(currentProduct.SID),
			"images", paramsProduct.Images)
	}
	if currentProduct.Price != paramsProduct.Price {
		ErrUpdate := conn.Send(
			"hset",
			"product:object:"+strconv.Itoa(currentProduct.SID),
			"price", paramsProduct.Price)
		if ErrUpdate == nil {
			UpdateIndexForProduct(
				conn, "price",
				strconv.Itoa(currentProduct.Price),
				strconv.Itoa(paramsProduct.Price),
				strconv.Itoa(currentProduct.SID),
			)
		}
	}
	if currentProduct.Quantity != paramsProduct.Quantity {
		ErrUpdate := conn.Send(
			"hset",
			"product:object:"+strconv.Itoa(currentProduct.SID),
			"quantity", paramsProduct.Quantity)
		if ErrUpdate == nil {
			UpdateIndexForProduct(
				conn, "quantity",
				strconv.Itoa(currentProduct.Quantity),
				strconv.Itoa(paramsProduct.Quantity),
				strconv.Itoa(currentProduct.SID),
			)
		}
	}
	if currentProduct.Status != paramsProduct.Status {
		ErrUpdate := conn.Send(
			"hset",
			"product:object:"+strconv.Itoa(currentProduct.SID),
			"status", paramsProduct.Status)
		if ErrUpdate == nil {
			UpdateIndexForProduct(
				conn, "status",
				strconv.Itoa(currentProduct.Status),
				strconv.Itoa(paramsProduct.Status),
				strconv.Itoa(currentProduct.SID),
			)
		}
	}
	if currentProduct.Createdby != paramsProduct.Createdby {
		_, ErrUpdate := conn.Do(
			"hset",
			"product:object:"+strconv.Itoa(currentProduct.SID),
			"createdby", paramsProduct.Createdby)
		if ErrUpdate == nil {
			UpdateIndexForProduct(
				conn, "createdby",
				strconv.Itoa(currentProduct.Createdby),
				strconv.Itoa(paramsProduct.Createdby),
				strconv.Itoa(currentProduct.SID),
			)
		}
	}
	currentProduct, _ = GetProductBySID(conn, strconv.Itoa(currentProduct.SID))
	return currentProduct, nil
}

// DeleteProductBySID func
func DeleteProductBySID(conn redis.Conn, Product models.Product) error {
	errDelete := conn.Send("del", "product:object:"+strconv.Itoa(Product.SID))
	if errDelete != nil {
		return errDelete
	}
	DeleteIndexForProduct(conn, "title", Product.Title, strconv.Itoa(Product.SID))
	DeleteIndexForProduct(conn, "categoryid", strconv.Itoa(Product.CategoryID), strconv.Itoa(Product.SID))
	DeleteIndexForProduct(conn, "price", strconv.Itoa(Product.Price), strconv.Itoa(Product.SID))
	DeleteIndexForProduct(conn, "quantity", strconv.Itoa(Product.Quantity), strconv.Itoa(Product.SID))
	DeleteIndexForProduct(conn, "status", strconv.Itoa(Product.Status), strconv.Itoa(Product.SID))
	DeleteIndexForProduct(conn, "createdby", strconv.Itoa(Product.Createdby), strconv.Itoa(Product.SID))
	return nil
}

//======================function process index for Product===================================

// GetCurrentIndexForProduct func
func GetCurrentIndexForProduct(conn redis.Conn, key string, field string) string {
	productExis, errExist := redis.Int(conn.Do("hexists", key, field))
	if errExist != nil {
		return ""
	}
	if productExis == 0 {
		return ""
	}
	valueIndex, _ := redis.String(conn.Do("hget", key, field))
	return valueIndex
}

// SortProductBySID func
func SortProductBySID(Products []*models.Product) []*models.Product {
	sort.Slice(Products, func(i, j int) bool {
		switch strings.Compare(strconv.Itoa(Products[i].SID), strconv.Itoa(Products[j].SID)) {
		case -1:
			return true
		case 1:
			return false
		}
		return Products[i].SID > Products[j].SID
	})
	return Products
}

// CreateIndexForProduct func
func CreateIndexForProduct(conn redis.Conn, key string, field string, value string) error {
	productExis := GetCurrentIndexForProduct(conn, key, field)
	if productExis == "" {
		_, errSetIndex := conn.Do("hset", key, field, value)
		if errSetIndex != nil {
			return errSetIndex
		}
		return nil
	}
	result := productExis + "," + value
	_, errSetIndex := conn.Do("hset", key, field, result)
	if errSetIndex != nil {
		return errSetIndex
	}
	return nil
}

// AddNewIndexForProduct func
func AddNewIndexForProduct(conn redis.Conn, Product models.Product) error {
	errTitle := CreateIndexForProduct(conn, "product:index:title", Product.Title, strconv.Itoa(Product.SID))
	if errTitle != nil {
		conn.Do("discard")
		return errTitle
	}
	errCategorySID := CreateIndexForProduct(conn, "product:index:categoryid", strconv.Itoa(Product.CategoryID), strconv.Itoa(Product.SID))
	if errCategorySID != nil {
		conn.Do("discard")
		return errCategorySID
	}
	errPrice := CreateIndexForProduct(conn, "product:index:price", strconv.Itoa(Product.Price), strconv.Itoa(Product.SID))
	if errPrice != nil {
		conn.Do("discard")
		return errPrice
	}
	errQuantity := CreateIndexForProduct(conn, "product:index:quantity", strconv.Itoa(Product.Quantity), strconv.Itoa(Product.SID))
	if errQuantity != nil {
		conn.Do("discard")
		return errQuantity
	}
	errStatus := CreateIndexForProduct(conn, "product:index:status", strconv.Itoa(Product.Status), strconv.Itoa(Product.SID))
	if errStatus != nil {
		conn.Do("discard")
		return errStatus
	}
	errCreatedby := CreateIndexForProduct(conn, "product:index:createdby", strconv.Itoa(Product.Createdby), strconv.Itoa(Product.SID))
	if errCreatedby != nil {
		conn.Do("discard")
		return errCreatedby
	}
	return nil
}

// UpdateIndexForProduct func
func UpdateIndexForProduct(conn redis.Conn, key string, currentField string, newField string, sid string) error {
	currentIndex := GetCurrentIndexForProduct(conn, "product:index:"+key, currentField)
	newIndex := GetCurrentIndexForProduct(conn, "product:index:"+key, newField)
	// check current index product
	if currentIndex != "" {
		arrValueIndex := strings.Split(currentIndex, ",")
		if len(arrValueIndex) == 1 {
			conn.Send("hdel", "product:index:"+key, currentField)
		} else {
			posItem := libs.SearchItemInArray(arrValueIndex, sid)
			if posItem > -1 {
				arrNewValueIndex := libs.RemoveItemInArray(arrValueIndex, posItem)
				strValueIndex := strings.Join(arrNewValueIndex, ",")
				conn.Send("hset", "product:index:"+key, currentField, strValueIndex)
			}
		}
	}
	// check new index product
	if newIndex == "" {
		errSetIndex := conn.Send("hset", "product:index:"+key, newField, sid)
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
	errSetIndex := conn.Send("hset", "product:index:"+key, newField, strValueIndex)
	if errSetIndex != nil {
		return errSetIndex
	}
	return nil
}

// DeleteIndexForProduct func
func DeleteIndexForProduct(conn redis.Conn, key string, field string, sid string) string {
	currentIndex := GetCurrentIndexForProduct(conn, "product:index:"+key, field)
	if currentIndex == "" {
		return "Index not found !"
	}
	arrValueIndex := strings.Split(currentIndex, ",")
	if len(arrValueIndex) == 1 {
		conn.Do("hdel", "product:index:"+key, field)
	} else {
		posItem := libs.SearchItemInArray(arrValueIndex, sid)
		arrNewValueIndex := libs.RemoveItemInArray(arrValueIndex, posItem)
		strValueIndex := strings.Join(arrNewValueIndex, ",")
		conn.Send("hset", "product:index:"+key, field, strValueIndex)
	}
	return ""
}
