package models

type UserModel struct {
	ID 				uint		`gorm:"not null; primaryKey; AUTO_INCREMENT;"`
	Username		string		`json:"username" gorm:"not null" binding:"required"`
	Password		string		`json:"password" gorm:"not null" binding:"required"`
}