package models

type UserRepository interface {
	FindByUsernameAndPassword(dbName, username, encodedPassword string) (*User, error)
	RegisterUser(dbName string, user *RegisterRequest) error
	VerifyUser(dbName, email, securityCode string) (bool, error)
	ResendVerifyCode(dbName, email, securityCode string) error
	SoftDeleteUser(dbName, username, password string) error
	HardDeleteUser(dbName, username, password string) error
}
