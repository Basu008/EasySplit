package app

import (
	"net/http"
	"testing"

	"github.com/Basu008/EasySplit.git/model"
	"github.com/Basu008/EasySplit.git/schema"
	"github.com/stretchr/testify/assert"
)

func setupTestUser(t *testing.T) (*App, User) {
	testApp := NewTestApp(getTestConfig())
	userService, serviceErr := InitUser(&UserImplOpts{
		App: testApp,
		DB:  testApp.Postgres.DB,
	})
	assert.Nil(t, serviceErr)

	testApp.User = userService
	return testApp, userService
}

func createTestUser(t *testing.T, testApp *App) *model.User {
	opts := schema.SignupOpts{
		FullName:    "Jane Doe",
		Username:    "janedoe123",
		Password:    "StrongPassword123!",
		PhoneNumber: "9876543210",
		CountryCode: "+91",
		Email:       "janedoe@example.com",
	}
	user, err := testApp.User.SignupUser(&opts)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	return user
}

func TestSignUpUser_OK(t *testing.T) {
	testApp, userService := setupTestUser(t)
	defer cleanUpDB(testApp.Postgres.DB, &model.User{})
	opts := schema.SignupOpts{
		FullName:    "John Doe",
		Username:    "johndoe123",
		Password:    "StrongP@ssw0rd!",
		PhoneNumber: "9876543210",
		CountryCode: "+91",
		Email:       "johndoe@example.com",
	}
	user, customErr := userService.SignupUser(&opts)
	assert.Nil(t, customErr)
	if assert.NotNil(t, user) {
		assert.Equal(t, "johndoe123", user.Username)
	}
}

func TestLoginUser_OK(t *testing.T) {
	testApp, userService := setupTestUser(t)
	defer cleanUpDB(testApp.Postgres.DB, &model.User{})
	user := createTestUser(t, testApp)

	loginOpts := schema.LoginOpts{
		Username: user.Username,
		Password: "StrongPassword123!",
	}
	claim, err := userService.LoginUser(&loginOpts)

	assert.Nil(t, err)
	assert.NotNil(t, claim)
	assert.Equal(t, user.PhoneNumber, claim.PhoneNumber)
	assert.Equal(t, user.ID, claim.ID)
}

func TestGetUserByID_OK(t *testing.T) {
	testApp, userService := setupTestUser(t)
	defer cleanUpDB(testApp.Postgres.DB, &model.User{})
	user := createTestUser(t, testApp)

	foundUser, err := userService.GetUserByID(user.ID)

	assert.Nil(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Username, foundUser.Username)
}

func TestGetUserByPhoneNo_OK(t *testing.T) {
	testApp, userService := setupTestUser(t)
	defer cleanUpDB(testApp.Postgres.DB, &model.User{})
	user := createTestUser(t, testApp)

	foundUser, err := userService.GetUserByPhoneNo(user.PhoneNumber)

	assert.Nil(t, err)
	assert.Equal(t, user.PhoneNumber, foundUser.PhoneNumber)
	assert.Equal(t, user.Username, foundUser.Username)
}

func TestGetUserByUsername_OK(t *testing.T) {
	testApp, userService := setupTestUser(t)
	defer cleanUpDB(testApp.Postgres.DB, &model.User{})
	user := createTestUser(t, testApp)

	foundUser, err := userService.GetUserByUsername(user.Username)

	assert.Nil(t, err)
	assert.Equal(t, user.Username, foundUser.Username)
	assert.Equal(t, user.Email, foundUser.Email)
}

func TestUpdateUser_OK(t *testing.T) {
	testApp, userService := setupTestUser(t)
	defer cleanUpDB(testApp.Postgres.DB, &model.User{})
	user := createTestUser(t, testApp)

	updateOpts := schema.UpdateUserOpts{
		ID:       user.ID,
		Username: "updatedusername123",
		Email:    "updated@example.com",
	}
	err := userService.UpdateUser(&updateOpts)

	assert.Nil(t, err)

	// Verify updates
	updatedUser, getErr := userService.GetUserByID(user.ID)
	assert.Nil(t, getErr)
	assert.Equal(t, "updatedusername123", updatedUser.Username)
	assert.Equal(t, "updated@example.com", updatedUser.Email)
}

func TestLoginUser_WrongPassword(t *testing.T) {
	testApp, userService := setupTestUser(t)
	user := createTestUser(t, testApp)

	loginOpts := schema.LoginOpts{
		Username: user.Username,
		Password: "WrongPassword!",
	}
	claim, err := userService.LoginUser(&loginOpts)

	assert.Nil(t, claim)
	assert.NotNil(t, err)
	assert.Equal(t, "incorrect password", err.Message)
	assert.Equal(t, http.StatusBadRequest, err.Code)
}

func TestGetUserByID_NotFound(t *testing.T) {
	_, userService := setupTestUser(t)

	foundUser, err := userService.GetUserByID(99999) // random non-existent ID

	assert.NotNil(t, err)
	assert.Equal(t, uint(0), foundUser.ID)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Equal(t, model.InvalidPhoneNo, err.Message)
}

func TestGetUserByPhoneNo_NotFound(t *testing.T) {
	_, userService := setupTestUser(t)

	foundUser, err := userService.GetUserByPhoneNo("0000000000") // invalid phone

	assert.NotNil(t, err)
	assert.Equal(t, "", foundUser.Username)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Equal(t, model.InvalidPhoneNo, err.Message)
}

func TestGetUserByUsername_NotFound(t *testing.T) {
	_, userService := setupTestUser(t)

	foundUser, err := userService.GetUserByUsername("nonexistentusername") // invalid username

	assert.NotNil(t, err)
	assert.Equal(t, "", foundUser.Username)
	assert.Equal(t, http.StatusBadRequest, err.Code)
	assert.Equal(t, model.InvalidUsername, err.Message)
}

func TestUpdateUser_UserNotFound(t *testing.T) {
	_, userService := setupTestUser(t)

	updateOpts := schema.UpdateUserOpts{
		ID:       99999, // non-existent ID
		Username: "shouldnotexist",
		Email:    "shouldnotexist@example.com",
	}
	err := userService.UpdateUser(&updateOpts)

	assert.Nil(t, err)
}
