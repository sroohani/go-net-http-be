package session

import (
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// func printUsers() {
// 	db.RLock()
// 	fmt.Println("Users:")
// 	for k, v := range db.usersMap {
// 		fmt.Println("----------------")
// 		fmt.Printf("User ID: %v\n", k)
// 		fmt.Printf("Email: %v\n", v.Email)
// 		fmt.Printf("Session Token: %v\n", v.SessionToken)
// 		fmt.Println("----------------")
// 	}
// 	db.RUnlock()
// }

func handleSignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body CredentialsRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil || !isEmailValid(body.Email) || !isPasswordValid(body.Password) {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, JsonMap{"message": "Malformed sigup data"})
		return
	}

	userId, _, _ := findUser(body.Email)
	if userId != "" {
		w.WriteHeader(http.StatusConflict)
		writeJson(w, JsonMap{"message": "User already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), db.bcryptCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, JsonMap{"message": "Internal server error"})
		return
	}

	registerUser(body.Email, string(hashedPassword))
	w.WriteHeader(http.StatusCreated)
	writeJson(w, JsonMap{"message": "Signup successful"})

	// printUsers()
}

func handleLogIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body CredentialsRequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil || !isEmailValid(body.Email) || !isPasswordValid(body.Password) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		writeJson(w, JsonMap{"message": "Wrong credentials"})
		return
	}

	userId, existingUser, err := findUser(body.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		writeJson(w, JsonMap{"message": "User not found"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(existingUser.Password.password), []byte(body.Password)) != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		writeJson(w, JsonMap{"message": "Wrong credentials"})
		return
	}

	sessionToken := generateSessionToken(userId)
	setSession(userId, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		MaxAge:   24 * 3600,
		Path:     "/",
		HttpOnly: true,
		// Secure:   true, // HTTPS only
		SameSite: http.SameSiteStrictMode,
	})
	w.WriteHeader(http.StatusOK)
	writeJson(w, JsonMap{
		"message": "Login successful",
	})

	// printUsers()
}

func handleLogOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var sessionToken string

	cookie, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			writeJson(w, JsonMap{
				"message": "Unauthorized",
			})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			writeJson(w, JsonMap{
				"message": "Internal server error",
			})
		}
		return
	} else {
		sessionToken = cookie.Value
	}

	// Session token still needs sanity check

	isValid := isSessionValid(sessionToken)
	if !isValid {
		w.WriteHeader(http.StatusUnprocessableEntity)
		writeJson(w, JsonMap{"message": "Unauthorized"})
		return
	}

	userId, err := userIdFromSessionToken(sessionToken)
	if err != nil {
		writeJson(w, JsonMap{
			"message": "Internal server error",
		})
		return
	}
	setSession(userId, "")

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		// Secure:   true, // HTTPS only
		SameSite: http.SameSiteStrictMode,
	})
	w.WriteHeader(http.StatusOK)
	writeJson(w, JsonMap{"message": "Logout successful"})

	// printUsers()
}

func handleDropOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var sessionToken string

	cookie, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			writeJson(w, JsonMap{
				"message": "Unauthorized",
			})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			writeJson(w, JsonMap{
				"message": "Internal server error",
			})
		}
		return
	} else {
		sessionToken = cookie.Value
	}

	// Session token still needs sanity check

	isValid := isSessionValid(sessionToken)
	if !isValid {
		w.WriteHeader(http.StatusUnprocessableEntity)
		writeJson(w, JsonMap{"message": "Unauthorized"})
		return
	}

	userId, err := userIdFromSessionToken(sessionToken)
	if err != nil {
		writeJson(w, JsonMap{
			"message": "Internal server error",
		})
		return
	}

	unregisterUser(userId)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		// Secure:   true, // HTTPS only
		SameSite: http.SameSiteStrictMode,
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Dropout successful"}`))

	// printUsers()
}
