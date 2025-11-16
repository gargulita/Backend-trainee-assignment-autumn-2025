package httpapi

import (
	"backend-trainee-assignment/internal/app"
	"encoding/json"
	"net/http"
	"strings"
)

type teamMemberDTO struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type teamAddRequest struct {
	TeamName string          `json:"team_name"`
	Members  []teamMemberDTO `json:"members"`
}

type teamResponse struct {
	TeamName string          `json:"team_name"`
	Members  []teamMemberDTO `json:"members"`
}

type teamDeactivateRequest struct {
    TeamName string `json:"team_name"`
}

type teamDeactivateResponse struct {
    TeamName              string   `json:"team_name"`
    DeactivatedUserIDs    []string `json:"deactivated_user_ids"`
    UpdatedPullRequestIDs []string `json:"updated_pull_request_ids"`
}


func (h *Handler) handleTeamAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, http.MethodPost)
		return
	}

	var req teamAddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: errorBody{
				Code:    "BAD_REQUEST",
				Message: "invalid JSON",
			},
		})
		return
	}

	req.TeamName = strings.TrimSpace(req.TeamName)
	if req.TeamName == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: errorBody{
				Code:    "BAD_REQUEST",
				Message: "team_name is required",
			},
		})
		return
	}

	inputMembers := make([]app.TeamMemberInput, 0, len(req.Members))
	for _, m := range req.Members {
		m.UserID = strings.TrimSpace(m.UserID)
		m.Username = strings.TrimSpace(m.Username)
		if m.UserID == "" || m.Username == "" {
			writeJSON(w, http.StatusBadRequest, errorResponse{
				Error: errorBody{
					Code:    "BAD_REQUEST",
					Message: "user_id and username are required",
				},
			})
			return
		}
		inputMembers = append(inputMembers, app.TeamMemberInput{
			UserID:   m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		})
	}

	team, err := h.svc.CreateTeam(r.Context(), req.TeamName, inputMembers)
	if err != nil {
		writeAppError(w, err)
		return
	}

	resp := teamResponse{
		TeamName: team.Name,
		Members:  make([]teamMemberDTO, 0, len(team.Members)),
	}
	for _, u := range team.Members {
		resp.Members = append(resp.Members, teamMemberDTO{
			UserID:   u.ID,
			Username: u.Username,
			IsActive: u.IsActive,
		})
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"team": resp,
	})
}

func (h *Handler) handleTeamGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, http.MethodGet)
		return
	}

	teamName := strings.TrimSpace(r.URL.Query().Get("team_name"))
	if teamName == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: errorBody{
				Code:    "BAD_REQUEST",
				Message: "team_name query param is required",
			},
		})
		return
	}

	team, err := h.svc.GetTeam(r.Context(), teamName)
	if err != nil {
		writeAppError(w, err)
		return
	}

	resp := teamResponse{
		TeamName: team.Name,
		Members:  make([]teamMemberDTO, 0, len(team.Members)),
	}
	for _, u := range team.Members {
		resp.Members = append(resp.Members, teamMemberDTO{
			UserID:   u.ID,
			Username: u.Username,
			IsActive: u.IsActive,
		})
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) handleTeamDeactivate(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        methodNotAllowed(w, http.MethodPost)
        return
    }

    var req teamDeactivateRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSON(w, http.StatusBadRequest, errorResponse{
            Error: errorBody{
                Code:    "BAD_REQUEST",
                Message: "invalid JSON",
            },
        })
        return
    }

    req.TeamName = strings.TrimSpace(req.TeamName)
    if req.TeamName == "" {
        writeJSON(w, http.StatusBadRequest, errorResponse{
            Error: errorBody{
                Code:    "BAD_REQUEST",
                Message: "team_name is required",
            },
        })
        return
    }

    res, err := h.svc.DeactivateTeamUsersAndReassignOpenPRs(r.Context(), req.TeamName)
    if err != nil {
        writeAppError(w, err)
        return
    }

    resp := teamDeactivateResponse{
        TeamName:              res.TeamName,
        DeactivatedUserIDs:    append([]string(nil), res.DeactivatedUserIDs...),
        UpdatedPullRequestIDs: append([]string(nil), res.UpdatedPullRequestIDs...),
    }

    writeJSON(w, http.StatusOK, resp)
}
