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
	jwtService service.InterfaceJWTService,
) {
	user := r.Group("/api/users")
	user.Use(middleware.Authentication(jwtService))

	// --- User Routes ---
	user.PATCH("/update-profile/:id", userController.UpdateUser)
	user.GET("/get-detail-user/:id", userController.GetUserByID)
	user.DELETE("/delete-profile/:id", userController.DeleteUser)
	
}
