package model

type AppError struct {
	Error   error  `json:"-"`
	Message string `json:"message"`
	Code    int    `json:"-"`
}
