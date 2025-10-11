package routes

import (
	"Agora/controller"
	"Agora/middleware"
	"Agora/service"
	"Agora/websocket"

	"github.com/gin-gonic/gin"
)

func UserRoutes(
	r *gin.Engine,
	userController controller.IUserController,
	commentController controller.ICommentController,
	jwtService service.IJWTService,
) {
	api := r.Group("/api")
	api.Use(middleware.Authentication(jwtService))

	// --- USERS ---
	users := api.Group("/users")
	{
		users.GET("", userController.GetAllUser)       // GET /api/users
		users.GET("/:id", userController.GetUserByID)    // GET /api/users/:id
		users.POST("", userController.CreateUser)       // POST /api/users
		users.PATCH("/:id", userController.UpdateUser)     // PUT /api/users/:id
		users.DELETE("/:id", userController.DeleteUser)  // DELETE /api/users/:id
	}

	// --- COMMENTS ---
	comments := api.Group("/comments")
	{
		comments.GET("", websocket.ServeWs)
		comments.POST("", commentController.CreateComment) 
	}
}
