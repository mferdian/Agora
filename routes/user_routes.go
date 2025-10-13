package routes

import (
	"Agora/controller"
	"Agora/middleware"
	"Agora/service"

	"github.com/gin-gonic/gin"
)

func UserRoutes(
	r *gin.Engine,
	userController controller.IUserController,
	jwtService service.IJWTService,
) {
	api := r.Group("/api")
	api.Use(middleware.Authentication(jwtService))

	// === USER ROUTES ===
	users := api.Group("/users")
	{
		users.GET("", userController.GetAllUser)
		users.GET("/:id", userController.GetUserByID)
		users.POST("", userController.CreateUser)
		users.PATCH("/:id", userController.UpdateUser)
		users.DELETE("/:id", userController.DeleteUser)
	}
}
