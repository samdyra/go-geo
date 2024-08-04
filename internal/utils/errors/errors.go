package errors

import "errors"

var (
    ErrInvalidInput       = errors.New("invalid input")
    ErrUserAlreadyExists  = errors.New("user already exists")
    ErrUserNotFound       = errors.New("user not found")
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrInternalServer     = errors.New("internal server error")
    ErrUnauthorized       = errors.New("unauthorized")
)

type APIError struct {
    Type    string `json:"type"`
    Message string `json:"message"`
}

func NewAPIError(err error) APIError {
    switch err {
    case ErrInvalidInput:
        return APIError{Type: "INVALID_INPUT", Message: err.Error()}
    case ErrUserAlreadyExists:
        return APIError{Type: "USER_ALREADY_EXISTS", Message: err.Error()}
    case ErrUserNotFound:
        return APIError{Type: "USER_NOT_FOUND", Message: err.Error()}
    case ErrInvalidCredentials:
        return APIError{Type: "INVALID_CREDENTIALS", Message: err.Error()}
    case ErrUnauthorized:
        return APIError{Type: "UNAUTHORIZED", Message: err.Error()}
    default:
        return APIError{Type: "INTERNAL_SERVER_ERROR", Message: "An unexpected error occurred"}
    }
}