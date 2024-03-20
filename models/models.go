package models

// User schema of the user table
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Image    string `json:"image"`
	Password string `json:"password"`
}

type UserRegister struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Image    string `json:"image"`
	Password string `json:"password"`
}
type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
