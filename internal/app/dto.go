package app

import "backend-trainee-assignment/internal/domain"

type TeamMemberInput struct {
    UserID   string
    Username string
    IsActive bool
}

type TeamWithMembers struct {
    Name    string
    Members []*domain.User
}

type DeactivateTeamResult struct {
    TeamName              string
    DeactivatedUserIDs    []string
    UpdatedPullRequestIDs []string
}