package middleware

import (
	"Agora/constants"
	"Agora/logging"
	"Agora/service"
	"Agora/utils"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authentication(jwtService service.IJWTService) gin.HandlerFunc {
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

		// Simpan ke context bawaan request (untuk service)
		newCtx := context.WithValue(ctx.Request.Context(), constants.ContextUserIDKey, claims.UserID)
		newCtx = context.WithValue(newCtx, constants.ContextRoleKey, claims.Role)
		ctx.Request = ctx.Request.WithContext(newCtx)

		// Simpan ke gin.Context (untuk middleware lain seperti AuthorizeRole)
		ctx.Set("id", claims.UserID)
		ctx.Set("role", claims.Role)

		logging.Log.Infof("Authenticated request - UserID: %s, Role: %s", claims.UserID, claims.Role)

		ctx.Next()
	}
}
