package middleware

import (
	"Agora/constants"
	"Agora/logging"
	"Agora/service"
	"Agora/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authentication(jwtService service.InterfaceJWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			logging.Log.Warn("Authorization header not found")
			res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_PROSES_REQUEST, constants.MESSAGE_FAILED_TOKEN_NOT_FOUND, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			logging.Log.Warn("Authorization header format invalid")
			res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_PROSES_REQUEST, constants.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, claims, err := jwtService.ValidateToken(tokenStr)
		if err != nil || !token.Valid {
			logging.Log.Warnf("Invalid token: %v", err)
			res := utils.BuildResponseFailed(constants.MESSAGE_FAILED_PROSES_REQUEST, constants.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		logging.Log.Infof("Authenticated request - UserID: %s, Role: %s", claims.UserID, claims.Role)

		ctx.Set("Authorization", tokenStr)
		ctx.Set("id", claims.ID)
		ctx.Set("role", claims.Role)

		ctx.Next()
	}
}
