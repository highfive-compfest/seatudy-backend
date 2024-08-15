package wallet

import (
	"github.com/gin-gonic/gin"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/middleware"
	"github.com/highfive-compfest/seatudy-backend/internal/response"
	"net/http"
)

type RestController struct {
	uc     *UseCase
	midtUc IMidtransUseCase
}

func NewRestController(engine *gin.Engine, uc *UseCase, midtUc IMidtransUseCase) {
	controller := &RestController{uc: uc, midtUc: midtUc}

	walletGroup := engine.Group("/v1/wallets")
	{
		walletGroup.POST("/top-up",
			middleware.Authenticate(),
			middleware.RequireEmailVerified(),
			middleware.RequireRole("student"),
			controller.TopUp(),
		)
		walletGroup.POST("/verify-payment/midtrans",
			controller.VerifyPayment(),
		)
		walletGroup.GET("/balance",
			middleware.Authenticate(),
			controller.GetBalance(),
		)
		walletGroup.GET("/midtrans-transactions",
			middleware.Authenticate(),
			controller.GetMidtransTransactions(),
		)
	}
}

func (c *RestController) TopUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req TopUpRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			err2 := apierror.ErrValidation.Build()
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		res, err := c.uc.TopUp(ctx, &req)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "TOP_UP_REQUEST_SUCCESS", res).Send(ctx)
	}
}

func (c *RestController) VerifyPayment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var notificationPayload map[string]any
		if err := ctx.ShouldBind(&notificationPayload); err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		err := c.midtUc.VerifyPayment(notificationPayload)
		if err != nil {
			ctx.Status(apierror.GetHttpStatus(err))
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func (c *RestController) GetBalance() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := c.uc.GetBalance(ctx)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "GET_BALANCE_SUCCESS", res).Send(ctx)
	}
}

func (c *RestController) GetMidtransTransactions() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req GetMidtransTransactionsRequest
		if err := ctx.ShouldBindQuery(&req); err != nil {
			err2 := apierror.ErrValidation.Build()
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), err.Error()).Send(ctx)
			return
		}

		res, err := c.uc.GetMidtransTransactionsByUser(ctx, &req)
		if err != nil {
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), apierror.GetPayload(err)).Send(ctx)
			return
		}

		response.NewRestResponse(http.StatusOK, "GET_MIDTRANS_TRANSACTIONS_SUCCESS", res).Send(ctx)
	}
}
