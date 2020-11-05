package models

type UserModel struct {
	Username	string	`gorm:"size:50;not null;" json:"username" form:"username" binding:"required"`
	Password	string	`gorm:"not null;" json:"password" form:"password" binding:"password"`
}

func (b *UserModel) TableName() string {
	return "usermodel"
}
