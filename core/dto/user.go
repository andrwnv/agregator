package dto

type BaseUserInfo struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Email      string `json:"email"`
	BirthDay   int    `json:"birth_day,omitempty"`
	Verified   bool   `json:"verified"`
}

type CreateUser struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}
