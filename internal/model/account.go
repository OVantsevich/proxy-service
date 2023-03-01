// Package model account model
package model

import "time"

// Account model
type Account struct {
	ID      string    `json:"id"`
	User    string    `json:"user"`
	Amount  float64   `json:"amount"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}
