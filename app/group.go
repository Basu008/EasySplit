package app

import (
	"errors"
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
		Name:    opts.Name,
		OwnerID: opts.OwnerID,
		Type:    opts.Type,
	}
	if err := tx.Create(&group).Error; err != nil {
		tx.Rollback()
		return &model.Error{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}
	members := []model.GroupMember{}
	for _, memberID := range opts.MemberIDs {
		members = append(members, model.GroupMember{
			GroupID:  group.ID,
			MemberID: memberID,
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
	group.ID = id
	if err := gi.DB.Preload("Members.Member").First(&group).Error; err != nil {
		return nil, &model.Error{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return groupResponse(&group), nil
}

func (gi *GroupImpl) GetGroups(ownerID uint, page int) ([]schema.GroupResponse, *model.Error) {
	var groups []model.Group
	limit := gi.App.Config.Limit
	offset := limit * page
	if err := gi.DB.Preload("Members.Member").
		Where("owner_id = ?", ownerID).
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
		groupResp = append(groupResp, *groupResponse(&group))
	}
	return groupResp, nil
}

func (gi *GroupImpl) EditGroup(opts *schema.EditGroupInfoOpts) *model.Error {
	group := model.Group{
		ID: opts.ID,
	}
	update := make(map[string]interface{})
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

func groupResponse(group *model.Group) *schema.GroupResponse {
	groupResp := schema.GroupResponse{
		ID:            group.ID,
		Name:          group.Name,
		Type:          group.Type,
		TotalExpense:  group.TotalExpense,
		SettledAmount: group.SettledAmount,
	}
	members := []schema.MemberResponse{}
	for _, resp := range group.Members {
		member := resp.Member
		members = append(members, schema.MemberResponse{
			ID:       member.ID,
			Username: *member.Username,
		})
	}
	groupResp.Members = members
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
