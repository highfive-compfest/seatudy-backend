package user

import "github.com/gin-gonic/gin"

type RestController struct {
	uc *UseCase
}

func NewRestController(engine *gin.Engine, uc *UseCase) *RestController {
	return &RestController{uc: uc}
}
