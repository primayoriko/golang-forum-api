package models

// ErrorResponse is optional body message for any failure response
type ErrorResponse struct {
	Message string `json:"message" example:"failure to connect to db"`
}

// NewErrorResponse is constructor for create new instance of ErrorResponse
func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{Message: message}
}
