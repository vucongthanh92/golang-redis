package main

import (
	"os"

	libs "github.com/TIG/api-redis/helpers"
	"github.com/TIG/api-redis/routers"
	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}

func getPort() string {
	p := os.Getenv("HOST_PORT")
	if p != "" {
		return ":" + p
	}
	return ":3030"
}

// CORSMiddleware func
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	port := getPort()
	gin.SetMode(os.Getenv("GIN_MODE"))
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.Use(location.Default())
	// check table role for user
	libs.CheckTableUserRole()
	rg := r.Group("/api")
	rg.Use(CORSMiddleware())
	{
		routers.UserRoute(rg)
		routers.CategoryRoute(rg)
		routers.ProductRoute(rg)
	}
	r.Run(port)
}
