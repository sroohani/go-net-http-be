package session

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"slices"

	"time"
)

func GenerateSecret() {
	secret := make([]byte, 32)
	rand.Read(secret)
	setSecret(secret)
}

func SetBcryptCost(cost int) {
	setBcryptCost(cost)
}

// <user_id><session_id><expiry_time><signature> -> Base64-encoded
func generateSessionToken(userId string) string {
	var token = []byte(userId) // User ID - 36 bytes (Assumed UUIDv4)

	sessionId := make([]byte, 32)
	rand.Read(sessionId)

	expiryTime := []byte(time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05")) // Expiry time - 19 bytes

	token = append(token, sessionId...)
	token = append(token, expiryTime...)

	mac := hmac.New(sha256.New, db.secret)
	mac.Write(token)
	signature := mac.Sum(nil)
	token = append(token, signature...)

	return base64.URLEncoding.EncodeToString(token)
}

func isSessionValid(token string) bool {
	tokenBytes, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return false
	}

	expiryTime, err := time.Parse("2006-01-02 15:04:05", string(tokenBytes[68:87]))
	if err != nil {
		return false
	}

	receivedSignature := tokenBytes[87:]
	mac := hmac.New(sha256.New, db.secret)
	mac.Write(tokenBytes[:87])
	calculatedSignature := mac.Sum(nil)
	return slices.Equal(receivedSignature, calculatedSignature) && expiryTime.After(time.Now())
}

func userIdFromSessionToken(token string) (string, error) {
	tokenBytes, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}

	return string(tokenBytes[:36]), nil
}
