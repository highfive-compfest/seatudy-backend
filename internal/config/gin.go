package config

import "github.com/gin-gonic/gin"

func NewGin() *gin.Engine {
	return gin.Default()
}
