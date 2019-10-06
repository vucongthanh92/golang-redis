package models

import (
	"strconv"

	libs "github.com/TIG/api-redis/helpers"
)

// Product struct
type Product struct {
	SID         int    `redis:"sid"`
	Title       string `redis:"title"`
	Description string `redis:"description"`
	CategoryID  int    `redis:"categoryid"`
	Images      string `redis:"images"`
	Price       int    `redis:"price"`
	Quantity    int    `redis:"quantity"`
	Status      int    `redis:"status"`
	Createdby   int    `redis:"createdby"`
}

// ConvertParamsToProduct func
func (product *Product) ConvertParamsToProduct(params map[string]interface{}) {
	var (
		res   interface{}
		value string
	)
	if params["sid"] != nil {
		value, res = libs.PassValueFromJSONToObject("sid", params)
		if res != nil {
			product.SID, _ = strconv.Atoi(value)
		}
	}
	if params["title"] != nil {
		value, res = libs.PassValueFromJSONToObject("title", params)
		if res != nil {
			product.Title = value
		}
	}
	if params["description"] != nil {
		value, res = libs.PassValueFromJSONToObject("description", params)
		if res != nil {
			product.Description = value
		}
	}
	if params["categoryid"] != nil {
		value, res = libs.PassValueFromJSONToObject("categoryid", params)
		if res != nil {
			product.CategoryID, _ = strconv.Atoi(value)
		}
	}
	if params["images"] != nil {
		value, res = libs.PassValueFromJSONToObject("images", params)
		if res != nil {
			product.Images = value
		}
	}
	if params["price"] != nil {
		value, res = libs.PassValueFromJSONToObject("price", params)
		if res != nil {
			product.Price, _ = strconv.Atoi(value)
		}
	}
	if params["quantity"] != nil {
		value, res = libs.PassValueFromJSONToObject("quantity", params)
		if res != nil {
			product.Quantity, _ = strconv.Atoi(value)
		}
	}
	if params["status"] != nil {
		value, res = libs.PassValueFromJSONToObject("status", params)
		if res != nil {
			product.Status, _ = strconv.Atoi(value)
		}
	}
	if params["createdby"] != nil {
		value, res = libs.PassValueFromJSONToObject("createdby", params)
		if res != nil {
			product.Createdby, _ = strconv.Atoi(value)
		}
	}
	return
}
