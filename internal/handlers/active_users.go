package handlers

import (
	"encoding/json"
	"net/http"

	"web-forum/internal/models"
	"web-forum/internal/services"
	"web-forum/internal/utils"
	"web-forum/pkg/logger"
)

type ActiveHandler struct {
	activeService *services.ActiveService
}

func NewActiveHandler(activeService *services.ActiveService) *ActiveHandler {
	return &ActiveHandler{activeService: activeService}
}

func (h *ActiveHandler) GetActiveUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithJSON(w, http.StatusMethodNotAllowed, models.Error{Message: "Method Not Allowed", Code: http.StatusMethodNotAllowed})
		return
	}
	var user models.ClientRegistration
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.LogWithDetails(err)
		utils.RespondWithJSON(w, http.StatusBadRequest, models.Error{Message: "Bad Request", Code: http.StatusBadRequest})
		return
	}
	// fmt.Println(user.SenderId,"user id")

	activeUsers, err1 := h.activeService.GetActiveUsers(user.SenderId)
	if err1 != nil {
		logger.LogWithDetails(err1)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, activeUsers)
}
