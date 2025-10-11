package main

import (
	cmd "Agora/command"
	"Agora/config/database"
	"Agora/controller"
	"Agora/logging"
	"Agora/middleware"
	"Agora/repository"
	"Agora/routes"
	"Agora/service"
	"Agora/websocket"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// ==== Set up logger ====
	logging.SetUpLogger()
	logging.Log.Info("Logger initialized")

	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// DB
	db := database.SetUpPostgreSQLConnection()
	defer database.ClosePostgreSQLConnection(db)

	// Seeder command
	if len(os.Args) > 1 {
		cmd.Command(db)
		return
	}

	// ==== WebSocket Setup ====
	go websocket.WsHub.Run()

	var (
		jwtService = service.NewJWTService()

		userRepo       = repository.NewUserRepository(db)
		userService    = service.NewUserService(userRepo, jwtService)
		userController = controller.NewUserController(userService)

		proposalRepo       = repository.NewProposalRepository(db)
		proposalService    = service.NewProposalService(proposalRepo, userRepo, jwtService)
		proposalController = controller.NewProposalController(proposalService)

		commentRepo       = repository.NewCommentRepository(db)
		commentService    = service.NewCommentService(commentRepo)
		commentController = controller.NewCommentController(commentService)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	routes.PublicRoutes(server, userController)
	routes.AdminRoutes(server, userController, jwtService, proposalController)
	routes.UserRoutes(server, userController, commentController, jwtService)

	server.Static("/assets", "./assets")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	var serve string
	if os.Getenv("APP_ENV") == "localhost" {
		serve = "127.0.0.1:" + port
	} else {
		serve = ":" + port
	}

	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
