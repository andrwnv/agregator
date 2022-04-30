package dto

type BaseUserInfo struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Email      string `json:"email"`
	BirthDay   int    `json:"birth_day,omitempty"`
	Verified   bool   `json:"verified"`
	PhotoUrl   string `json:"photo_url"`
}

type CreateUser struct {
	FirstName  string `json:"first_name" binding:"required"`
	SecondName string `json:"second_name" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type UpdateUser struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	BirthDay   int    `json:"birth_day,omitempty"`
	Password   string `json:"password"`
	PhotoUrl   string `json:"photo_url"`
}
