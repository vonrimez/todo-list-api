package models

type User struct {
	ID   int
	Name string
	JWT  string
}

type UserOutput struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Pass string `json:"password"`
}

type UserRegisterInput struct {
	Name  *string `json:"name" binding:"required"`
	Email *string `json:"email" binding:"required"`
	Pass  *string `json:"password" binding:"required"`
}

type UserLoginInput struct {
	Email *string `json:"email" binding:"required"`
	Pass  *string `json:"password" binding:"required"`
}
