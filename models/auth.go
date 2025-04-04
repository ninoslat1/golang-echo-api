// models/auth/auth.go
package models

type AuthService interface {
	Login(dbName string, loginReq *LoginRequest) (*User, error)
	Register(dbName string, registerReq *RegisterRequest) (*RegisterResponse, error)
	VerifyUser(dbName, email, securityCode string) (bool, error)
	SoftDeleteUser(dbName string, loginReq *LoginRequest) error
	HardDeleteUser(dbName string, loginReq *LoginRequest) error
}

type LoginRequest struct {
	UserName string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Cookie  string `json:"cookie,omitempty"`
}

type VerifyRequest struct {
	Email        string `json:"email" form:"email"`
	SecurityCode string `json:"security_code" form:"security_code"`
}

type RegisterRequest struct {
	UserName     string `json:"username" binding:"required" gorm:"column:UserName" form:"username"`
	Password     string `json:"password" binding:"required"  gorm:"column:Password" form:"password"`
	Email        string `json:"email" binding:"required,email"  gorm:"column:Email" form:"email"`
	UserCode     string `json:"usercode" binding:"required"  gorm:"column:UserCode" form:"usercode"`
	LogIn        int64  `json:"login" form:"login"  gorm:"column:LogIn"`
	SecurityCode string `json:"securitycode" form:"securitycode"  gorm:"column:SecurityCode" binding:"required"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}
