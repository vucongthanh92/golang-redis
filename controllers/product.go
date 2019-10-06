package controllers

import (
	"strconv"

	"github.com/TIG/api-redis/database"
	"github.com/TIG/api-redis/models"
	"github.com/TIG/api-redis/services"
	"github.com/gin-gonic/gin"
)

// CreateProduct function
func CreateProduct(c *gin.Context) {
	conn := database.Connection()
	defer conn.Close()
	newProductSID, errExis := services.CheckProductNewID(conn)
	var (
		status  = 200
		Product models.Product
		errors  = make([]string, 0)
	)
	if errExis != nil {
		status = 500
		errors = append(errors, errExis.Error())
	} else {
		Product = services.GetParamProductFromBody(c)
		Product.SID, _ = strconv.Atoi(newProductSID)

		// insert product into redis
		errInsert := services.AddProduct(conn, Product)
		if errInsert != nil {
			status = 500
			errors = append(errors, errInsert.Error())
		}
		errIndex := services.AddNewIndexForProduct(conn, Product)
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
				"data":    Product,
				"message": "success",
			})
		}
	}
}

// GetProducts func
func GetProducts(c *gin.Context) {
	var (
		Products = make([]*models.Product, 0)
		Err      error
	)
	conn := database.Connection()
	defer conn.Close()
	Products, Err = services.GetAllProducts(conn)
	if Err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"error":  Err.Error(),
		})
	} else {
		Products = services.SortProductBySID(Products)
		c.JSON(200, gin.H{
			"status":  200,
			"data":    Products,
			"message": "success",
		})
	}
}

// GetProduct func
func GetProduct(c *gin.Context) {
	var (
		Product models.Product
		Err     error
	)
	conn := database.Connection()
	ProductID := c.Param("id")
	Product, Err = services.GetProductBySID(conn, ProductID)
	if Err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"error":  Err,
		})
	} else {
		if Product.SID == 0 {
			c.JSON(403, gin.H{
				"status": 403,
				"error":  "Product not found",
			})
		} else {
			c.JSON(200, gin.H{
				"status":  200,
				"data":    Product,
				"message": "success",
			})
		}
	}
}

// UpdateProduct func
func UpdateProduct(c *gin.Context) {
	var (
		ParamsProduct models.Product
		Product       models.Product
		Err           error
	)
	conn := database.Connection()
	ParamsProduct = services.GetParamProductFromBody(c)
	Product, errProduct := services.GetProductBySID(conn, strconv.Itoa(ParamsProduct.SID))
	if errProduct != nil || Product.SID == 0 {
		c.JSON(500, gin.H{
			"status": 500,
			"error":  "Product not found",
		})
	} else {
		Product, Err = services.UpdateProductBySID(conn, Product, ParamsProduct)
		if Err != nil {
			c.JSON(500, gin.H{
				"status":  500,
				"message": Err,
			})
		} else {
			c.JSON(200, gin.H{
				"status":  200,
				"data":    Product,
				"message": "success",
			})
		}
	}
}

// DeleteProduct func
func DeleteProduct(c *gin.Context) {
	var (
		Product models.Product
		Err     error
	)
	conn := database.Connection()
	ProductID := c.Param("id")
	Product, Err = services.GetProductBySID(conn, ProductID)
	if Err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"error":  Err,
		})
	} else {
		if Product.SID == 0 {
			c.JSON(403, gin.H{
				"status": 403,
				"error":  "Product not found",
			})
		} else {
			errDelete := services.DeleteProductBySID(conn, Product)
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
