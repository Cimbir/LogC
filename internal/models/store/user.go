package models

type User struct {
	Id       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	IsAdmin  bool   `json:"is_admin" db:"is_admin"`
}
