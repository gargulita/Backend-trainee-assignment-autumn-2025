package memory

import (
	"backend-trainee-assignment/internal/domain"
	"context"
	"sync"
)

type InMemoryStore struct {
	mu sync.RWMutex

	teams        map[string]struct{}
	users        map[string]*domain.User
	pullRequests map[string]*domain.PullRequest
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		teams:        make(map[string]struct{}),
		users:        make(map[string]*domain.User),
		pullRequests: make(map[string]*domain.PullRequest),
	}
}

func (s *InMemoryStore) CreateTeam(_ context.Context, name string, members []*domain.User) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.teams[name]; exists {
		return false
	}

	s.teams[name] = struct{}{}
	for _, u := range members {
		if u == nil {
			continue
		}
		copyUser := *u
		s.users[u.ID] = &copyUser
	}

	return true
}

func (s *InMemoryStore) TeamExists(_ context.Context, name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.teams[name]
	return exists
}

func (s *InMemoryStore) ListUsersByTeam(_ context.Context, teamName string) ([]*domain.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.teams[teamName]; !exists {
		return nil, false
	}

	var res []*domain.User
	for _, u := range s.users {
		if u.TeamName == teamName {
			copyUser := *u
			res = append(res, &copyUser)
		}
	}

	return res, true
}

func (s *InMemoryStore) GetUserByID(_ context.Context, id string) (*domain.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.users[id]
	if !ok {
		return nil, false
	}
	copyUser := *u
	return &copyUser, true
}

func (s *InMemoryStore) SaveUser(_ context.Context, user *domain.User) {
	if user == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	copyUser := *user
	s.users[user.ID] = &copyUser
}

func (s *InMemoryStore) SetUserIsActive(_ context.Context, id string, isActive bool) (*domain.User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.users[id]
	if !ok {
		return nil, false
	}

	u.IsActive = isActive
	copyUser := *u
	return &copyUser, true
}

func (s *InMemoryStore) CreatePullRequest(_ context.Context, pr *domain.PullRequest) bool {
	if pr == nil {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.pullRequests[pr.ID]; exists {
		return false
	}

	copyPR := *pr
	s.pullRequests[pr.ID] = &copyPR
	return true
}

func (s *InMemoryStore) GetPullRequestByID(_ context.Context, id string) (*domain.PullRequest, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pr, ok := s.pullRequests[id]
	if !ok {
		return nil, false
	}
	copyPR := *pr
	if pr.AssignedReviewers != nil {
		copyPR.AssignedReviewers = append([]string(nil), pr.AssignedReviewers...)
	}
	return &copyPR, true
}

func (s *InMemoryStore) UpdatePullRequest(_ context.Context, pr *domain.PullRequest) bool {
	if pr == nil {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.pullRequests[pr.ID]; !exists {
		return false
	}

	copyPR := *pr
	if pr.AssignedReviewers != nil {
		copyPR.AssignedReviewers = append([]string(nil), pr.AssignedReviewers...)
	}
	s.pullRequests[pr.ID] = &copyPR
	return true
}

func (s *InMemoryStore) ListPullRequests(_ context.Context) []*domain.PullRequest {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*domain.PullRequest, 0, len(s.pullRequests))
	for _, pr := range s.pullRequests {
		copyPR := *pr
		if pr.AssignedReviewers != nil {
			copyPR.AssignedReviewers = append([]string(nil), pr.AssignedReviewers...)
		}
		res = append(res, &copyPR)
	}
	return res
}
