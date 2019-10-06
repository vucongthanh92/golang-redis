package controllers

import (
	"strconv"

	"github.com/TIG/api-redis/database"
	libs "github.com/TIG/api-redis/helpers"
	"github.com/TIG/api-redis/models"
	"github.com/TIG/api-redis/services"
	"github.com/gin-gonic/gin"
)

// CreateUser function
func CreateUser(c *gin.Context) {
	conn := database.Connection()
	defer conn.Close()
	newUserID, errExis := services.CheckUserNewID(conn)
	var (
		status = 200
		User   models.User
		errors = make([]string, 0)
	)
	if errExis != nil {
		status = 500
		errors = append(errors, errExis.Error())
	} else {
		User = services.GetParamFromBody(c)
		User.SID, _ = strconv.Atoi(newUserID)

		// insert user into redis
		errInsert := services.AddUser(conn, User)
		if errInsert != nil {
			status = 500
			errors = append(errors, errInsert.Error())
		}
		errIndex := services.AddNewIndexForUser(conn, User)
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
				"data":    User,
				"message": "success",
			})
		}
	}
}

// GetUsers func
func GetUsers(c *gin.Context) {
	libs.CheckTableUserRole()
	var (
		Users = make([]*models.User, 0)
		Err   error
	)
	conn := database.Connection()
	defer conn.Close()
	Users, Err = services.GetAllUsers(conn)
	if Err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"error":  Err.Error(),
		})
	} else {
		Users = services.SortUserBySID(Users)
		c.JSON(200, gin.H{
			"status":  200,
			"data":    Users,
			"message": "success",
		})
	}
}

// GetUser func
func GetUser(c *gin.Context) {
	var (
		User models.User
		Err  error
	)
	conn := database.Connection()
	UserID := c.Param("id")
	User, Err = services.GetUserBySID(conn, UserID)
	if Err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"error":  Err,
		})
	} else {
		if User.SID == 0 {
			c.JSON(500, gin.H{
				"status": 403,
				"error":  "User not found",
			})
		} else {
			c.JSON(200, gin.H{
				"status": 200,
				"data":   User,
			})
		}
	}
}

// UpdateUser func
func UpdateUser(c *gin.Context) {
	var (
		ParamsUser models.User
		User       models.User
		Err        error
	)
	conn := database.Connection()
	ParamsUser = services.GetParamFromBody(c)
	User, errUser := services.GetUserBySID(conn, strconv.Itoa(ParamsUser.SID))
	if errUser != nil || User.SID == 0 {
		c.JSON(500, gin.H{
			"status": 500,
			"error":  "User not found",
		})
	} else {
		User, Err = services.UpdateUserDatabase(conn, User, ParamsUser)
		if Err != nil {
			c.JSON(500, gin.H{
				"status":  500,
				"message": Err,
			})
		} else {
			c.JSON(200, gin.H{
				"status":  200,
				"data":    User,
				"message": "success",
			})
		}
	}
}

// DeleteUser func
func DeleteUser(c *gin.Context) {
	var (
		User models.User
		Err  error
	)
	conn := database.Connection()
	UserID := c.Param("id")
	User, Err = services.GetUserBySID(conn, UserID)
	if Err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"error":  Err,
		})
	} else {
		if User.SID == 0 {
			c.JSON(403, gin.H{
				"status": 403,
				"error":  "User not found",
			})
		} else {
			errDelete := services.DeleteUserBySID(conn, User)
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
