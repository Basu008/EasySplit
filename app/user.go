package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Basu008/EasySplit.git/model"
	"github.com/Basu008/EasySplit.git/schema"
	"github.com/Basu008/EasySplit.git/server/auth"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User interface {
	SignupUser(opts *schema.SignupOpts) (*model.User, *model.Error)
	LoginUser(opts *schema.LoginOpts) (*auth.UserClaim, *model.Error)
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

func (ui *UserImpl) SignupUser(opts *schema.SignupOpts) (*model.User, *model.Error) {
	user := model.User{
		FullName:    opts.FullName,
		Username:    opts.Username,
		CountryCode: opts.CountryCode,
		PhoneNumber: opts.PhoneNumber,
		Email:       opts.Email,
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(opts.Password), 4)
	if err != nil {
		return nil, model.NewError(err, "", http.StatusBadRequest)
	}
	user.Password = string(encryptedPassword)
	if err := ui.DB.Create(&user).Error; err != nil {
		return nil, model.NewError(err, "", http.StatusInternalServerError)
	}
	return &user, nil
}

func (ui *UserImpl) LoginUser(opts *schema.LoginOpts) (*auth.UserClaim, *model.Error) {
	user, customErr := ui.GetUserByUsername(opts.Username)
	if customErr != nil {
		return nil, customErr
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(opts.Password))
	if err != nil {
		return nil, model.NewError(err, "incorrect password", http.StatusBadRequest)
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
