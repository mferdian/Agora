package routes

import (
	"Agora/controller"

	"github.com/gin-gonic/gin"
)

func PublicRoutes(r *gin.Engine, userController controller.IUserController) {
	public := r.Group("/api/users")
	public.POST("/register", userController.Register)
	public.POST("/login", userController.Login)
}
