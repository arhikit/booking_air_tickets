package terr

import (
	"net/http"
)

// Unauthorized represents unauthorized error.
func Unauthorized() *Error {
	return &Error{
		Code:           "UNAUTHORIZED",
		HTTPStatusCode: http.StatusUnauthorized,
	}
}

// Forbidden represents forbidden error.
func Forbidden() *Error {
	return &Error{
		Code:           "FORBIDDEN",
		HTTPStatusCode: http.StatusForbidden,
	}
}

// Conflict represents conflict error.
func Conflict(code string, message string) *Error {
	return &Error{
		Code:           code,
		HTTPStatusCode: http.StatusConflict,
		Message:        message,
	}
}

// NotFound represents not found error.
func NotFound(message string) *Error {
	return &Error{
		Code:           "NOT_FOUND",
		HTTPStatusCode: http.StatusNotFound,
		Message:        message,
	}
}

// InternalServerError represents internal server error.
func InternalServerError(code, message string) *Error {
	return &Error{
		Code:           code,
		HTTPStatusCode: http.StatusInternalServerError,
		Message:        message,
	}
}

// BadRequest represents bad request with invalid input data.
func BadRequest(code, message string) *Error {
	return &Error{
		Code:           code,
		HTTPStatusCode: http.StatusBadRequest,
		Message:        message,
	}
}

// SQLDatabaseError represents sql database error.
func SQLDatabaseError(err error) *Error {
	return &Error{
		Code:           `SQL_DATABASE_ERROR`,
		HTTPStatusCode: http.StatusInternalServerError,
		Message:        err.Error(),
	}
}
