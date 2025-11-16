package httpapi

import (
	"backend-trainee-assignment/internal/app"
	"encoding/json"
	"errors"
	"net/http"
)

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error errorBody `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(v)
}

func writeAppError(w http.ResponseWriter, err error) {
	var appErr *app.AppError
	if errors.As(err, &appErr) {
		status := httpStatusFromCode(appErr.Code)
		resp := errorResponse{
			Error: errorBody{
				Code:    string(appErr.Code),
				Message: appErr.Message,
			},
		}
		writeJSON(w, status, resp)
		return
	}

	writeJSON(w, http.StatusInternalServerError, errorResponse{
		Error: errorBody{
			Code:    "INTERNAL_ERROR",
			Message: "internal server error",
		},
	})
}

func httpStatusFromCode(code app.ErrorCode) int {
	switch code {
	case app.ErrorCodeTeamExists:
		return http.StatusBadRequest
	case app.ErrorCodePRExists:
		return http.StatusConflict
	case app.ErrorCodePRMerged,
		app.ErrorCodeNotAssigned,
		app.ErrorCodeNoCandidate:
		return http.StatusConflict
	case app.ErrorCodeNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func methodNotAllowed(w http.ResponseWriter, allowed string) {
	w.Header().Set("Allow", allowed)
	writeJSON(w, http.StatusMethodNotAllowed, errorResponse{
		Error: errorBody{
			Code:    "METHOD_NOT_ALLOWED",
			Message: "method not allowed",
		},
	})
}
