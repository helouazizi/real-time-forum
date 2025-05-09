package handlers

// import (
// 	"net/http"

// 	"web-forum/internal/models"
// 	"web-forum/internal/services"
// 	"web-forum/internal/utils"
// )

// type HomeHandler struct {
// 	HomeService *services.HomeService
// }

// func NewHomeHandler(HomeService *services.HomeService) *HomeHandler {
// 	return &HomeHandler{HomeService: HomeService}
// }

// func (h *HomeHandler) FetchPosts(w http.ResponseWriter, r *http.Request) {
// 	// if r.URL.Path != "/" {
// 	// 	utils.RespondWithJSON(w, http.StatusNotFound, models.Error{Message: "Page Not Found", Code: http.StatusNotFound})
// 	// 	return
// 	// }
// 	if r.Method != http.MethodGet {
// 		utils.RespondWithJSON(w, http.StatusMethodNotAllowed, models.Error{Message: "Method Not Allowed", Code: http.StatusMethodNotAllowed})
// 		return
// 	}
// 	Posts, err := h.HomeService.FetchPosts()
// 	if err.Code != http.StatusOK {
// 		utils.RespondWithJSON(w, http.StatusInternalServerError, models.Error{
// 			Message: "Internal server error",
// 			Code:    http.StatusInternalServerError,
// 		})
// 		return
// 	}
// 	utils.RespondWithJSON(w, http.StatusOK, Posts)
// }
