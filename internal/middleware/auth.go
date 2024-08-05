package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/highfive-compfest/seatudy-backend/internal/apierror"
	"github.com/highfive-compfest/seatudy-backend/internal/jwtoken"
	"github.com/highfive-compfest/seatudy-backend/internal/response"
	"strings"
	"time"
)

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearer := ctx.GetHeader("Authorization")
		if bearer == "" {
			err := apierror.ErrTokenEmpty
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), nil).Send(ctx)
			ctx.Abort()
			return
		}

		tokenSlice := strings.Split(bearer, " ")
		if len(tokenSlice) != 2 {
			err := apierror.ErrTokenInvalid
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), nil).Send(ctx)
			ctx.Abort()
			return
		}

		token := tokenSlice[1]

		claims, err := jwtoken.DecodeAccessJWT(token)
		if err != nil {
			err2 := apierror.ErrTokenInvalid
			response.NewRestResponse(apierror.GetHttpStatus(err2), err2.Error(), nil).Send(ctx)
			ctx.Abort()
			return
		}

		if claims.Issuer != "seatudy-backend-accesstoken" {
			err := apierror.ErrTokenInvalid
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), nil).Send(ctx)
			ctx.Abort()
			return
		}

		if claims.ExpiresAt.Time.Before(time.Now()) {
			err := apierror.ErrTokenExpired
			response.NewRestResponse(apierror.GetHttpStatus(err), err.Error(), nil).Send(ctx)
			ctx.Abort()
			return
		}

		ctx.Set("user.id", claims.Subject)
		ctx.Set("user.email", claims.Email)
		ctx.Set("user.is_email_verified", claims.IsEmailVerified)
		ctx.Set("user.name", claims.Name)
		ctx.Set("user.role", claims.Role)
		ctx.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Implementation here
	}
}
