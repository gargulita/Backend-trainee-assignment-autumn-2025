package httpapi

import (
    "net/http"
)

type statsResponse struct {
    ReviewAssignments map[string]int `json:"review_assignments"`
    PRStatuses        map[string]int `json:"pr_statuses"`
}

func (h *Handler) handleStats(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        methodNotAllowed(w, http.MethodGet)
        return
    }

    stats, err := h.svc.GetStats(r.Context())
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, errorResponse{
            Error: errorBody{
                Code:    "INTERNAL_ERROR",
                Message: err.Error(),
            },
        })
        return
    }

    statuses := make(map[string]int)
    for st, cnt := range stats.PRStatuses {
        statuses[string(st)] = cnt
    }

    resp := statsResponse{
        ReviewAssignments: stats.ReviewAssignments,
        PRStatuses:        statuses,
    }

    writeJSON(w, http.StatusOK, resp)
}
