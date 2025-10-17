package middleware

import (
	"errors"
	"fmt"
	"github.com/WhiteMaks/go_academy_portal_service/internal/service"
	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	AuthorizationPayload = "authorization_payload"
)

func AuthMiddleware(tokenMaker util.TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader("authorization")
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is required")

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, service.PrepareErrorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 {
			err := errors.New("invalid authorization header format")

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, service.PrepareErrorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != "bearer" {
			err := fmt.Errorf("invalid authorization header type %s", authorizationType)

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, service.PrepareErrorResponse(err))
			return
		}

		token := fields[1]
		payload, err := tokenMaker.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, service.PrepareErrorResponse(err))
			return
		}

		ctx.Set(AuthorizationPayload, payload)
		ctx.Next()
	}
}
