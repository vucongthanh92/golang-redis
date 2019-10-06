package routers

import (
	"github.com/TIG/api-redis/controllers"
	"github.com/gin-gonic/gin"
)

// UserRoute route
func UserRoute(r *gin.RouterGroup) {
	r.GET("/users", controllers.GetUsers)
	r.GET("/user/:id", controllers.GetUser)
	r.POST("/user", controllers.CreateUser)
	r.PUT("/user", controllers.UpdateUser)
	r.DELETE("/user/:id", controllers.DeleteUser)
}

// CategoryRoute route
func CategoryRoute(r *gin.RouterGroup) {
	r.GET("/categories", controllers.GetCategories)
	r.GET("/category/:id", controllers.GetCategory)
	r.POST("/category", controllers.CreateCategory)
	r.POST("/search/category/", controllers.SearchCategory)
	r.PUT("/category", controllers.UpdateCategory)
	r.DELETE("/category/:id", controllers.DeleteCategory)
}

// ProductRoute route
func ProductRoute(r *gin.RouterGroup) {
	r.GET("/products", controllers.GetProducts)
	r.GET("/product/:id", controllers.GetProduct)
	r.POST("/product", controllers.CreateProduct)
	r.PUT("/product", controllers.UpdateProduct)
	r.DELETE("/product/:id", controllers.DeleteProduct)
}
