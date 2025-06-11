package session

import "errors"

var (
	ErrorInvalidUser  = errors.New("Invalid user")
	ErrorUserNotFound = errors.New("User not found")
)
