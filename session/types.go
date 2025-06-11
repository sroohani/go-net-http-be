package session

type Password struct {
	password string
}

func (p Password) MarshalJSON() ([]byte, error) {
	return []byte{}, nil
}

type User struct {
	Email        string   `json:"email"`
	Password     Password `json:"password"`
	SessionToken string   `json:"-"`
}

type CredentialsRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Logout struct {
	SessionToken string `json:"session_token"`
}

type JsonMap = map[string]string
