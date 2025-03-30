package repositories

import (
	"echo-api/configs"
	models "echo-api/models"
	"echo-api/utils"
	"errors"
	"fmt"

	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	// securityCode, err := utils.GenerateSecurityCode()
	// if err != nil {
	// 	log.Error("Failed to generate security code:", err)
	// 	return false, err
	// }

	// user.SecurityCode = securityCode
	// user.LogIn = 0
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // Rollback jika panic
		}
	}()

	var existingUser models.User
	err = tx.Raw("SELECT * FROM myuser WHERE Email = ? AND FOR UPDATE", user.Email).
		Scan(&existingUser).Error
	if err == nil {
		tx.Rollback()
		return errors.New("Email already registered")
	}

	err = tx.Table("myuser").Create(&user).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// if user.Email == "" {
	// 	return false, errors.New("Email cannot be null")
	// }

	// err = utils.SendVerificationEmail(user.Email, securityCode)
	// if err != nil {
	// 	message := fmt.Sprintf("Failed to send verification email to %s", user.UserCode)
	// 	r.log.Error(message, err)
	// 	return false, errors.New("Failed to send verification email")
	// }
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

	if user.LogIn == 1 {
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
		Clauses(clause.Locking{Strength: "UPDATE"}).
		// Select("UserCode, Login").
		Raw("SELECT UserCode, LogIn FROM myuser WHERE UserName = ? AND Password LIKE ?", username, encodedPassword[:len(encodedPassword)-2]+"%").
		// Where("UserName = ? AND Password LIKE ?", username, encodedPassword[:len(encodedPassword)-2]+"%").
		Scan(&user).Error

	if err != nil {
		return nil, err
	}

	if user.LogIn == 0 {
		return nil, errors.New(fmt.Sprintf("Account %s not verified, please verify your account first", user.UserCode))
	}

	return &user, nil
}

func FindByUsernameAndPassword(dbName, username, encodedPassword string) (*models.User, error) {
	db, err := configs.RunDatabase(dbName)
	if err != nil {
		localLog.Error("Database connection failed:", err)
		return nil, err
	}

	db = db.Session(&gorm.Session{SkipDefaultTransaction: true})

	var user models.User
	err = db.Table("myuser").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		// Select("UserCode, Login").
		Raw("SELECT UserCode FROM myuser WHERE UserName = ? AND Password LIKE ? AND LogIn = 1", username, encodedPassword[:len(encodedPassword)-2]+"%").
		// Where("UserName = ? AND Password LIKE ?", username, encodedPassword[:len(encodedPassword)-2]+"%").
		Scan(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
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
