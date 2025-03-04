package app

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Basu008/EasySplit.git/model"
	"github.com/Basu008/EasySplit.git/schema"
	"github.com/Basu008/EasySplit.git/server/auth"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User interface {
	//Login
	LoginUser(opts *schema.PhoneNoLogin) *model.Error
	ConfirmOTP(opts *schema.ConfirmOTPOpts) (*auth.UserClaim, *model.Error)
	//Get user
	GetUserByID(id uint) (model.User, *model.Error)
	GetUserByPhoneNo(phoneNumber string) (model.User, *model.Error)
	GetUserByUsername(username string) (model.User, *model.Error)
	//Edit User
	UpdateUser(opts *schema.UpdateUserOpts) *model.Error

	MigrateUser() error
}

type UserImpl struct {
	App *App
	DB  *gorm.DB
}

type UserImplOpts struct {
	App *App
	DB  *gorm.DB
}

func InitUser(opts *UserImplOpts) (User, error) {
	ui := UserImpl{
		App: opts.App,
		DB:  opts.DB,
	}
	if err := ui.MigrateUser(); err != nil {
		log.Fatal(err)
		return nil, errors.New("unable to migrate User")
	}
	return &ui, nil
}

func (ui *UserImpl) LoginUser(opts *schema.PhoneNoLogin) *model.Error {
	otp := ui.generateOTP()
	user := model.User{
		PhoneNumber: opts.PhoneNumber,
		CountryCode: opts.CountryCode,
		OTP:         otp,
	}
	update := map[string]any{
		model.OTP: otp,
	}
	err := ui.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: model.PhoneNumber,
			},
		},
		DoUpdates: clause.Assignments(update),
	}).Create(&user).Error
	if err != nil {
		if err == gorm.ErrCheckConstraintViolated {
			return model.NewError(err, "invalid phone_number", http.StatusBadRequest)
		}
		return model.NewError(err, "", http.StatusInternalServerError)
	}
	return nil
}

func (ui *UserImpl) ConfirmOTP(opts *schema.ConfirmOTPOpts) (*auth.UserClaim, *model.Error) {
	user, customErr := ui.GetUserByPhoneNo(opts.PhoneNumber)
	if customErr != nil {
		return nil, customErr
	}
	if user.OTP != opts.OTP {
		newErr := model.NewError(nil, "invalid otp", http.StatusBadRequest)
		return nil, newErr
	}
	updates := make(map[string]any)
	updates[model.OTP] = "-"
	if !user.PhoneVerified {
		updates[model.PhoneVerified] = true
	}
	err := ui.DB.Model(&user).Updates(updates).Error
	if err != nil {
		newErr := model.NewError(err, "", http.StatusInternalServerError)
		return nil, newErr
	}
	claim := auth.UserClaim{
		ID:          user.ID,
		Plan:        user.Plan,
		PhoneNumber: user.PhoneNumber,
	}
	return &claim, nil
}

func (ui *UserImpl) GetUserByID(id uint) (model.User, *model.Error) {
	var user model.User
	err := ui.DB.First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errBody := model.NewError(err, model.InvalidPhoneNo, http.StatusBadRequest)
			return user, errBody
		}
		return user, model.NewError(err, "", http.StatusInternalServerError)
	}
	return user, nil
}

func (ui *UserImpl) GetUserByPhoneNo(phoneNumber string) (model.User, *model.Error) {
	var user model.User
	whereQuery := fmt.Sprintf("%s = ?", model.PhoneNumber)
	err := ui.DB.Where(whereQuery, phoneNumber).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		errBody := model.NewError(err, model.InvalidPhoneNo, http.StatusBadRequest)
		return user, errBody
	}
	return user, nil
}

func (ui *UserImpl) GetUserByUsername(username string) (model.User, *model.Error) {
	var user model.User
	whereQuery := fmt.Sprintf("%s = ?", model.Username)
	err := ui.DB.Where(whereQuery, username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errBody := model.NewError(err, model.InvalidUsername, http.StatusBadRequest)
			return user, errBody
		}
		return user, model.NewError(err, "", http.StatusInternalServerError)
	}
	return user, nil
}

func (ui *UserImpl) UpdateUser(opts *schema.UpdateUserOpts) *model.Error {
	user := model.User{
		ID: opts.ID,
	}
	updates := make(map[string]any)
	if opts.Username != "" {
		updates[model.Username] = opts.Username
	}
	if opts.Email != "" {
		updates[model.Email] = opts.Email
	}
	err := ui.DB.Model(&user).Updates(
		updates,
	).Error
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			return model.NewError(nil, "duplicate email/username not allowed", http.StatusBadRequest)
		}
		if err == gorm.ErrRecordNotFound {
			return model.NewError(nil, "user doesn't exists", http.StatusBadRequest)
		}
		return model.NewError(err, "", http.StatusInternalServerError)
	}
	return nil
}

func (ui *UserImpl) MigrateUser() error {
	err := ui.DB.AutoMigrate(&model.User{})
	return err
}

func (ui *UserImpl) generateOTP() string {
	length := ui.App.Config.OTPLength
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return ""
	}
	otpChars := ui.App.Config.OTPChars
	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}
	return string(buffer)
}
