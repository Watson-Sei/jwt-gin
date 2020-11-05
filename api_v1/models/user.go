package models

type UserModel struct {
	Username		string		`form:"username" binding:"required" gorm:"unique;not null"`
	Password		string		`form:"password" binding:"required" gorm:"not null"`
}