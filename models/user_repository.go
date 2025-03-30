package models

type UserRepository interface {
	FindByUsernameAndPassword(dbName, username, encodedPassword string) (*User, error)
	RegisterUser(dbName string, user *RegisterRequest) (bool, error)
	VerifyUser(dbName, email, securityCode string) (bool error)
}
