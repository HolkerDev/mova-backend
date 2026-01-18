package router

import (
	"mova-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	r.GET("/", handler.HelloWorld)

	return r
}
