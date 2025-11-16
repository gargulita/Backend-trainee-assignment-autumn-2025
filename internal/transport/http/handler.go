package httpapi

import (
	"backend-trainee-assignment/internal/app"
	"net/http"
)

type Handler struct {
	svc *app.Service
}

func NewHandler(svc *app.Service) http.Handler {
	h := &Handler{svc: svc}

	mux := http.NewServeMux()

	mux.HandleFunc("/team/add", h.handleTeamAdd)
	mux.HandleFunc("/team/get", h.handleTeamGet)
	mux.HandleFunc("/team/deactivate", h.handleTeamDeactivate)

	mux.HandleFunc("/users/setIsActive", h.handleUserSetIsActive)
	mux.HandleFunc("/users/getReview", h.handleUserGetReview)

	mux.HandleFunc("/pullRequest/create", h.handlePullRequestCreate)
	mux.HandleFunc("/pullRequest/merge", h.handlePullRequestMerge)
	mux.HandleFunc("/pullRequest/reassign", h.handlePullRequestReassign)

	mux.HandleFunc("/health", h.handleHealth)

	mux.HandleFunc("/stats", h.handleStats)


	return mux
}
