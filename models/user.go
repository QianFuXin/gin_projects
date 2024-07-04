package models

type User struct {
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"name" binding:"required,oneof=male female prefer_not_to"`
	Email string `json:"email"`
}
