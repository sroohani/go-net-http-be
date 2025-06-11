package session

import (
	"encoding/json"
	"io"
)

// Custom email validation logic
func isEmailValid(email string) bool {
	return email != ""
}

// Custom password validation logic
func isPasswordValid(password string) bool {
	return password != ""
}

func writeJson(w io.Writer, m JsonMap) {
	json.NewEncoder(w).Encode(m)
}
