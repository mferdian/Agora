package routes

import (
	"Agora/constants"
	"Agora/controller"
	"Agora/middleware"
	"Agora/service"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(r *gin.Engine, userController controller.IUserController,
	jwtService service.InterfaceJWTService) {
	admin := r.Group("/api/admin")
	admin.Use(middleware.Authentication(jwtService))
	admin.Use(middleware.AuthorizeRole(constants.ENUM_ROLE_ADMIN))

	// User management
	admin.POST("/create-user", userController.CreateUser)
	admin.GET("/get-all-users", userController.GetAllUser)
}
