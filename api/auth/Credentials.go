package auth

// Credentials is data used for login
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
