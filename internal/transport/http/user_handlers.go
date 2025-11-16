package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	
)

type setIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type userDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type setIsActiveResponse struct {
	User userDTO `json:"user"`
}

type pullRequestShortDTO struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

type userGetReviewResponse struct {
	UserID       string               `json:"user_id"`
	PullRequests []pullRequestShortDTO `json:"pull_requests"`
}

func (h *Handler) handleUserSetIsActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}

	var req setIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: errorBody{
				Code:    "BAD_REQUEST",
				Message: "invalid JSON",
			},
		})
		return
	}

	req.UserID = strings.TrimSpace(req.UserID)
	if req.UserID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: errorBody{
				Code:    "BAD_REQUEST",
				Message: "user_id is required",
			},
		})
		return
	}

	user, err := h.svc.SetUserIsActive(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		writeAppError(w, err)
		return
	}

	resp := setIsActiveResponse{
		User: userDTO{
			UserID:   user.ID,
			Username: user.Username,
			TeamName: user.TeamName,
			IsActive: user.IsActive,
		},
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) handleUserGetReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}

	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))
	if userID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: errorBody{
				Code:    "BAD_REQUEST",
				Message: "user_id query param is required",
			},
		})
		return
	}

	prs := h.svc.GetUserReviewPullRequests(r.Context(), userID)

	resp := userGetReviewResponse{
		UserID:       userID,
		PullRequests: make([]pullRequestShortDTO, 0, len(prs)),
	}

	for _, pr := range prs {
		resp.PullRequests = append(resp.PullRequests, pullRequestShortDTO{
			PullRequestID:   pr.ID,
			PullRequestName: pr.Name,
			AuthorID:        pr.AuthorID,
			Status:          string(pr.Status),
		})
	}

	writeJSON(w, http.StatusOK, resp)
}
