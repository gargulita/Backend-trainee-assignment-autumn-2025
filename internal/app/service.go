package app

import (
    "backend-trainee-assignment/internal/domain"
    "context"
    "errors"
    "math/rand"
    "time"
)

type Service struct {
    store Store
    rand  *rand.Rand
}

func NewService(store Store, r *rand.Rand) *Service {
    if r == nil {
        r = rand.New(rand.NewSource(time.Now().UnixNano()))
    }
    return &Service{store: store, rand: r}
}


func (s *Service) CreateTeam(ctx context.Context, teamName string, members []TeamMemberInput) (*TeamWithMembers, error) {
    if teamName == "" {
        return nil, errors.New("teamName is empty")
    }

    users := make([]*domain.User, 0, len(members))
    for _, m := range members {
        if m.UserID == "" || m.Username == "" {
            return nil, errors.New("user_id and username are required")
        }

        users = append(users, &domain.User{
            ID:       m.UserID,
            Username: m.Username,
            TeamName: teamName,
            IsActive: m.IsActive,
        })
    }

    created := s.store.CreateTeam(ctx, teamName, users)
    if !created {
        return nil, NewAppError(ErrorCodeTeamExists, "team_name already exists")
    }

    teamMembers, _ := s.store.ListUsersByTeam(ctx, teamName)
    return &TeamWithMembers{Name: teamName, Members: teamMembers}, nil
}

func (s *Service) GetTeam(ctx context.Context, teamName string) (*TeamWithMembers, error) {
    if !s.store.TeamExists(ctx, teamName) {
        return nil, NewAppError(ErrorCodeNotFound, "resource not found")
    }

    members, _ := s.store.ListUsersByTeam(ctx, teamName)
    return &TeamWithMembers{Name: teamName, Members: members}, nil
}

func (s *Service) SetUserIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error) {
    user, ok := s.store.SetUserIsActive(ctx, userID, isActive)
    if !ok {
        return nil, NewAppError(ErrorCodeNotFound, "resource not found")
    }
    return user, nil
}


func (s *Service) CreatePullRequest(ctx context.Context, id, name, authorID string) (*domain.PullRequest, error) {
    if id == "" || name == "" || authorID == "" {
        return nil, errors.New("pull_request_id, pull_request_name and author_id are required")
    }

    author, ok := s.store.GetUserByID(ctx, authorID)
    if !ok {
        return nil, NewAppError(ErrorCodeNotFound, "resource not found")
    }

    allMembers, _ := s.store.ListUsersByTeam(ctx, author.TeamName)

    candidates := make([]*domain.User, 0)
    for _, u := range allMembers {
        if u.IsActive && u.ID != author.ID {
            candidates = append(candidates, u)
        }
    }

    reviewers := s.pickRandomReviewers(candidates, 2)

    pr := &domain.PullRequest{
        ID:                id,
        Name:              name,
        AuthorID:          authorID,
        Status:            domain.StatusOpen,
        AssignedReviewers: reviewers,
    }

    if !s.store.CreatePullRequest(ctx, pr) {
        return nil, NewAppError(ErrorCodePRExists, "PR id already exists")
    }

    return pr, nil
}

func (s *Service) GetUserReviewPullRequests(ctx context.Context, userID string) []*domain.PullRequest {
    all := s.store.ListPullRequests(ctx)
    res := make([]*domain.PullRequest, 0)

    for _, pr := range all {
        if containsString(pr.AssignedReviewers, userID) {
            res = append(res, pr)
        }
    }

    return res
}


func (s *Service) MergePullRequest(ctx context.Context, id string) (*domain.PullRequest, error) {
    pr, ok := s.store.GetPullRequestByID(ctx, id)
    if !ok {
        return nil, NewAppError(ErrorCodeNotFound, "resource not found")
    }

    if pr.Status != domain.StatusMerged {
        now := time.Now().UTC()
        pr.Status = domain.StatusMerged
        pr.MergedAt = &now
        s.store.UpdatePullRequest(ctx, pr)
    }

    return pr, nil
}

func (s *Service) ReassignReviewer(ctx context.Context, prID, oldUserID string) (*domain.PullRequest, string, error) {
    if prID == "" || oldUserID == "" {
        return nil, "", errors.New("pull_request_id and old_user_id are required")
    }

    pr, ok := s.store.GetPullRequestByID(ctx, prID)
    if !ok {
        return nil, "", NewAppError(ErrorCodeNotFound, "resource not found")
    }

    if pr.Status == domain.StatusMerged {
        return nil, "", NewAppError(ErrorCodePRMerged, "cannot reassign on merged PR")
    }

    reviewerIndex := -1
    for i, id := range pr.AssignedReviewers {
        if id == oldUserID {
            reviewerIndex = i
            break
        }
    }
    if reviewerIndex == -1 {
        return nil, "", NewAppError(ErrorCodeNotAssigned, "reviewer is not assigned to this PR")
    }

    reviewer, ok := s.store.GetUserByID(ctx, oldUserID)
    if !ok {
        return nil, "", NewAppError(ErrorCodeNotFound, "reviewer not found")
    }

    allMembers, _ := s.store.ListUsersByTeam(ctx, reviewer.TeamName)

    candidates := make([]*domain.User, 0)
    for _, u := range allMembers {
        if !u.IsActive {
            continue
        }
        if u.ID == oldUserID {
            continue
        }
        if containsString(pr.AssignedReviewers, u.ID) {
            continue
        }
        candidates = append(candidates, u)
    }

    if len(candidates) == 0 {
        return nil, "", NewAppError(ErrorCodeNoCandidate, "no active replacement candidate in team")
    }

    newReviewer := s.pickRandomReviewers(candidates, 1)[0]
    pr.AssignedReviewers[reviewerIndex] = newReviewer
    s.store.UpdatePullRequest(ctx, pr)

    return pr, newReviewer, nil
}


func (s *Service) pickRandomReviewers(candidates []*domain.User, limit int) []string {
    if len(candidates) == 0 || limit <= 0 {
        return nil
    }

    if len(candidates) <= limit {
        res := make([]string, 0, len(candidates))
        for _, u := range candidates {
            res = append(res, u.ID)
        }
        return res
    }

    idx := s.rand.Perm(len(candidates))[:limit]

    res := make([]string, 0, limit)
    for _, i := range idx {
        res = append(res, candidates[i].ID)
    }
    return res
}

func containsString(list []string, target string) bool {
    for _, v := range list {
        if v == target {
            return true
        }
    }
    return false
}
