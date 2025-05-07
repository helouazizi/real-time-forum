package handlers

import (
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
	if r.Method != http.MethodGet {
		utils.RespondWithJSON(w, http.StatusMethodNotAllowed, models.Error{Message: "Method Not Allowed", Code: http.StatusMethodNotAllowed})
		return
	}

	activeUsers, err1 := h.activeService.GetActiveUsers()
	if err1 != nil {
		logger.LogWithDetails(err1)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, activeUsers)
}
