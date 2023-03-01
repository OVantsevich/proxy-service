// Package model User model
package model

import "time"

// User model info
// @Description User account information
type User struct {
	ID       string    `json:"id"`
	Login    string    `json:"login" validate:"required,alphanum,gte=5,lte=20"`
	Email    string    `json:"email" validate:"required,email" format:"email"`
	Password string    `json:"password" validate:"required"`
	Name     string    `json:"name" validate:"required,alpha,gte=2,lte=25"`
	Age      int32     `json:"age" validate:"required,gte=0,lte=100"`
	Created  time.Time `json:"created"`
}
