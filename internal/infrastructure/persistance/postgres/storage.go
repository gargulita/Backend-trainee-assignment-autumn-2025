package postgres

import (
	"backend-trainee-assignment/internal/domain"
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type PostgresStore struct {
    db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
    return &PostgresStore{db: db}
}

func (s *PostgresStore) CreateTeam(ctx context.Context, name string, members []*domain.User) bool {
   tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false
	}

	defer func() {
		_ = tx.Rollback() 
	}()

    var exists bool
    err = tx.QueryRowContext(ctx, `SELECT TRUE FROM teams WHERE name=$1`, name).Scan(&exists)
    if err == nil {
        return false
    }

    _, err = tx.ExecContext(ctx, `INSERT INTO teams (name) VALUES ($1)`, name)
    if err != nil {
        return false
    }

    for _, u := range members {
        _, err = tx.ExecContext(ctx,
            `INSERT INTO users (id, username, team_name, is_active)
             VALUES ($1,$2,$3,$4)
             ON CONFLICT (id) DO UPDATE 
                SET username=EXCLUDED.username,
                    team_name=EXCLUDED.team_name,
                    is_active=EXCLUDED.is_active`,
            u.ID, u.Username, name, u.IsActive,
        )
        if err != nil {
            return false
        }
    }

    return tx.Commit() == nil
}

func (s *PostgresStore) TeamExists(ctx context.Context, name string) bool {
    var exists bool
    err := s.db.QueryRowContext(ctx,
        `SELECT TRUE FROM teams WHERE name=$1`, name).Scan(&exists)
    return err == nil
}

func (s *PostgresStore) ListUsersByTeam(ctx context.Context, teamName string) ([]*domain.User, bool) {
    rows, err := s.db.QueryContext(ctx,
        `SELECT id, username, team_name, is_active 
           FROM users WHERE team_name=$1`, teamName)
    if err != nil {
        return nil, false
    }
    defer rows.Close()

    users := []*domain.User{}
    for rows.Next() {
        u := domain.User{}
        if err := rows.Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
            return nil, false
        }
        users = append(users, &u)
    }

    return users, true
}

func (s *PostgresStore) GetUserByID(ctx context.Context, id string) (*domain.User, bool) {
    u := domain.User{}
    err := s.db.QueryRowContext(ctx,
        `SELECT id, username, team_name, is_active FROM users WHERE id=$1`, id).
        Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive)
    if err != nil {
        return nil, false
    }
    return &u, true
}

func (s *PostgresStore) SaveUser(ctx context.Context, user *domain.User) {
    _, _ = s.db.ExecContext(ctx,
        `INSERT INTO users (id, username, team_name, is_active)
         VALUES ($1,$2,$3,$4)
         ON CONFLICT (id) DO UPDATE 
            SET username=EXCLUDED.username,
                team_name=EXCLUDED.team_name,
                is_active=EXCLUDED.is_active`,
        user.ID, user.Username, user.TeamName, user.IsActive)
}

func (s *PostgresStore) SetUserIsActive(ctx context.Context, id string, isActive bool) (*domain.User, bool) {
    _, err := s.db.ExecContext(ctx,
        `UPDATE users SET is_active=$1 WHERE id=$2`, isActive, id)
    if err != nil {
        return nil, false
    }
    return s.GetUserByID(ctx, id)
}

func (s *PostgresStore) CreatePullRequest(ctx context.Context, pr *domain.PullRequest) bool {
    reviewers := pr.AssignedReviewers
    if reviewers == nil {
        reviewers = []string{}
    }

    _, err := s.db.ExecContext(ctx,
        `INSERT INTO pull_requests (id, name, author_id, status, reviewers)
         VALUES ($1,$2,$3,$4,$5)`,
        pr.ID, pr.Name, pr.AuthorID, pr.Status, pq.StringArray(reviewers),
    )
    return err == nil
}

func (s *PostgresStore) GetPullRequestByID(ctx context.Context, id string) (*domain.PullRequest, bool) {
    pr := domain.PullRequest{}
    var reviewers pq.StringArray
    var mergedAt sql.NullTime

    err := s.db.QueryRowContext(ctx,
        `SELECT id, name, author_id, status, reviewers, merged_at 
         FROM pull_requests WHERE id=$1`,
        id,
    ).Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &reviewers, &mergedAt)

    if err != nil {
        return nil, false
    }

    pr.AssignedReviewers = reviewers
    if mergedAt.Valid {
        pr.MergedAt = &mergedAt.Time
    }
    return &pr, true
}

func (s *PostgresStore) UpdatePullRequest(ctx context.Context, pr *domain.PullRequest) bool {
    _, err := s.db.ExecContext(ctx,
        `UPDATE pull_requests 
            SET name=$2, author_id=$3, status=$4, reviewers=$5, merged_at=$6 
          WHERE id=$1`,
        pr.ID, pr.Name, pr.AuthorID, pr.Status,
        pq.StringArray(pr.AssignedReviewers), pr.MergedAt,
    )
    return err == nil
}

func (s *PostgresStore) ListPullRequests(ctx context.Context) []*domain.PullRequest {
    rows, err := s.db.QueryContext(ctx,
        `SELECT id, name, author_id, status, reviewers, merged_at FROM pull_requests`)
    if err != nil {
        return nil
    }
    defer rows.Close()

    list := []*domain.PullRequest{}
    for rows.Next() {
        pr := domain.PullRequest{}
        var reviewers pq.StringArray
        var mergedAt sql.NullTime
        _ = rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &reviewers, &mergedAt)
        pr.AssignedReviewers = reviewers
        if mergedAt.Valid {
            pr.MergedAt = &mergedAt.Time
        }
        list = append(list, &pr)
    }
    return list
}
