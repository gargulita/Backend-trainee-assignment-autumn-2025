package app

type ErrorCode string

const (
    ErrorCodeTeamExists  ErrorCode = "TEAM_EXISTS"
    ErrorCodePRExists    ErrorCode = "PR_EXISTS"
    ErrorCodePRMerged    ErrorCode = "PR_MERGED"
    ErrorCodeNotAssigned ErrorCode = "NOT_ASSIGNED"
    ErrorCodeNoCandidate ErrorCode = "NO_CANDIDATE"
    ErrorCodeNotFound    ErrorCode = "NOT_FOUND"
    ErrorCodeBadRequest  ErrorCode = "BAD_REQUEST"
)

type AppError struct {
    Code    ErrorCode
    Message string
}

func (e *AppError) Error() string {
    return e.Message
}

func NewAppError(code ErrorCode, msg string) *AppError {
    return &AppError{
        Code:    code,
        Message: msg,
    }
}
