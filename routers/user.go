package routers

import (
	"github.com/TIG/api-redis/controllers"
	"github.com/gin-gonic/gin"
)

// UserRoute function
func UserRoute(r *gin.RouterGroup) {
	r.GET("/user", controllers.GetAllUsers)
	r.GET("/user/:id", controllers.GetUserByID)
	r.POST("/user", controllers.CreateUser)
	r.PUT("/user", controllers.UpdateUser)
	r.DELETE("/user/:id", controllers.DeleteUser)
}
