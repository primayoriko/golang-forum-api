package models

import "github.com/dgrijalva/jwt-go"

// Claims is login result as an auth token
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
