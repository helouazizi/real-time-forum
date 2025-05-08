package app

import (
	"database/sql"
	"log"

	"web-forum/internal/database"
	"web-forum/internal/handlers"
	"web-forum/internal/repository"
	"web-forum/internal/services"
	"web-forum/pkg/config"
	"web-forum/pkg/logger"

	_ "github.com/mattn/go-sqlite3"
)

type Application struct {
	DB                 *sql.DB
	Home               *handlers.HomeHandler
	UserHandler        *handlers.UserHandler
	PostHandler        *handlers.PostHandler
	ChatHandler        *handlers.ChatHandler
	ActiveUsersHandler *handlers.ActiveHandler
}

func NewApp(config *config.Configuration) *Application {
	// Connect to SQLite3
	db, err := sql.Open("sqlite3", config.DB_PATH)
	if err != nil {
		logger.LogWithDetails(err)
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Run database migrations to craete tables
	database.Migrate(db)

	// Initialize repositorys
	homeRepo := repository.NewHomeRepository(db)
	userMethods := repository.NewUserRepository(db)
	postMrthods := repository.NewPostRepository(db)
	chatMethods := repository.NewChatRepo(db)
	activeUsers := repository.NewActiveRepo(db)

	// Initialize services
	homeService := services.NewHomeService(homeRepo)
	userService := services.NewUserService(userMethods)
	postServices := services.NewPostService(postMrthods)
	chatServices := services.NewChatService(chatMethods)
	activeServices := services.NewActiveRepo(activeUsers)

	// Initialize handlers
	homeHandler := handlers.NewHomeHandler(homeService)
	userHandler := handlers.NewUserHandler(userService)
	postHandler := handlers.NewPostHandler(postServices)
	// this for chst handler
	hub := handlers.NewHub()
	chatHandler := handlers.NewChatHandler(hub, chatServices)
	activeHandlers := handlers.NewActiveHandler(activeServices)
	go hub.Run()

	return &Application{
		DB:                 db,
		Home:               homeHandler,
		UserHandler:        userHandler,
		PostHandler:        postHandler,
		ChatHandler:        chatHandler,
		ActiveUsersHandler: activeHandlers,
	}
}
