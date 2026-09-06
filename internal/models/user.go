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
	Name  *string `json:"name" validate:"required,min=2,max=50"`
	Email *string `json:"email" validate:"required,email"`
	Pass  *string `json:"password" validate:"required,min=8,max=255"`
}

type UserLoginInput struct {
	Email *string `json:"email" validate:"required,email"`
	Pass  *string `json:"password" validate:"required"`
}
