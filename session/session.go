package session

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"
)

var sessionStore = struct {
	sync.RWMutex
	m map[string]time.Time // This maps session IDs to their respective expiration times
}{m: make(map[string]time.Time)}

func createSession() string {
	sessionID := make([]byte, 32)
	rand.Read(sessionID)
	id := base64.URLEncoding.EncodeToString(sessionID)

	sessionStore.Lock()
	sessionStore.m[id] = time.Now().Add(24 * time.Hour)
	sessionStore.Unlock()

	return id
}

func isSessionValid(sessionID string) bool {
	sessionStore.RLock()
	expiry, exists := sessionStore.m[sessionID]
	sessionStore.RUnlock()

	return exists && time.Now().Before(expiry)
}
