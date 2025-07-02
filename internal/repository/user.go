package repository

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	"web-forum/internal/models"
	"web-forum/internal/utils"
	"web-forum/pkg/logger"

	"github.com/google/uuid"
)

type UserMethods interface {
	CreateUser(user models.User) models.Error
	Login(user models.UserLogin) (models.UserLogin, models.Error)
	Logout(token string) models.Error
	// UpdateUser(user models.User) (models.User, models.Error)
	GetUserInfo(token string) (models.User, models.Error)
	IsUsernameOrEmailTaken(username, email string) models.Error
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user models.User) models.Error {
	user.Email = strings.TrimSpace(user.Email)
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.Nickname = strings.TrimSpace(user.Nickname)
	// lets hash the pass
	hashedPass, err := utils.HashPassWord(user.Password)
	if err != nil {
		logger.LogWithDetails(err)
		return models.Error{
			Message: "Internal server error",
			Code:    http.StatusInternalServerError,
		}
	}

	// check email existnace and nick name
	err1 := r.IsUsernameOrEmailTaken(user.Nickname, user.Email)
	if err1.Code != http.StatusOK {
		logger.LogWithDetails(fmt.Errorf(err1.Message))
		return err1
	}

	// Proceed to insert
	query := `
	INSERT INTO users (
		nickname, age, gender, first_name, last_name, email, password_hash 
	) VALUES (?, ?, ?, ?, ?, ?, ? )
	`
	_, err = r.db.Exec(query,
		template.HTMLEscapeString(user.Nickname),

		user.Age,
		template.HTMLEscapeString(user.Gender),
		template.HTMLEscapeString(user.FirstName),
		template.HTMLEscapeString(user.LastName),
		template.HTMLEscapeString(user.Email),
		hashedPass,
	)
	if err != nil {
		logger.LogWithDetails(err)
		return models.Error{
			Message: "Internal server error",
			Code:    http.StatusInternalServerError,
		}
	}

	// finaly seccefully create a user
	return models.Error{
		Message: "seccefully created your account",
		Code:    http.StatusCreated,
	}
}

func (r *UserRepository) Login(user models.UserLogin) (models.UserLogin, models.Error) {
	query := `SELECT password_hash FROM users WHERE nickname = ?`
	Updatequery := `UPDATE users SET session_token = ?, session_expires_at = ?, is_active = true WHERE nickname = ?`
	isEmail := utils.ValidEmail(user.LoginId)
	if isEmail {
		query = fmt.Sprintf(`SELECT password_hash FROM users WHERE %s = ?`, "email")
		Updatequery = fmt.Sprintf(`UPDATE users SET session_token = ?, session_expires_at = ? , is_active = true WHERE %s = ?`, "email")
	}
	var hash string
	err := r.db.QueryRow(query, user.LoginId).Scan(&hash)
	if err != nil {
		return models.UserLogin{}, models.Error{
			Message: "Bad Request",
			Code:    http.StatusBadRequest,
			UserErrors: models.UserInputErrors{
				HasError: true,
				Nickname: "Invalid nickname or email",
			},
		}
	}

	errCompare := utils.ComparePass([]byte(hash), []byte(user.Password))
	if errCompare != nil {
		logger.LogWithDetails(errCompare)
		return models.UserLogin{}, models.Error{
			Message: "Bad Request",
			Code:    http.StatusBadRequest,
			UserErrors: models.UserInputErrors{
				HasError: true,
				Pass:     "Invalid password",
			},
		}
	}

	// Generate a new token
	newToken := uuid.New().String()

	//  Update the token in database
	_, errUpdate := r.db.Exec(Updatequery, newToken, time.Now().Add(24*time.Hour), user.LoginId) // expires after 24h
	if errUpdate != nil {
		logger.LogWithDetails(errUpdate)
		return models.UserLogin{}, models.Error{
			Message: "Internal Sererver Error",
			Code:    http.StatusInternalServerError,
		}
	}

	user.SessionToken = newToken

	return user, models.Error{
		Message: "Seccefully Loged in",
		Code:    http.StatusOK,
	}
}

func (r *UserRepository) GetUserInfo(token string) (models.User, models.Error) {
	var userInfo models.User
	query := `
		SELECT id, age, gender, first_name, last_name, nickname, email, created_at, updated_at 
		FROM users 
		WHERE session_token = ? AND session_expires_at > CURRENT_TIMESTAMP
	`

	err := r.db.QueryRow(query, token).Scan(
		&userInfo.ID,
		&userInfo.Age,
		&userInfo.Gender,
		&userInfo.FirstName,
		&userInfo.LastName,
		&userInfo.Nickname,
		&userInfo.Email,
		&userInfo.CreatedAt,
		&userInfo.UpdatedAt,
	)
	if err != nil {
		logger.LogWithDetails(err)
		return models.User{}, models.Error{
			Message: "Internal Server Error",
			Code:    http.StatusInternalServerError,
		}
	}

	return userInfo, models.Error{
		Message: "Valid Token",
		Code:    http.StatusOK,
	}
}

func (r *UserRepository) Logout(token string) models.Error {
	query := `UPDATE users SET session_token = NULL, session_expires_at = NULL , is_active = false WHERE session_token = ?`
	result, err := r.db.Exec(query, token)
	if err != nil {
		logger.LogWithDetails(err)
		return models.Error{Message: "Internal server error", Code: http.StatusInternalServerError}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.LogWithDetails(err)
		return models.Error{Message: "Internal server error", Code: http.StatusInternalServerError}
	}

	if rowsAffected == 0 {
		// Token was invalid or already cleared
		return models.Error{Message: "Invalid or expired session", Code: http.StatusUnauthorized}
	}

	return models.Error{Message: "Successfully logged out", Code: http.StatusOK}
}

// func (r *UserRepository) UpdateUser(user models.User) (models.User, models.Error) {

// 	/////////////// have some work here we should as the user for the previouse credential
// 	query := `
// 	UPDATE users
// 	SET username = ?, email = ?, password_hash = ?, bio = ?, updated_at = CURRENT_TIMESTAMP
// 	WHERE id = ?
// 	`
// 	_, err := r.db.Exec(query,
// 		user.Username,
// 		user.Email,
// 		user.PasswordHash,
// 		user.Bio,
// 		user.ID,
// 	)
// 	if err != nil {
// 		logger.LogWithDetails(err)
// 		return models.User{}, models.Error{
// 			Message: "Internal server error",
// 			Code:    http.StatusInternalServerError,
// 		}
// 	}
// 	return user, models.Error{
// 		Message: "seccefully updated information",
// 		Code:    http.StatusOK, // 200
// 	}
// }

func (r *UserRepository) IsUsernameOrEmailTaken(username, email string) models.Error {
	query := `
	SELECT COUNT(*) FROM users
	WHERE nickname = ? OR email = ?
	`
	var count int
	err := r.db.QueryRow(query, username, email).Scan(&count)
	if err != nil {
		logger.LogWithDetails(err)
		return models.Error{
			Message: "Internal server error",
			Code:    http.StatusInternalServerError,
		}
	}

	if count > 0 {
		return models.Error{
			Message:    "Bad Request",
			Code:       http.StatusBadRequest,
			UserErrors: models.UserInputErrors{HasError: true, Email: "Email already taken", Nickname: "Nickname already taken"},
		}
	}

	// No match found
	return models.Error{
		Message: "Username and email are available",
		Code:    http.StatusOK,
	}
}
