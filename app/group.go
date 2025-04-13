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

type Group interface {
	//Create
	CreateGroup(opts *schema.CreateGroupOpts) *model.Error
	//Get
	GetGroupByID(id uint) (*schema.GroupResponse, *model.Error)
	GetGroups(ownerID uint, page int) ([]schema.GroupResponse, *model.Error)
	//Edit
	EditGroup(opts *schema.EditGroupInfoOpts) *model.Error
	AddGroupMembers(opts *schema.AddGroupMembersOpts) *model.Error
	RemoveGroupMember(opts *schema.RemoveGroupMemberOpts) *model.Error
	MigrateGroup() error
}

type GroupImplOpts struct {
	App *App
	DB  *gorm.DB
}

type GroupImpl struct {
	App *App
	DB  *gorm.DB
}

func InitGroup(opts *GroupImplOpts) (Group, error) {
	ei := GroupImpl{
		App: opts.App,
		DB:  opts.DB,
	}
	if err := ei.MigrateGroup(); err != nil {
		log.Print(err)
		return nil, errors.New("unable to migrate Group")
	}
	if err := ei.MigrateGroupMember(); err != nil {
		log.Print(err)
		return nil, errors.New("unable to migrate Group Member")
	}
	return &ei, nil
}

func (gi *GroupImpl) CreateGroup(opts *schema.CreateGroupOpts) *model.Error {
	tx := gi.DB.Begin()
	group := model.Group{
		Name:      opts.Name,
		CreatedBy: opts.OwnerID,
		Type:      opts.Type,
	}
	if err := tx.Create(&group).Error; err != nil {
		tx.Rollback()
		return &model.Error{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}
	members := []model.GroupMember{
		{
			GroupID: group.ID,
			UserID:  opts.OwnerID,
		},
	}
	for _, userID := range opts.UserIDs {
		members = append(members, model.GroupMember{
			GroupID: group.ID,
			UserID:  userID,
		})
	}
	if err := tx.CreateInBatches(members, len(members)).Error; err != nil {
		tx.Rollback()
		return &model.Error{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}
	err := tx.Commit().Error
	if err != nil {
		return &model.Error{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return nil
}

func (gi *GroupImpl) GetGroupByID(id uint) (*schema.GroupResponse, *model.Error) {
	var group model.Group
	var members []model.GroupMember

	err := gi.DB.First(&group, id).Error
	if err != nil {
		return nil, &model.Error{
			Err:  fmt.Errorf("group not found: %w", err),
			Code: http.StatusBadRequest,
		}
	}

	err = gi.DB.Preload("User").Where("group_id = ?", id).Find(&members).Error
	if err != nil {
		return nil, &model.Error{
			Err:  fmt.Errorf("failed to fetch members: %w", err),
			Code: http.StatusBadRequest,
		}
	}
	return groupResponse(&group, members), nil
}

func (gi *GroupImpl) GetGroups(ownerID uint, page int) ([]schema.GroupResponse, *model.Error) {
	var groups []model.Group
	limit := gi.App.Config.Limit
	offset := limit * page
	if err := gi.DB.
		Where("created_by = ?", ownerID).
		Offset(offset).
		Limit(limit).
		Find(&groups).
		Error; err != nil {
		return nil, &model.Error{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	groupResp := []schema.GroupResponse{}
	for _, group := range groups {
		groupResp = append(groupResp, *groupResponse(&group, nil))
	}
	return groupResp, nil
}

func (gi *GroupImpl) EditGroup(opts *schema.EditGroupInfoOpts) *model.Error {
	group := model.Group{
		ID: opts.ID,
	}
	update := make(map[string]any)
	if opts.Name != "" {
		update[model.GroupName] = opts.Name
	}
	if opts.Type != "" {
		update[model.GroupType] = opts.Type
	}
	err := gi.DB.Model(&group).Updates(update).Error
	if err != nil {
		return &model.Error{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return nil
}

func (gi *GroupImpl) AddGroupMembers(opts *schema.AddGroupMembersOpts) *model.Error {
	members := []model.GroupMember{}
	for _, userID := range opts.UserIDs {
		members = append(members, model.GroupMember{
			GroupID: opts.ID,
			UserID:  userID,
		})
	}
	if err := gi.DB.CreateInBatches(members, len(members)).Error; err != nil {
		return &model.Error{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}
	return nil
}

func (gi *GroupImpl) RemoveGroupMember(opts *schema.RemoveGroupMemberOpts) *model.Error {
	if err := gi.DB.Where("group_id = ? AND user_id = ?", opts.ID, opts.UserID).Delete(&model.GroupMember{}).Error; err != nil {
		return &model.Error{
			Err:  fmt.Errorf("failed to delete member: %w", err),
			Code: http.StatusInternalServerError,
		}
	}
	return nil
}

func groupResponse(group *model.Group, members []model.GroupMember) *schema.GroupResponse {
	groupResp := schema.GroupResponse{
		ID:   group.ID,
		Name: group.Name,
		Type: group.Type,
	}
	membersData := []schema.MemberResponse{}
	// members := []schema.UserResponse{}
	for _, resp := range members {
		member := resp.User
		membersData = append(membersData, schema.MemberResponse{
			ID:       member.ID,
			Username: member.Username,
			FullName: member.FullName,
		})
	}
	groupResp.Members = membersData
	return &groupResp
}

func (gi *GroupImpl) MigrateGroup() error {
	err := gi.DB.AutoMigrate(&model.Group{})
	return err
}

func (gi *GroupImpl) MigrateGroupMember() error {
	err := gi.DB.AutoMigrate(&model.GroupMember{})
	return err
}
