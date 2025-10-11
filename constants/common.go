package constants


type ContextKey string

const (
	ENUM_ROLE_ADMIN = "admin"
	ENUM_ROLE_USER  = "user"

	ENUM_RUN_PRODUCTION = "production"
	ENUM_RUN_TESTING    = "testing"

	ENUM_PAGINATION_LIMIT = 10
	ENUM_PAGINATION_PAGE  = 1

	ContextUserIDKey ContextKey = "user_id"
	ContextRoleKey   ContextKey = "role"
)