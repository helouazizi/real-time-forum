package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"web-forum/internal/models"
	"web-forum/internal/services"
	"web-forum/internal/utils"
	"web-forum/pkg/logger"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}


func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithJSON(w, http.StatusMethodNotAllowed, models.Error{Message: "Method Not Allowed", Code: http.StatusMethodNotAllowed})
		return
	}
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.LogWithDetails(err)
		utils.RespondWithJSON(w, http.StatusBadRequest, models.Error{Message: "Bad Request", Code: http.StatusBadRequest})
		return
	}

	// lets validate the user inputs
	err := utils.ValidateUserInputs(user)
	if err.Code == http.StatusBadRequest {
		logger.LogWithDetails(fmt.Errorf("invalid user credentials"))
		utils.RespondWithJSON(w, err.Code, err)
		return
	}
	err1 := h.userService.CreateUser(user)
	if err1.Code != http.StatusCreated {
		logger.LogWithDetails(fmt.Errorf(err.Message))
		utils.RespondWithJSON(w, err.Code, err1)
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, models.SuccesMessage{Message: "Seccefully created your account"})
}


func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithJSON(w, http.StatusMethodNotAllowed, models.Error{Message: "Methos Not Allowed", Code: http.StatusMethodNotAllowed})
		return
	}
	var user models.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.LogWithDetails(err)
		utils.RespondWithJSON(w, http.StatusBadRequest, models.Error{Message: "Bad Request", Code: http.StatusBadRequest})
		return
	}

	LogedUser, err := h.userService.Login(user)
	if err.Code != http.StatusOK {
		logger.LogWithDetails(fmt.Errorf(err.Message))
		utils.RespondWithJSON(w, err.Code, err)
		return
	}

	// After successful login, set the cookie
	cookie := &http.Cookie{Name: "Token", Value: LogedUser.SessionToken, MaxAge: 3600, HttpOnly: true, SameSite: http.SameSiteStrictMode, Path: "/", Secure: false}
	http.SetCookie(w, cookie)

	utils.RespondWithJSON(w, http.StatusOK, models.SuccesMessage{Message: "Seccefuly loged in"})
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithJSON(w, http.StatusMethodNotAllowed, models.Error{Message: "Methos Not Allowed", Code: http.StatusMethodNotAllowed})
		return
	}

	token, err := utils.GetToken(r, "Token")
	if err.Code != http.StatusOK {
		logger.LogWithDetails(fmt.Errorf(err.Message))
		utils.RespondWithJSON(w, err.Code, err)
		return
	}

	err1 := h.userService.Logout(token)
	if err1.Code != http.StatusOK {
		logger.LogWithDetails(fmt.Errorf(err1.Message))
		utils.RespondWithJSON(w, err1.Code, err)
		return
	}

	// After successful login, set the cookie
	cookiee := &http.Cookie{Name: "Token", Value: "", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteStrictMode, Path: "/", Secure: false}
	http.SetCookie(w, cookiee)

	utils.RespondWithJSON(w, http.StatusOK, models.SuccesMessage{Message: "Seccefuly loged out"})
}

func (h *UserHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithJSON(w, http.StatusMethodNotAllowed, models.Error{Message: "Methos Not Allowed", Code: http.StatusMethodNotAllowed})
		return
	}
	token, err := utils.GetToken(r, "Token")
	if err.Code != http.StatusOK {
		logger.LogWithDetails(fmt.Errorf(err.Message))
		utils.RespondWithJSON(w, err.Code, err)
		return
	}

	userInfo, err1 := h.userService.GetUserInfo(token)
	if err1.Code != http.StatusOK {
		logger.LogWithDetails(fmt.Errorf(err1.Message))
		utils.RespondWithJSON(w, err1.Code,err1)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, userInfo)
}







// func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		utils.RespondWithError(w, models.Error{Message: "Methos Not Allowed", Code: http.StatusMethodNotAllowed})
// 		return
// 	}
// 	var user models.User
// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		logger.LogWithDetails(err)
// 		utils.RespondWithError(w, models.Error{Message: "Bad Request", Code: http.StatusBadRequest})
// 		return
// 	}
// 	updatedUser, err := h.userService.UpdateUser(user)
// 	if err.Code != http.StatusOK {
// 		logger.LogWithDetails(fmt.Errorf(err.Message))
// 		utils.RespondWithError(w, err)
// 		return
// 	}
// 	utils.RespondWithJSON(w, http.StatusCreated, updatedUser)
// }