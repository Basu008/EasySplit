package app

import (
	"testing"

	"github.com/Basu008/EasySplit.git/model"
	"github.com/Basu008/EasySplit.git/schema"
	"github.com/stretchr/testify/assert"
)

func setupTestFriend(t *testing.T) (*App, Friend) {
	testApp := NewTestApp(getTestConfig())
	t.Cleanup(func() { cleanUpDB(testApp.Postgres.DB) })
	friendService, serviceErr := InitFriend(&FriendImplOpts{
		App: testApp,
		DB:  testApp.Postgres.DB,
	})
	assert.Nil(t, serviceErr)
	testApp.Friend = friendService

	userService, serviceErr := InitUser(&UserImplOpts{
		App: testApp,
		DB:  testApp.Postgres.DB,
	})
	assert.Nil(t, serviceErr)
	testApp.User = userService
	return testApp, friendService
}

func createFriendRequest(t *testing.T, friendService Friend, senderID, receiverID uint) {
	opts := &schema.FriendRequestOpts{
		SenderUserID:   senderID,
		ReceiverUserID: receiverID,
	}
	err := friendService.SendFriendRequest(opts)
	assert.Nil(t, err)
}

func createTestUsers(t *testing.T, testApp *App) []uint {
	ids := []uint{}
	opts := schema.SignupOpts{
		FullName:    "Jane Doe",
		Username:    "janedoe123",
		Password:    "StrongPassword123!",
		PhoneNumber: "9876543210",
		CountryCode: "+91",
		Email:       "janedoe@example.com",
	}
	user1, err := testApp.User.SignupUser(&opts)
	assert.Nil(t, err)
	assert.NotNil(t, user1)
	ids = append(ids, user1.ID)
	opts2 := schema.SignupOpts{
		FullName:    "Doe Jane",
		Username:    "doejane123",
		Password:    "StrongPassword123!",
		PhoneNumber: "9999999999",
		CountryCode: "+91",
		Email:       "doejane@example.com",
	}
	user2, err := testApp.User.SignupUser(&opts2)
	assert.Nil(t, err)
	assert.NotNil(t, user2)
	ids = append(ids, user2.ID)
	return ids
}

func TestSendFriendRequest_OK(t *testing.T) {
	testApp, friendService := setupTestFriend(t)
	ids := createTestUsers(t, testApp)
	senderID, userID := ids[0], ids[1]
	opts := &schema.FriendRequestOpts{
		SenderUserID:   senderID,
		ReceiverUserID: userID,
	}
	err := friendService.SendFriendRequest(opts)

	assert.Nil(t, err)
}

func TestSendFriendRequest_AlreadyExists(t *testing.T) {
	testApp, friendService := setupTestFriend(t)
	ids := createTestUsers(t, testApp)
	senderID, receiverID := ids[0], ids[1]
	createFriendRequest(t, friendService, senderID, receiverID)
	err := friendService.SendFriendRequest(&schema.FriendRequestOpts{
		SenderUserID:   receiverID,
		ReceiverUserID: senderID,
	})

	assert.NotNil(t, err)
	assert.Equal(t, model.RequestAlreadyExists, err.Message)
	assert.Equal(t, 400, err.Code)
}

func TestUpdateFriendRequest_Accept_OK(t *testing.T) {
	testApp, friendService := setupTestFriend(t)
	ids := createTestUsers(t, testApp)
	senderID, receiverID := ids[0], ids[1]
	createFriendRequest(t, friendService, senderID, receiverID)

	updateOpts := &schema.UpdateFriendRequestOpts{
		SenderUserID:   senderID,
		ReceiverUserID: receiverID,
		RequestStatus:  model.Accepted,
	}
	err := friendService.UpdateFriendRequest(updateOpts)

	assert.Nil(t, err)
}

func TestUpdateFriendRequest_Reject_OK(t *testing.T) {
	testApp, friendService := setupTestFriend(t)
	ids := createTestUsers(t, testApp)
	senderID, receiverID := ids[0], ids[1]
	createFriendRequest(t, friendService, senderID, receiverID)

	updateOpts := &schema.UpdateFriendRequestOpts{
		SenderUserID:   senderID,
		ReceiverUserID: receiverID,
		RequestStatus:  model.Rejected,
	}
	err := friendService.UpdateFriendRequest(updateOpts)

	assert.Nil(t, err)
}

func TestGetAllFriends_OK(t *testing.T) {
	testApp, friendService := setupTestFriend(t)
	ids := createTestUsers(t, testApp)
	senderID, receiverID := ids[0], ids[1]
	createFriendRequest(t, friendService, senderID, receiverID)
	updateOpts := &schema.UpdateFriendRequestOpts{
		SenderUserID:   senderID,
		ReceiverUserID: receiverID,
		RequestStatus:  model.Accepted,
	}
	_ = friendService.UpdateFriendRequest(updateOpts)

	friends := friendService.GetAllFriends(senderID, 0)

	assert.NotNil(t, friends)
	assert.GreaterOrEqual(t, len(friends), 0)
}

func TestGetFriendStatus_OK(t *testing.T) {
	testApp, friendService := setupTestFriend(t)
	ids := createTestUsers(t, testApp)
	senderID, receiverID := ids[0], ids[1]
	createFriendRequest(t, friendService, senderID, receiverID)

	friend, err := friendService.GetFriendStatus(senderID, receiverID)

	assert.Nil(t, err)
	assert.NotNil(t, friend)
	assert.Equal(t, senderID, friend.SenderUserID)
	assert.Equal(t, receiverID, friend.ReceiverUserID)
}

func TestGetFriendStatus_NotFound(t *testing.T) {
	testApp, friendService := setupTestFriend(t)
	ids := createTestUsers(t, testApp)
	senderID, receiverID := ids[0], ids[1]
	friend, err := friendService.GetFriendStatus(senderID, receiverID)

	assert.NotNil(t, err)
	assert.Nil(t, friend)
	assert.Equal(t, model.RequestDoesntExist, err.Message)
	assert.Equal(t, 400, err.Code)
}
