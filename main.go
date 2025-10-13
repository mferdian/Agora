package main

import (
	"Agora/command"
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
	// ==== Setup Logger ====
	logging.SetUpLogger()
	logging.Log.Info("Logger initialized")

	// ==== Load .env ====
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	// ==== Database Connection ====
	db := database.SetUpPostgreSQLConnection()
	defer database.ClosePostgreSQLConnection(db)

	// ==== Seeder Mode ====
	if len(os.Args) > 1 {
		cmd.Command(db)
		return
	}

	// ==== Repository Layer ====
	userRepo := repository.NewUserRepository(db)
	proposalRepo := repository.NewProposalRepository(db)
	commentRepo := repository.NewCommentRepository(db)

	// ==== Service Layer ====
	jwtService := service.NewJWTService()
	userService := service.NewUserService(userRepo, jwtService)
	proposalService := service.NewProposalService(proposalRepo, userRepo, jwtService)
	commentService := service.NewCommentService(commentRepo, userRepo, proposalRepo, jwtService)

	// ==== WebSocket Hub ====
	hub := websocket.NewHub(commentService)
	go hub.Run()

	// ==== Controller Layer ====
	userController := controller.NewUserController(userService)
	proposalController := controller.NewProposalController(proposalService)

	// ==== WebSocket Handler ====
	wsHandler := &websocket.WsHandler{
		Hub: hub,
	}

	// ==== Gin Server Setup ====
	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	// ==== Routes ====
	// ==== Routes ====
	routes.PublicRoutes(server, userController)
	routes.AdminRoutes(server, userController, jwtService, proposalController)
	routes.UserRoutes(server, userController, jwtService)
	routes.WebSocketRoutes(server, wsHandler)

	// ==== Static Files ====
	server.Static("/assets", "./assets")

	// ==== Port Config ====
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	addr := ":" + port
	if os.Getenv("APP_ENV") == "localhost" {
		addr = "127.0.0.1:" + port
	}

	log.Printf("🚀 Server running at http://%s", addr)
	if err := server.Run(addr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
