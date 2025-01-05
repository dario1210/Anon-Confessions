package app

import (
	"anon-confessions/cmd/internal/comments"
	"anon-confessions/cmd/internal/config"
	"anon-confessions/cmd/internal/db"
	"anon-confessions/cmd/internal/middleware"
	"anon-confessions/cmd/internal/posts"
	"anon-confessions/cmd/internal/user"
	"anon-confessions/cmd/internal/websocket"
	"anon-confessions/docs"
	"fmt"
	"log"

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

func NewApp() (*App, error) {
	log.Println("Loading configuration...")
	cfg := config.LoadConfig()

	log.Println("Initializing database connection & Migrations ...")
	dbConn, err := db.DbConnection(cfg.DB.File)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.RunMigrations(&cfg.Migrations); err != nil {
		return nil, fmt.Errorf("database migrations failed: %w", err)
	}

	// Websocket
	hub := websocket.NewHub()
	go hub.Run()

	// MIDDLEWARE
	authMiddleware := middleware.Authentication(dbConn)

	// Repositories
	userRepo := user.NewSQLiteUserRepository(dbConn)
	postsRepo := posts.NewSQLitePostsRepository(dbConn)
	commentsRepo := comments.NewSQLiteCommentsRepository(dbConn)

	// Services
	userService := user.NewUserService(userRepo)
	postsService := posts.NewPostsService(postsRepo, hub)
	commentsService := comments.NewCommentsService(commentsRepo, hub)

	// Handlers
	userHandler := user.NewUserHandler(userService)
	postsHandler := posts.NewPostsHandler(postsService)
	commentsHandler := comments.NewCommentsHandler(commentsService, postsService)

	handlers := &HandlerContainer{
		UserHandler:     userHandler,
		PostsHandler:    postsHandler,
		CommentsHandler: commentsHandler,
	}

	router := setupRouter(handlers, authMiddleware, hub)

	app := &App{
		Config: *cfg,
		DB:     dbConn,
		Router: router,
	}

	return app, nil

}

func (a *App) Run() error {
	log.Printf("Starting the HTTP server on %v...", a.Config.Port)
	if err := a.Router.Run(fmt.Sprintf(":%v", a.Config.Port)); err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}
	return nil
}

func setupRouter(h *HandlerContainer, authMiddleware gin.HandlerFunc, hub *websocket.Hub) *gin.Engine {
	router := gin.Default()

	// Swagger documentation route
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

	return router
}
