package app

import (
	"testing"

	"github.com/Basu008/EasySplit.git/model"
	"github.com/Basu008/EasySplit.git/schema"
	"github.com/stretchr/testify/assert"
)

func TestSignUpUser_OK(t *testing.T) {
	testApp := NewTestApp(getTestConfig())
	defer cleanUpDB(testApp.Postgres.DB, &model.User{})
	var serviceErr error
	testApp.User, serviceErr = InitUser(&UserImplOpts{
		App: testApp,
		DB:  testApp.Postgres.DB,
	})
	opts := schema.SignupOpts{
		FullName:    "John Doe",
		Username:    "johndoe123",
		Password:    "StrongP@ssw0rd!",
		PhoneNumber: "9876543210",
		CountryCode: "+91",
		Email:       "johndoe@example.com",
	}
	user, customErr := testApp.User.SignupUser(&opts)
	assert.Nil(t, customErr)
	assert.Nil(t, serviceErr)
	if assert.NotNil(t, user) {
		assert.Equal(t, "johndoe123", user.Username)
	}
}
