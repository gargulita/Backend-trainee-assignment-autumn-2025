package app

import (
    "backend-trainee-assignment/internal/domain"
    "context"
)

type Store interface {
    CreateTeam(ctx context.Context, name string, members []*domain.User) bool
    TeamExists(ctx context.Context, name string) bool
    ListUsersByTeam(ctx context.Context, teamName string) ([]*domain.User, bool)

    GetUserByID(ctx context.Context, id string) (*domain.User, bool)
    SaveUser(ctx context.Context, user *domain.User)
    SetUserIsActive(ctx context.Context, id string, isActive bool) (*domain.User, bool)

    CreatePullRequest(ctx context.Context, pr *domain.PullRequest) bool
    GetPullRequestByID(ctx context.Context, id string) (*domain.PullRequest, bool)
    UpdatePullRequest(ctx context.Context, pr *domain.PullRequest) bool
    ListPullRequests(ctx context.Context) []*domain.PullRequest

     GetStats(ctx context.Context) (*domain.Stats, error)
}
