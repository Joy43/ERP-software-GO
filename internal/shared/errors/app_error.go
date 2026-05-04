package errors

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

// helpers
func BadRequest(msg string) *AppError {
	return &AppError{Code: 400, Message: msg}
}

func NotFound(msg string) *AppError {
	return &AppError{Code: 404, Message: msg}
}

func Conflict(msg string) *AppError {
	return &AppError{Code: 409, Message: msg}
}

func Internal(msg string) *AppError {
	return &AppError{Code: 500, Message: msg}
}