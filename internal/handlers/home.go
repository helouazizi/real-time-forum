package handlers

import (
	"log"
	"net/http"
	"os"
	"text/template"
	"web-forum/pkg/logger"
)

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

//	func (h *HomeHandler) FetchPosts(w http.ResponseWriter, r *http.Request) {
//		// if r.URL.Path != "/" {
//		// 	utils.RespondWithJSON(w, http.StatusNotFound, models.Error{Message: "Page Not Found", Code: http.StatusNotFound})
//		// 	return
//		// }
//		if r.Method != http.MethodGet {
//			utils.RespondWithJSON(w, http.StatusMethodNotAllowed, models.Error{Message: "Method Not Allowed", Code: http.StatusMethodNotAllowed})
//			return
//		}
//		Posts, err := h.HomeService.FetchPosts()
//		if err.Code != http.StatusOK {
//			utils.RespondWithJSON(w, http.StatusInternalServerError, models.Error{
//				Message: "Internal server error",
//				Code:    http.StatusInternalServerError,
//			})
//			return
//		}
//		utils.RespondWithJSON(w, http.StatusOK, Posts)
//	}
func parsefile(url string) (*template.Template, error) {
	templ, err := template.ParseFiles(url)

	return templ, err

}
func Servstatique(w http.ResponseWriter, r *http.Request) {
	type eror struct {
		ErrorMessage string
		StatusCode   int
	}
	templ, err := parsefile("front-end/error.html")
	if err != nil {
		log.Fatal("internale server error")
	}
	if r.Method != http.MethodGet {
		templ.Execute(w, eror{ErrorMessage: "Method Not Allowed", StatusCode: http.StatusMethodNotAllowed})
		return
	}

	path := r.URL.Path[1:]
	fileinfo, err := os.Stat(path)
	if err != nil || fileinfo.IsDir() {
		logger.LogWithDetails(err)

		templ.Execute(w, eror{ErrorMessage: "Page Not Found", StatusCode: http.StatusNotFound})
		return
	}
	http.ServeFile(w, r, path)
}
