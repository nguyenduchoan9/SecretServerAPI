package model

import "time"

type HashKey struct {
	ID             int       `json:"-"`
	Hash           string    `json:"hash"`
	SecretText     string    `json:"secretText"`
	CreatedAt      time.Time `json:"createdAt"`
	ExpireAt       time.Time `json:"expiresAt"`
	RemainingViews int       `json:"remainingViews"`
}
