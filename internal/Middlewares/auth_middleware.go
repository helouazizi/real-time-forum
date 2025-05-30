package middlewares

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
	"web-forum/internal/auth"
	"web-forum/internal/models"
	"web-forum/internal/utils"
	"web-forum/pkg/logger"
)

var auth_rateLimiter = auth.NewRateLimiter(30, time.Second) // 30 requests per second limit
func AuthMiddleware(next http.Handler, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Check if the request contains a token
		ip := r.RemoteAddr

		// Check rate limit (applies to both login and registration)
		if auth_rateLimiter.CheckRateLimit(ip) {

			logger.LogWithDetails(fmt.Errorf("too many requeste"))
			utils.RespondWithJSON(w, http.StatusTooManyRequests, models.Error{ Code : 429, Message: "too many request please try later"})
			return
		}

		token, err := utils.GetToken(r, "Token")
		if err.Code != http.StatusOK {
			logger.LogWithDetails(fmt.Errorf(err.Message))
			utils.RespondWithJSON(w, err.Code, err)
			return
		}
		const query = `
			SELECT id, session_expires_at 
			FROM users 
			WHERE session_token = ?
		`

		var id int
		var exipration time.Time
		err1 := db.QueryRow(query, token).Scan(&id, &exipration)
		if err1 == sql.ErrNoRows {
			logger.LogWithDetails(fmt.Errorf("invalid token"))
			utils.RespondWithJSON(w, http.StatusUnauthorized, models.Error{
				Message: "Invalid or expired token",
				Code:    http.StatusUnauthorized,
			})
			return
		}

		if err1 != nil {
			logger.LogWithDetails(err1)
			utils.RespondWithJSON(w, http.StatusInternalServerError, models.Error{
				Message: "Internal Server Error",
				Code:    http.StatusInternalServerError,
			})
			return
		}

		if time.Now().After(exipration) {
			logger.LogWithDetails(fmt.Errorf("token expired at %v", exipration))
			utils.RespondWithJSON(w, http.StatusUnauthorized, models.Error{
				Message: "Session has expired. Please log in again.",
				Code:    http.StatusUnauthorized,
			})
			return
		}

		// Proceed to next handler
		next.ServeHTTP(w, r)
	})
}
