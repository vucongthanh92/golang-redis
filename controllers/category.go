package controllers

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/TIG/api-redis/database"
	"github.com/TIG/api-redis/models"
	"github.com/TIG/api-redis/services"
	"github.com/gin-gonic/gin"
)

// CreateCategory function
func CreateCategory(c *gin.Context) {
	conn := database.Connection()
	defer conn.Close()
	newCategoryID, errExis := services.CheckCategoryNewID(conn)
	var (
		status   = 200
		Category models.Category
		errors   = make([]string, 0)
	)
	if errExis != nil {
		status = 500
		errors = append(errors, errExis.Error())
	} else {
		Category = services.GetParamCategoryFromBody(c)
		Category.SID, _ = strconv.Atoi(newCategoryID)

		// insert category into redis
		errInsert := services.AddCategory(conn, Category)
		if errInsert != nil {
			status = 500
			errors = append(errors, errInsert.Error())
		}
		errIndex := services.AddNewIndexForCategory(conn, Category)
		if errIndex != nil {
			status = 500
			errors = append(errors, errIndex.Error())
		}

		if errInsert != nil {
			c.JSON(500, gin.H{
				"status":  500,
				"Error: ": errors,
			})
		} else {
			c.JSON(200, gin.H{
				"status":  status,
				"data":    Category,
				"message": "success",
			})
		}
	}
}

// GetCategories func
func GetCategories(c *gin.Context) {
	var (
		Categories = make([]*models.Category, 0)
		Err        error
	)
	conn := database.Connection()
	defer conn.Close()
	Categories, Err = services.GetAllCategories(conn)
	if Err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"error":  Err.Error(),
		})
	} else {
		Categories = services.SortCategoryBySID(Categories)
		c.JSON(200, gin.H{
			"status":  200,
			"data":    Categories,
			"message": "success",
		})
	}
}

// GetCategory func
func GetCategory(c *gin.Context) {
	var (
		Category models.Category
		Err      error
	)
	conn := database.Connection()
	CategoryID := c.Param("id")
	Category, Err = services.GetCategoryBySID(conn, CategoryID)
	if Err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"error":  Err,
		})
	} else {
		if Category.SID == 0 {
			c.JSON(403, gin.H{
				"status": 403,
				"error":  "Category not found",
			})
		} else {
			c.JSON(200, gin.H{
				"status":  200,
				"data":    Category,
				"message": "success",
			})
		}
	}
}

// UpdateCategory func
func UpdateCategory(c *gin.Context) {
	var (
		ParamsCategory models.Category
		Category       models.Category
		Err            error
	)
	conn := database.Connection()
	ParamsCategory = services.GetParamCategoryFromBody(c)
	Category, errCategory := services.GetCategoryBySID(conn, strconv.Itoa(ParamsCategory.SID))
	if errCategory != nil || Category.SID == 0 {
		c.JSON(500, gin.H{
			"status": 500,
			"error":  "Category not found",
		})
	} else {
		Category, Err = services.UpdateCategoryBySID(conn, Category, ParamsCategory)
		if Err != nil {
			c.JSON(500, gin.H{
				"status":  500,
				"message": Err,
			})
		} else {
			c.JSON(200, gin.H{
				"status":  200,
				"data":    Category,
				"message": "success",
			})
		}
	}
}

// DeleteCategory func
func DeleteCategory(c *gin.Context) {
	var (
		Category models.Category
		Err      error
	)
	conn := database.Connection()
	CategoryID := c.Param("id")
	Category, Err = services.GetCategoryBySID(conn, CategoryID)
	if Err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"error":  Err,
		})
	} else {
		if Category.SID == 0 {
			c.JSON(403, gin.H{
				"status": 403,
				"error":  "Category not found",
			})
		} else {
			errDelete := services.DeleteCategoryBySID(conn, Category)
			if errDelete != nil {
				c.JSON(500, gin.H{
					"status": 500,
					"error":  errDelete,
				})
			} else {
				c.JSON(200, gin.H{
					"status":  200,
					"message": "success",
				})
			}
		}
	}
}

// SearchCategory func
func SearchCategory(c *gin.Context) {
	var (
		params     models.SearchCategory
		Categories = make([]models.Category, 0)
		Err        error
	)
	// get params in request
	body, _ := ioutil.ReadAll(c.Request.Body)
	var jsonObject map[string]interface{}
	json.Unmarshal([]byte(string(body)), &jsonObject)
	params.ConvertParamsToSearchCategory(jsonObject)
	conn := database.Connection()

	// Retrieve data from the found result
	Categories, Err = services.SearchCategoryByField(conn, params)
	if Err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"error":  Err.Error(),
		})
	}
	if len(Categories) < 1 {
		c.JSON(403, gin.H{
			"status":  403,
			"message": "No results found",
		})
	} else {
		c.JSON(200, gin.H{
			"status":  200,
			"data":    Categories,
			"message": "success",
		})
	}
}
