package httpapi

import (
	"backend-trainee-assignment/internal/domain"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type prCreateRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type prMergeRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type prReassignRequest struct {
	PullRequestID       string `json:"pull_request_id"`
	OldUserID           string `json:"old_user_id"`        
	LegacyOldReviewerID string `json:"old_reviewer_id"`   
}

type prDTO struct {
	PullRequestID   string     `json:"pull_request_id"`
	PullRequestName string     `json:"pull_request_name"`
	AuthorID        string     `json:"author_id"`
	Status          string     `json:"status"`
	AssignedReviewers []string `json:"assigned_reviewers"`
	MergedAt        string `json:"mergedAt,omitempty"`
}

type prCreateResponse struct {
	PR prDTO `json:"pr"`
}

type prMergeResponse struct {
	PR prDTO `json:"pr"`
}

type prReassignResponse struct {
	PR         prDTO  `json:"pr"`
	ReplacedBy string `json:"replaced_by"`
}

func (h *Handler) handlePullRequestCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}

	var req prCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: errorBody{
				Code:    "BAD_REQUEST",
				Message: "invalid JSON",
			},
		})
		return
	}

	req.PullRequestID = strings.TrimSpace(req.PullRequestID)
	req.PullRequestName = strings.TrimSpace(req.PullRequestName)
	req.AuthorID = strings.TrimSpace(req.AuthorID)

	if req.PullRequestID == "" || req.PullRequestName == "" || req.AuthorID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: errorBody{
				Code:    "BAD_REQUEST",
				Message: "pull_request_id, pull_request_name and author_id are required",
			},
		})
		return
	}

	pr, err := h.svc.CreatePullRequest(r.Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	resp := prCreateResponse{
		PR: toPRDTO(pr),
	}
	writeJSON(w, http.StatusCreated, resp)
}

func (h *Handler) handlePullRequestMerge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}

	var req prMergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: errorBody{
				Code:    "BAD_REQUEST",
				Message: "invalid JSON",
			},
		})
		return
	}

	req.PullRequestID = strings.TrimSpace(req.PullRequestID)
	if req.PullRequestID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: errorBody{
				Code:    "BAD_REQUEST",
				Message: "pull_request_id is required",
			},
		})
		return
	}

	pr, err := h.svc.MergePullRequest(r.Context(), req.PullRequestID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	resp := prMergeResponse{
		PR: toPRDTO(pr),
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) handlePullRequestReassign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}

	var req prReassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: errorBody{
				Code:    "BAD_REQUEST",
				Message: "invalid JSON",
			},
		})
		return
	}

	req.PullRequestID = strings.TrimSpace(req.PullRequestID)
	req.OldUserID = strings.TrimSpace(req.OldUserID)
	req.LegacyOldReviewerID = strings.TrimSpace(req.LegacyOldReviewerID)

	if req.OldUserID == "" && req.LegacyOldReviewerID != "" {
		req.OldUserID = req.LegacyOldReviewerID
	}

	if req.PullRequestID == "" || req.OldUserID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: errorBody{
				Code:    "BAD_REQUEST",
				Message: "pull_request_id and old_user_id are required",
			},
		})
		return
	}

	pr, replacedBy, err := h.svc.ReassignReviewer(r.Context(), req.PullRequestID, req.OldUserID)
	if err != nil {
		writeAppError(w, err)
		return
	}

	resp := prReassignResponse{
		PR:         toPRDTO(pr),
		ReplacedBy: replacedBy,
	}
	writeJSON(w, http.StatusOK, resp)
}

func toPRDTO(pr *domain.PullRequest) prDTO {
    if pr == nil {
        return prDTO{}
    }

    dto := prDTO{
        PullRequestID:     pr.ID,
        PullRequestName:   pr.Name,
        AuthorID:          pr.AuthorID,
        Status:            string(pr.Status),
        AssignedReviewers: append([]string(nil), pr.AssignedReviewers...),
    }

    if pr.MergedAt != nil {
        dto.MergedAt = pr.MergedAt.Format(time.RFC3339) 
    }

    return dto
}

