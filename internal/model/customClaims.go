// Package model custom claims model
package model

import "github.com/golang-jwt/jwt/v4"

// CustomClaims RegisteredClaims with ID
type CustomClaims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}
