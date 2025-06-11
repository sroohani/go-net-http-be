package session

import (
	"sync"

	"github.com/google/uuid"
)

var users []*User
var db = struct {
	sync.RWMutex
	secret     []byte
	bcryptCost int
	usersMap   map[string]*User // User ID (UUIDv4) -> User
}{usersMap: make(map[string]*User)}

func findUser(email string) (string, *User, error) {
	db.RLock()
	defer db.RUnlock()
	for k, v := range db.usersMap {
		if email == v.Email {
			return k, v, nil
		}
	}

	return "", nil, ErrorUserNotFound
}

func registerUser(email, password string) *User {
	user := &User{
		Email:    email,
		Password: Password{password: password},
	}

	db.Lock()
	defer db.Unlock()
	db.usersMap[uuid.NewString()] = user

	return user
}

func unregisterUser(userId string) (*User, error) {
	db.Lock()
	defer db.Unlock()
	user, ok := db.usersMap[userId]
	if ok {
		delete(db.usersMap, userId)
		return user, nil
	}
	return nil, ErrorInvalidUser
}

func setSecret(secret []byte) {
	db.Lock()
	defer db.Unlock()
	db.secret = secret
}

func setBcryptCost(cost int) {
	db.Lock()
	defer db.Unlock()
	db.bcryptCost = cost
}

func setSession(userId string, token string) error {
	db.Lock()
	defer db.Unlock()
	_, ok := db.usersMap[userId]
	if !ok {
		return ErrorInvalidUser
	}
	db.usersMap[userId].SessionToken = token

	return nil
}
