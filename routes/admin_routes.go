package routes

import (
	"Agora/constants"
	"Agora/controller"
	"Agora/middleware"
	"Agora/service"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(
	r *gin.Engine,
	userController controller.IUserController,
	jwtService service.IJWTService,
	proposalController controller.IProposalController,
) {
	admin := r.Group("/api/admins")
	admin.Use(middleware.Authentication(jwtService))
	admin.Use(middleware.AuthorizeRole(constants.ENUM_ROLE_ADMIN))

	users := admin.Group("/users")
	{
		users.POST("", userController.CreateUser)       
		users.GET("", userController.GetAllUser)       
		users.GET("/:id", userController.GetUserByID)   
		users.PATCH("/:id", userController.UpdateUser)  
		users.DELETE("/:id", userController.DeleteUser) 
	}

	proposals := admin.Group("/proposals")
	{
		proposals.POST("", proposalController.CreateProposal)       
		proposals.GET("", proposalController.GetAllProposal)        
		proposals.GET("/:id", proposalController.GetProposalById)   
		proposals.PATCH("/:id", proposalController.UpdateProposal)  
		proposals.DELETE("/:id", proposalController.DeleteProposal) 
	}
}
