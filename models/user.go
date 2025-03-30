package models

type User struct {
	ID           uint    `gorm:"primaryKey;autoIncrement"`
	UserCode     string  `gorm:"column:UserCode;not null;unique;size:50"`
	UserName     string  `gorm:"column:UserName;not null;size:100"`
	Position     *string `gorm:"column:Position;size:100"`
	Telephone    *string `gorm:"column:Telephone;size:20"`
	Handphone    *string `gorm:"column:Handphone;size:20"`
	Email        *string `gorm:"column:Email;size:100"`
	Password     string  `gorm:"column:Password;not null;size:255"`
	SecurityCode string  `gorm:"column:SecurityCode;not null;size:255"`
	GroupID      uint    `gorm:"column:GroupID;not null"`
	Status       int32   `gorm:"column:Status;not null"`
	UserID       uint    `gorm:"column:UserID;not null"`
	LogIn        int32   `gorm:"column:LogIn;not null"`
}

func (User) TableName() string {
	return "myuser"
}
