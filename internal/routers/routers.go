package routers

import (
	"net/http"

	middlewares "web-forum/internal/Middlewares"
	"web-forum/internal/app"
)

func SetupRoutes(h *app.Application) *http.ServeMux {
	mux := http.NewServeMux()

	// Serve frontend index
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "front-end/index.html")
	})

	// Serve static files (JS/CSS)
	mux.Handle("/front-end/static/", http.StripPrefix("/front-end/static/", http.FileServer(http.Dir("front-end/static"))))

	// --- User routes ---
	mux.HandleFunc("/api/v1/users/register", h.UserHandler.CreateUser)
	mux.HandleFunc("/api/v1/users/login", h.UserHandler.Login)
	mux.HandleFunc("/api/v1/users/logout", h.UserHandler.Logout)
	mux.HandleFunc("/api/v1/users/info", h.UserHandler.GetUserInfo)

	// --- Post routes (with auth middleware) ---
	mux.HandleFunc("/api/v1/posts", h.PostHandler.FetchPosts)
	mux.Handle("/api/v1/posts/create", middlewares.AuthMiddleware(http.HandlerFunc(h.PostHandler.CreatePost), h.DB))
	mux.Handle("/api/v1/posts/react", middlewares.AuthMiddleware(http.HandlerFunc(h.PostHandler.ReactToPost), h.DB))
	mux.Handle("/api/v1/posts/addComment", middlewares.AuthMiddleware(http.HandlerFunc(h.PostHandler.CommentPost), h.DB))
	mux.Handle("/api/v1/posts/fetchComments", middlewares.AuthMiddleware(http.HandlerFunc(h.PostHandler.FetchComments), h.DB))
	mux.Handle("/api/v1/posts/filter", middlewares.AuthMiddleware(http.HandlerFunc(h.PostHandler.FilterPosts), h.DB))

	// --- Chat routes ---
	mux.Handle("/api/v1/chat", middlewares.AuthMiddleware(http.HandlerFunc(h.ChatHandler.HandleChat), h.DB))
	mux.Handle("/api/v1/active", middlewares.AuthMiddleware(http.HandlerFunc(h.ActiveUsersHandler.GetActiveUsers), h.DB))

	return mux
}
