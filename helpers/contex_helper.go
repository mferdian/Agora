package helpers

import (
	"Agora/constants"
	"context"
)

func GetUserID(ctx context.Context) string {
	if val := ctx.Value(constants.ContextUserIDKey); val != nil {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return ""
}

func GetUserRole(ctx context.Context) string {
	if val := ctx.Value(constants.ContextRoleKey); val != nil {
		if role, ok := val.(string); ok {
			return role
		}
	}
	return ""
}
