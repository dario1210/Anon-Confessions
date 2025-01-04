package app

import (
	"anon-confessions/cmd/internal/comments"
	"anon-confessions/cmd/internal/config"
	"anon-confessions/cmd/internal/db"
	"anon-confessions/cmd/internal/posts"
	"anon-confessions/cmd/internal/user"
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

	// Repositories
	userRepo := user.NewSQLiteUserRepository(dbConn)
	postsRepo := posts.NewSQLitePostsRepository(dbConn)
	commentsRepo := comments.NewSQLiteCommentsRepository(dbConn)

	// Services
	userService := user.NewUserService(userRepo)
	postsService := posts.NewPostsService(postsRepo)
	commentsService := comments.NewCommentsService(commentsRepo)

	// Handlers
	userHandler := user.NewUserHandler(userService)
	postsHandler := posts.NewPostsHandler(postsService)
	commentsHandler := comments.NewCommentsHandler(commentsService, postsService)

	handlers := &HandlerContainer{
		UserHandler:     userHandler,
		PostsHandler:    postsHandler,
		CommentsHandler: commentsHandler,
	}
	router := setupRouter(handlers, dbConn)

	app := &App{
		Config: *cfg,
		DB:     dbConn,
		Router: router,
		// Logger : logger,
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

func setupRouter(h *HandlerContainer, db *gorm.DB) *gin.Engine {
	router := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")
	{
		posts.RegisterPostRoutes(api, h.PostsHandler, db)
		user.RegisterUsersRoutes(api, h.UserHandler, db)
		comments.RegisterCommentsRoutes(api, h.CommentsHandler, db)
	}

	return router
}
