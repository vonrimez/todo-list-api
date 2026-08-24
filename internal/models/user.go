package models

type User struct {
	ID    int
	Name  string
	Email string
	Pass  string
}

type UserRegisterInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Pass  string `json:"password"`
}

type UserLoginInput struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Pass  string `json:"password"`
}
