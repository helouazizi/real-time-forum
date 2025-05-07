// internal/utils/sweet.go
package utils

import (
	"net/http"
	"regexp"
	"slices"
	"strings"
	"web-forum/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func ValidUsername(username string) bool {
	// Example: Allow only alphanumeric characters and underscores
	match, _ := regexp.MatchString("^[a-zA-Z0-9_]{3,50}$", username) // we can add length like {3,15} and remove the +
	return match && !isReservedUsername(username)
}

// IsValidEmail checks if an email is syntactically valid using regex.
func ValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// Function to check if a username is reserved
func isReservedUsername(username string) bool {
	reservedWords := []string{"admin", "root", "system", "superuser"}
	return slices.Contains(reservedWords, strings.ToLower(username))
}

func StrongPassword(password string) bool {

	hasLower := false
	hasUpper := false
	hasDigit := false

	// Loop through the password to check for lowercase letters and digits
	for _, char := range password {
		if char >= 'a' && char <= 'z' {
			hasLower = true
		}
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		}
		if char >= '0' && char <= '9' {
			hasDigit = true
		}
	}

	// Password is strong if it contains at least one lowercase letter and one digit
	return hasLower && hasDigit && hasUpper
}

func HashPassWord(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(bytes), err
}

func ComparePass(hashed, pass []byte) error {
	return bcrypt.CompareHashAndPassword(hashed, pass)
}

// this function validate the user inputs

func ValidateUserInputs(user models.User) models.Error {
	var errors models.Error
	// lets check the nickname first
	if !ValidUsername(user.Nickname) {
		errors.Code = http.StatusBadRequest
		errors.UserErrors.HasError = true
		errors.UserErrors.Nickname = "Invalid Nickname"
	}
	// lets check the email
	if !ValidEmail(user.Email) {
		errors.Code = http.StatusBadRequest
		errors.UserErrors.HasError = true
		errors.UserErrors.Email = "Invalid Email"
	}

	// lets check the pass
	if !StrongPassword(user.Password) {
		errors.Code = http.StatusBadRequest
		errors.UserErrors.HasError = true
		errors.UserErrors.Pass = "Password is too weak"
	}

	// check the age
	if user.Age < 1 || user.Age > 100 {
		errors.Code = http.StatusBadRequest
		errors.UserErrors.HasError = true
		errors.UserErrors.Age = "Invalid Age"
	}

	// check the gender
	if user.Gender != "male" && user.Gender != "female" {
		errors.Code = http.StatusBadRequest
		errors.UserErrors.HasError = true
		errors.UserErrors.Gender = "Invalid Gender"
	}

	// check the last and first names
	if user.LastName == "" || len(user.LastName) < 3 || len(user.LastName) > 50 {
		errors.Code = http.StatusBadRequest
		errors.UserErrors.HasError = true
		errors.UserErrors.LastName = "Invalid Lastname"
	}
	if user.FirstName == "" || len(user.FirstName) < 3 || len(user.FirstName) > 50 {
		errors.Code = http.StatusBadRequest
		errors.UserErrors.HasError = true
		errors.UserErrors.FirstName = "Invalid Firstname"
	}

	return errors
}

func GetToken(r *http.Request, name string) (string, models.Error) {
	cookie, err := r.Cookie(name)
	if err != nil || cookie.Value == "" {
		// Return unauthorized if there's no cookie or cookie is empty
		return "", models.Error{
			Message: "Token not found or invalid",
			Code:    http.StatusUnauthorized,
		}
	}
	token := cookie.Value

	return token, models.Error{Message: "User Has Token", Code: http.StatusOK}
}
