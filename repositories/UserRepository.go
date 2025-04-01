package repositories

import (
	"echo-api/configs"
	models "echo-api/models"
	"echo-api/utils"
	"errors"
	"fmt"

	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

var localLog = logrus.New()

type userRepository struct {
	log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) models.UserRepository {
	return &userRepository{log}
}

func (r *userRepository) RegisterUser(dbName string, user *models.RegisterRequest) error {
	db, err := configs.RunDatabase(dbName)
	if err != nil {
		localLog.Error("Database connection failed:", err)
		return err
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingUser models.User
	err = tx.Raw("SELECT * FROM myuser WHERE Email = ? AND FOR UPDATE", user.Email).
		Scan(&existingUser).Error
	if err != nil {
		tx.Rollback()
		return errors.New("Email already registered")
	}

	err = tx.Table("myuser").Create(&user).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		localLog.Error("Transaction commit failed:", err)
		return err
	}

	return nil
}

func (r *userRepository) VerifyUser(dbName, email, securityCode string) (bool, error) {
	db, err := configs.RunDatabase(dbName)
	if err != nil {
		localLog.Error("Database connection failed:", err)
		return false, err
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var user models.User
	err = tx.Raw("SELECT * FROM myuser WHERE Email = ? AND SecurityCode = ? FOR UPDATE", email, securityCode).
		Scan(&user).Error
	if err != nil {
		tx.Rollback()
		localLog.Error("User not found or already verified:", err)
		return false, errors.New("Invalid verification code or user already verified")
	}

	if user.LogIn.Int32 == 1 {
		tx.Rollback()
		return false, errors.New("User already verified")
	}

	result := tx.Model(&user).Update("LogIn", 1)
	if result.RowsAffected == 0 || result.Error != nil {
		tx.Rollback()
		localLog.Error("Failed to update LogIn:", result.Error)
		return false, errors.New("Failed to verify user")
	}

	if err := tx.Commit().Error; err != nil {
		localLog.Error("Transaction commit failed:", err)
		return false, err
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
		Raw("SELECT UserCode, LogIn FROM myuser WHERE UserName = ? AND Password LIKE ?", username, encodedPassword[:len(encodedPassword)-2]+"%").
		Scan(&user).Error

	if err != nil {
		return nil, err
	}

	if !user.LogIn.Valid {
		return nil, errors.New(fmt.Sprintf("Account %s is deactivated", user.UserCode))
	}

	if user.LogIn.Int32 == 0 {
		return nil, errors.New(fmt.Sprintf("Account %s not verified, please verify your account first", user.UserCode))
	}

	return &user, nil
}

func (r *userRepository) SoftDeleteUser(dbName, username, encodedPassword string) error {
	db, err := configs.RunDatabase(dbName)
	if err != nil {
		localLog.Error("Database connection failed:", err)
		return err
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user, err := r.FindByUsernameAndPassword(dbName, username, encodedPassword)
	if err != nil {
		tx.Rollback()
		return errors.New(err.Error())
	}

	err = tx.Raw("SELECT UserCode, LogIn FROM myuser WHERE UserName = ? AND Password LIKE ? FOR UPDATE",
		username, encodedPassword[:len(encodedPassword)-2]+"%").
		Scan(&user).Error
	if err != nil {
		tx.Rollback()
		return errors.New(err.Error())
	}

	result := tx.Table("myuser").Where("UserCode = ?", user.UserCode).Update("LogIn", nil)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("User not found or already soft deleted")
	}

	if err := tx.Commit().Error; err != nil {
		localLog.Error("Transaction commit failed:", err)
		return err
	}

	return nil
}

func (r *userRepository) ResendVerifyCode(dbName, email, securityCode string) error {
	db, err := configs.RunDatabase(dbName)
	if err != nil {
		localLog.Error("Database connection failed:", err)
		return err
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // Rollback jika panic
		}
	}()

	var user models.User
	err = tx.Raw("SELECT * FROM myuser WHERE Email = ? FOR UPDATE", email).
		Scan(&user).Error
	if err != nil {
		tx.Rollback()
		return errors.New("User not found")
	}

	securityCode, err = utils.GenerateSecurityCode()
	if err != nil {
		tx.Rollback()
		log.Error("Failed to generate security code:", err)
		return err
	}

	if err := tx.Model(&user).Update("SecurityCode", securityCode).Error; err != nil {
		tx.Rollback()
		return errors.New("Failed to update verification code")
	}

	if err := tx.Commit().Error; err != nil {
		localLog.Error("Transaction commit failed:", err)
		return err
	}

	// if err := utils.SendVerificationEmail(email, securityCode); err != nil {
	// 	return false, errors.New("Failed to send verification email")
	// }

	return nil
}
