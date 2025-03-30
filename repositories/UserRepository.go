package repositories

import (
	"echo-api/configs"
	models "echo-api/models"
	"echo-api/utils"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
)

var localLog = logrus.New()

type UserRepository interface {
	FindByUsernameAndPassword(dbName, username, encodedPassword string) (*models.User, error)
	RegisterUser(dbName string, user *models.RegisterRequest) (bool, error)
	VerifyUser(dbName, email, securityCode string) (bool, error)
}

type userRepository struct {
	log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) UserRepository {
	return &userRepository{log}
}

func (r *userRepository) RegisterUser(dbName string, user *models.RegisterRequest) (bool, error) {
	db, err := configs.RunDatabase(dbName)
	if err != nil {
		localLog.Error("Database connection failed:", err)
		return false, err
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	securityCode := fmt.Sprintf("%06d", rand.Intn(1000000))

	user.SecurityCode = securityCode
	user.LogIn = 0

	err = db.Table("myuser").Create(&user).Error
	if err != nil {
		return false, err
	}

	if user.Email == "" {
		return false, errors.New("Email cannot be null")
	}

	err = utils.SendVerificationEmail(user.Email, securityCode)
	if err != nil {
		message := fmt.Sprintf("Failed to send verification email to %s", user.UserCode)
		r.log.Error(message, err)
		return false, errors.New("Failed to send verification email")
	}

	return true, nil
}

func (r *userRepository) VerifyUser(dbName, email, securityCode string) (bool, error) {
	db, err := configs.RunDatabase(dbName)
	if err != nil {
		localLog.Error("Database connection failed:", err)
		return false, err
	}

	// Update langsung tanpa perlu memanggil First(&user)
	result := db.Table("myuser").
		Where("Email = ? AND SecurityCode = ?", email, securityCode).
		Updates(map[string]interface{}{"LogIn": 1})

	if result.RowsAffected == 0 {
		return false, errors.New("Invalid verification code or user already verified")
	}

	if result.Error != nil {
		return false, errors.New("Failed to verify user")
	}

	return true, nil
}

func (r *userRepository) FindByUsernameAndPassword(dbName, username, encodedPassword string) (*models.User, error) {
	db, err := configs.RunDatabase(dbName)
	if err != nil {
		localLog.Error("Database connection failed:", err)
		return nil, err
	}

	var user models.User
	err = db.Table("myuser").
		Select("UserCode, Login").
		Where("UserName = ? AND Password LIKE ?", username, encodedPassword[:len(encodedPassword)-2]+"%").
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
