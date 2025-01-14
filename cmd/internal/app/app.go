package app

import (
	"anon-confessions/cmd/internal/config"
	"anon-confessions/cmd/internal/db"
	"anon-confessions/cmd/internal/middleware"
	"anon-confessions/cmd/internal/modules/comments"
	"anon-confessions/cmd/internal/modules/posts"
	"anon-confessions/cmd/internal/modules/user"
	"anon-confessions/cmd/internal/websocket"
	"anon-confessions/docs"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type App struct {
	Config config.Config
	DB     *gorm.DB
	Router *gin.Engine
}

type HandlerContainer struct {
	UserHandler     *user.UserHandler
	PostsHandler    *posts.PostsHandler
	CommentsHandler *comments.CommentsHandler
}

// @title           Anonymous Confessions API
// @version         1.0
// @description     A privacy-focused backend service that allows users to:
// @description     • Post and manage anonymous confessions.
// @description     • React to posts with likes and comments.
// @description     • Leave comments on confessions.
// @description     • Receive real-time updates through WebSocket.
// @description
// @description     The API is designed with RESTful principles, uses SQLite for data storage, and ensures anonymity without storing personal information.

// @host            localhost: cfg.Port
// @BasePath        /api/v1
//
// @securityDefinitions.apikey AccountNumberAuth
// @in               header
// @name             X-Account-Number
// @description      A unique account number for user authentication.
//
// @security         AccountNumberAuth
func swaggerInfo() {}

// NewApp initializes the application by loading configuration, setting up database connection, running migrations, and initializing services and handlers.
func NewApp(cfg *config.Config) (*App, error) {
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.Port)

	slog.Info("Initializing database connection and running migrations...")
	dbConn, err := db.DbConnection(cfg.DB.File)
	if err != nil {
		slog.Error("Failed to connect to database", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	slog.Info("Starting WebSocket hub...")
	hub := websocket.NewHub()
	go hub.Run()

	// MIDDLEWARE
	slog.Info("Setting up middleware...")
	authMiddleware := middleware.Authentication(dbConn)

	// Repositories
	slog.Info("Initializing repositories...")
	userRepo := user.NewSQLiteUserRepository(dbConn)
	postsRepo := posts.NewSQLitePostsRepository(dbConn)
	commentsRepo := comments.NewSQLiteCommentsRepository(dbConn)

	// Services
	slog.Info("Initializing services...")
	userService := user.NewUserService(userRepo)
	postsService := posts.NewPostsService(postsRepo, hub)
	commentsService := comments.NewCommentsService(commentsRepo, hub)

	// Handlers
	slog.Info("Initializing handlers...")
	userHandler := user.NewUserHandler(userService)
	postsHandler := posts.NewPostsHandler(postsService)
	commentsHandler := comments.NewCommentsHandler(commentsService, postsService)

	handlers := &HandlerContainer{
		UserHandler:     userHandler,
		PostsHandler:    postsHandler,
		CommentsHandler: commentsHandler,
	}

	slog.Info("Setting up router...")
	router := setupRouter(handlers, authMiddleware, hub)

	slog.Info("Application initialized successfully")
	app := &App{
		Config: *cfg,
		DB:     dbConn,
		Router: router,
	}

	return app, nil
}

func (a *App) Run() error {
	slog.Info("Starting HTTP server", slog.String("port", a.Config.Port))
	if err := a.Router.Run(fmt.Sprintf(":%v", a.Config.Port)); err != nil {
		slog.Error("Server failed to start", slog.String("error", err.Error()))
		return fmt.Errorf("server failed to start: %w", err)
	}
	slog.Info("HTTP server started successfully")
	return nil
}

func setupRouter(h *HandlerContainer, authMiddleware gin.HandlerFunc, hub *websocket.Hub) *gin.Engine {
	router := gin.Default()

	// Swagger documentation route
	slog.Info("Setting up Swagger documentation routes...")
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Base API group
	api := router.Group("/api/v1")

	// Routes that require authentication
	authenticated := api.Group("/")
	authenticated.Use(authMiddleware)
	{
		posts.RegisterPostRoutes(authenticated, h.PostsHandler)
		comments.RegisterCommentsRoutes(authenticated, h.CommentsHandler)
	}

	// Routes that do not require authentication
	user.RegisterUsersRoutes(api, h.UserHandler)

	// WebSocket routes
	websocket.RegisterWebSocketRoutes(api, hub)

	slog.Info("Router setup complete")
	return router
}
