package models

import (
	"strconv"

	libs "github.com/TIG/api-redis/helpers"
)

// Category struct
type Category struct {
	SID         int    `redis:"sid"`
	Title       string `redis:"title"`
	Description string `redis:"description"`
	ParentID    int    `redis:"parentid"`
	Status      int    `redis:"status"`
	Createdby   int    `redis:"createdby"`
}

// SearchCategory struct
type SearchCategory struct {
	Keyword string
	Value   string
}

// ConvertParamsToCategory func
func (category *Category) ConvertParamsToCategory(params map[string]interface{}) {
	var (
		res   interface{}
		value string
	)
	if params["sid"] != nil {
		value, res = libs.PassValueFromJSONToObject("sid", params)
		if res != nil {
			category.SID, _ = strconv.Atoi(value)
		}
	}
	if params["title"] != nil {
		value, res = libs.PassValueFromJSONToObject("title", params)
		if res != nil {
			category.Title = value
		}
	}
	if params["description"] != nil {
		value, res = libs.PassValueFromJSONToObject("description", params)
		if res != nil {
			category.Description = value
		}
	}
	if params["parentid"] != nil {
		value, res = libs.PassValueFromJSONToObject("parentid", params)
		if res != nil {
			category.ParentID, _ = strconv.Atoi(value)
		}
	}
	if params["status"] != nil {
		value, res = libs.PassValueFromJSONToObject("status", params)
		if res != nil {
			category.Status, _ = strconv.Atoi(value)
		}
	}
	if params["createdby"] != nil {
		value, res = libs.PassValueFromJSONToObject("createdby", params)
		if res != nil {
			category.Createdby, _ = strconv.Atoi(value)
		}
	}
	return
}

// ConvertParamsToSearchCategory func
func (search *SearchCategory) ConvertParamsToSearchCategory(params map[string]interface{}) {
	var (
		res   interface{}
		value string
	)
	value, res = libs.PassValueFromJSONToObject("keyword", params)
	if res != nil {
		search.Keyword = value
	}
	value, res = libs.PassValueFromJSONToObject("value", params)
	if res != nil {
		search.Value = value
	}
	return
}
