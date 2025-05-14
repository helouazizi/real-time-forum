package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"web-forum/internal/models"
	"web-forum/internal/services"
	"web-forum/internal/utils"
	"web-forum/pkg/logger"
)

type PostHandler struct {
	PostService *services.PostService
}

func NewPostHandler(PostService *services.PostService) *PostHandler {
	return &PostHandler{PostService: PostService}
}
func (h *PostHandler) FetchPosts(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		utils.RespondWithJSON(w, http.StatusMethodNotAllowed, models.Error{Message: "Method Not Allowed", Code: http.StatusMethodNotAllowed})
		return
	}
	Posts, err := h.PostService.FetchPosts()
	if err.Code != http.StatusOK {
		utils.RespondWithJSON(w, http.StatusInternalServerError, models.Error{
			Message: "Internal server error",
			Code:    http.StatusInternalServerError,
		})
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, Posts)
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithJSON(w, http.StatusMethodNotAllowed, models.Error{Message: "Methos Not Allowed", Code: http.StatusMethodNotAllowed})
		return
	}
	var Post models.Post
	if err := json.NewDecoder(r.Body).Decode(&Post); err != nil {
		logger.LogWithDetails(err)
		utils.RespondWithJSON(w, http.StatusBadRequest, models.Error{Message: "Bad Request", Code: http.StatusBadRequest})
		return
	}
	// lets validate the post input
	err := utils.ValidPostInputs(Post)
	if err.UserErrors.HasError {
		logger.LogWithDetails(fmt.Errorf("invalid post input"))
		utils.RespondWithJSON(w, http.StatusBadRequest, err)
		return
	}

	// lets get the user id
	token, err := utils.GetToken(r, "Token")
	if err.Code != http.StatusOK {
		utils.RespondWithJSON(w, err.Code, err)
		return
	}

	userId, err := h.PostService.GetUserID(token)
	if err.Code != http.StatusOK {
		logger.LogWithDetails(fmt.Errorf(err.Message))
		utils.RespondWithJSON(w, err.Code, err)
		return
	}
	Post.UserID = userId

	err1 := h.PostService.CreatePost(Post)
	if err1.Code != http.StatusCreated {
		logger.LogWithDetails(fmt.Errorf(err.Message))
		// utils.RespondWithError(w, err)
		return
	}
	// our response
	utils.RespondWithJSON(w, http.StatusCreated, models.SuccesMessage{Message: "Post created successfully"})
}

func (h *PostHandler) ReactToPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithJSON(w, http.StatusMethodNotAllowed, models.Error{
			Message: "Method Not Allowed",
			Code:    http.StatusMethodNotAllowed,
		})
		return
	}

	// Parse and decode JSON body
	var reaction models.PostReaction
	if err := json.NewDecoder(r.Body).Decode(&reaction); err != nil {
		logger.LogWithDetails(err)
		utils.RespondWithJSON(w, http.StatusBadRequest, models.Error{
			Message: "Bad Request",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate reaction type
	if reaction.Reaction != "like" && reaction.Reaction != "dislike" {
		utils.RespondWithJSON(w, http.StatusBadRequest, models.Error{
			Message: "Bad Request",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Get token from cookie
	token, err := utils.GetToken(r, "Token")
	if err.Code != http.StatusOK {
		logger.LogWithDetails(fmt.Errorf(err.Message))
		utils.RespondWithJSON(w, err.Code, err)
		return
	}

	// Call service to react
	serviceErr := h.PostService.ReactToPost(token, reaction)
	if serviceErr.Code != http.StatusOK {
		logger.LogWithDetails(fmt.Errorf(serviceErr.Message))
		utils.RespondWithJSON(w, serviceErr.Code, serviceErr)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.SuccesMessage{
		Message: "Reaction recorded",
	})
}

func (h *PostHandler) CommentPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithJSON(w, http.StatusMethodNotAllowed, models.Error{
			Message: "Method Not Allowed",
			Code:    http.StatusMethodNotAllowed,
		})
		return
	}

	// Parse and decode JSON body
	var reaction models.PostReaction
	if err := json.NewDecoder(r.Body).Decode(&reaction); err != nil {
		logger.LogWithDetails(err)
		utils.RespondWithJSON(w, http.StatusBadRequest, models.Error{
			Message: "Bad Request",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate reaction type
	if reaction.Comment == "" || len(strings.Fields(reaction.Comment)) == 0 {
		utils.RespondWithJSON(w, http.StatusBadRequest, models.Error{
			Message: "Bad Request",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Get token from cookie
	token, err := utils.GetToken(r, "Token")
	if err.Code != http.StatusOK {
		logger.LogWithDetails(fmt.Errorf(err.Message))
		utils.RespondWithJSON(w, err.Code, err)
		return
	}

	// Call service to react
	serviceErr := h.PostService.AddComment(token, reaction)
	if serviceErr.Code != http.StatusOK {
		logger.LogWithDetails(fmt.Errorf(serviceErr.Message))
		utils.RespondWithJSON(w, serviceErr.Code, serviceErr)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.SuccesMessage{
		Message: "Comment Posted seccefully",
	})
}

func (h *PostHandler) FetchComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.RespondWithJSON(w, http.StatusMethodNotAllowed, models.Error{
			Message: "Method Not Allowed",
			Code:    http.StatusMethodNotAllowed,
		})
		return
	}
	postId := r.URL.Query().Get("postId")
	if postId == "" {
		utils.RespondWithJSON(w, http.StatusBadRequest, models.Error{
			Message: "Missing postId parameter",
			Code:    http.StatusBadRequest,
		})
		return
	}

	postID, err := strconv.Atoi(postId)
	if err != nil {
		utils.RespondWithJSON(w, http.StatusInternalServerError, models.Error{
			Message: "Failed to fetch comments",
			Code:    http.StatusBadRequest,
		})
		return
	}

	comments, err1 := h.PostService.GetCommentsByPostID(postID)
	if err1.Code != http.StatusOK {
		utils.RespondWithJSON(w, http.StatusInternalServerError, models.Error{
			Message: "Failed to fetch comments",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, comments)
}

func (h *PostHandler) FilterPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithJSON(w, http.StatusMethodNotAllowed, models.Error{Message: "Methos Not Allowed", Code: http.StatusMethodNotAllowed})
		return
	}
	var Post models.Post
	if err := json.NewDecoder(r.Body).Decode(&Post); err != nil {
		logger.LogWithDetails(err)
		utils.RespondWithJSON(w, http.StatusBadRequest, models.Error{Message: "Bad Request", Code: http.StatusBadRequest})
		return
	}
	filtredPosts, err1 := h.PostService.FilterPosts(Post.Categories)
	if err1.Code != http.StatusOK {
		logger.LogWithDetails(fmt.Errorf(err1.Message))
		utils.RespondWithJSON(w, err1.Code, err1)
		return
	}
	// our response
	utils.RespondWithJSON(w, http.StatusOK, filtredPosts)
}
