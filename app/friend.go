package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Basu008/EasySplit.git/model"
	"github.com/Basu008/EasySplit.git/schema"
	"gorm.io/gorm"
)

type Friend interface {
	SendFriendRequest(opts *schema.FriendRequestOpts) *model.Error
	UpdateFriendRequest(opts *schema.UpdateFriendRequestOpts) *model.Error
	GetAllFriends(userID uint, page int) []schema.GetAllFriendsResponse
	GetFriendStatus(userID, friendUserID uint) (*model.Friend, *model.Error)
	MigrateFriend() error
}

type FriendImpl struct {
	App *App
	DB  *gorm.DB
}

type FriendImplOpts struct {
	App *App
	DB  *gorm.DB
}

func InitFriend(opts *FriendImplOpts) (Friend, error) {
	fi := FriendImpl{
		App: opts.App,
		DB:  opts.DB,
	}
	if err := fi.MigrateFriend(); err != nil {
		log.Fatal(err)
		return nil, errors.New("unable to migrate Friend")
	}
	return &fi, nil
}

func (fi *FriendImpl) SendFriendRequest(opts *schema.FriendRequestOpts) *model.Error {
	friendModel := &model.Friend{}

	whereQuery := fmt.Sprintf("%s = ? AND %s = ?", model.SenderUserID, model.ReceiverUserID)

	var count int64
	fi.DB.Model(friendModel).Where(whereQuery, opts.ReceiverUserID, opts.SenderUserID).
		Count(&count)

	if count != 0 {
		return model.NewError(nil, model.RequestAlreadyExists, http.StatusBadRequest)
	}
	friend := model.Friend{
		SenderUserID:   opts.SenderUserID,
		ReceiverUserID: opts.ReceiverUserID,
		RequestStatus:  model.Requested,
	}
	err := fi.DB.Create(&friend).Error
	if err != nil {
		return model.NewError(err, "", http.StatusBadRequest)
	}
	return nil
}

func (fi *FriendImpl) UpdateFriendRequest(opts *schema.UpdateFriendRequestOpts) *model.Error {
	friendModel := &model.Friend{}
	whereQuery := fmt.Sprintf("%s = ? AND %s = ? AND %s = ?", model.SenderUserID, model.ReceiverUserID, model.RequestStatus)
	if opts.RequestStatus == model.Rejected {
		err := fi.DB.Delete(friendModel, whereQuery, opts.SenderUserID, opts.ReceiverUserID, model.Requested).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return model.NewError(err, model.RequestDoesntExist, http.StatusBadRequest)
			}
			return model.NewError(err, model.RequestProcessUnable, http.StatusBadRequest)
		}
	} else {
		err := fi.DB.Model(friendModel).
			Where(whereQuery, opts.SenderUserID, opts.ReceiverUserID, model.Requested).
			Update(model.RequestStatus, opts.RequestStatus).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return model.NewError(err, model.RequestDoesntExist, http.StatusBadRequest)
			}
			return model.NewError(err, model.RequestProcessUnable, http.StatusInternalServerError)
		}
	}
	return nil
}

func (fi *FriendImpl) GetAllFriends(userID uint, page int) []schema.GetAllFriendsResponse {
	friends := []schema.GetAllFriendsResponse{}
	offset := page * fi.App.Config.Limit
	fi.DB.Table("users").
		Select("users.id, users.username, users.phone_number").
		Joins("LEFT JOIN friends ON users.id = friends.receiver_user_id").
		Where("friends.sender_user_id = ? AND friends.request_status = ?", userID, model.Accepted).
		Offset(offset).
		Limit(fi.App.Config.Limit).
		Scan(&friends)
	return friends
}

func (fi *FriendImpl) GetFriendStatus(userID, friendUserID uint) (*model.Friend, *model.Error) {
	friend := model.Friend{}
	friend.SenderUserID = userID
	friend.ReceiverUserID = friendUserID
	err := fi.DB.First(&friend).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewError(err, model.RequestDoesntExist, http.StatusBadRequest)
		}
		return nil, model.NewError(err, "", http.StatusInternalServerError)
	}
	return &friend, nil
}

func (fi *FriendImpl) MigrateFriend() error {
	err := fi.DB.AutoMigrate(&model.Friend{})
	return err
}
