package auth

import (
	"github.com/gin-gonic/gin"
)

type RestController struct {
	uc *UseCase
}

func NewRestController(router *gin.Engine, uc *UseCase) *RestController {
	controller := &RestController{uc: uc}

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", controller.Register())
		authGroup.POST("/login", controller.Login())
		authGroup.POST("/refresh", controller.Refresh())
		authGroup.POST("/verify-email", controller.VerifyEmail())
	}

	return controller
}

func (c *RestController) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Implementation here
	}
}

func (c *RestController) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Implementation here
	}
}

func (c *RestController) Refresh() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Implementation here
	}
}

func (c *RestController) VerifyEmail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Implementation here
	}
}
