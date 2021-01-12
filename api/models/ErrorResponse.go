package models

// type User struct {
// 	ID   uint64 `json:"id" example:"100" description:"User identity"`
// 	Name string `json:"name" example:"Mikun"`
// }

// type UsersResponse struct {
// 	Data []User `json:"users" example:"[{\"id\":100, \"name\":\"Mikun\"}]"`
// }

// ErrorResponse is optional body message for any failure response
type ErrorResponse struct {
	Message string `json:"message" example:"failure to connect to db"`
}

// NewErrorResponse is constructor for create new instance of ErrorResponse
func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{Message: message}
}
