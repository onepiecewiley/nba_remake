package myErrors

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Detail  string    `json:"detail,omitempty"`
}

func NewError(code ErrorCode, msg string, detail string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Detail:  detail,
	}
}
