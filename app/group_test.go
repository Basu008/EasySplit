package app

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Basu008/EasySplit.git/model"
	"github.com/Basu008/EasySplit.git/schema"
	"github.com/stretchr/testify/assert"
)

func setupTestGroup(t *testing.T) (*App, Group) {
	testApp := NewTestApp(getTestConfig())
	t.Cleanup(func() { cleanUpDB(testApp.Postgres.DB) })
	groupService, serviceErr := InitGroup(&GroupImplOpts{
		App: testApp,
		DB:  testApp.Postgres.DB,
	})
	assert.Nil(t, serviceErr)
	testApp.Group = groupService

	userService, serviceErr := InitUser(&UserImplOpts{
		App: testApp,
		DB:  testApp.Postgres.DB,
	})
	assert.Nil(t, serviceErr)
	testApp.User = userService

	return testApp, groupService
}

func createGroupUser(t *testing.T, testApp *App) *model.User {
	signupOpts := schema.SignupOpts{
		FullName:    "Group Owner",
		Username:    RandString(8),
		Password:    "StrongPassword@123",
		PhoneNumber: RandPhoneNumber(),
		CountryCode: "+91",
		Email:       RandEmail(),
	}
	user, err := testApp.User.SignupUser(&signupOpts)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	return user
}

func createGroup(t *testing.T, testApp *App) (*schema.GroupResponse, uint) {
	owner := createGroupUser(t, testApp)
	createOpts := &schema.CreateGroupOpts{
		Name:    "Test Group",
		OwnerID: owner.ID,
		Type:    "friends",
		UserIDs: []uint{},
	}
	err := testApp.Group.CreateGroup(createOpts)
	assert.Nil(t, err)

	groups, fetchErr := testApp.Group.GetGroups(owner.ID, 0)
	assert.Nil(t, fetchErr)
	assert.NotEmpty(t, groups)
	return &groups[0], owner.ID
}

func TestCreateGroup_OK(t *testing.T) {
	testApp, groupService := setupTestGroup(t)
	owner := createGroupUser(t, testApp)
	var memberIDs []uint
	for range 2 {
		member := createGroupUser(t, testApp)
		memberIDs = append(memberIDs, member.ID)
	}
	createOpts := &schema.CreateGroupOpts{
		Name:    "New Group",
		OwnerID: owner.ID,
		Type:    "family",
		UserIDs: memberIDs,
	}
	err := groupService.CreateGroup(createOpts)

	assert.Nil(t, err)
}

func TestGetGroupByID_OK(t *testing.T) {
	testApp, groupService := setupTestGroup(t)
	group, _ := createGroup(t, testApp)

	foundGroup, err := groupService.GetGroupByID(group.ID)

	assert.Nil(t, err)
	assert.NotNil(t, foundGroup)
	assert.Equal(t, group.ID, foundGroup.ID)
	assert.Equal(t, "Test Group", foundGroup.Name)
}

func TestGetGroupByID_NotFound(t *testing.T) {
	_, groupService := setupTestGroup(t)
	foundGroup, err := groupService.GetGroupByID(9999)

	assert.Nil(t, foundGroup)
	assert.NotNil(t, err)
	assert.Contains(t, err.Err.Error(), "group not found")
	assert.Equal(t, 400, err.Code)
}

func TestGetGroups_OK(t *testing.T) {
	testApp, groupService := setupTestGroup(t)
	_, ownerID := createGroup(t, testApp)

	groups, err := groupService.GetGroups(ownerID, 0)

	assert.Nil(t, err)
	assert.NotNil(t, groups)
	assert.GreaterOrEqual(t, len(groups), 1)
}

func TestEditGroup_OK(t *testing.T) {
	testApp, groupService := setupTestGroup(t)
	group, _ := createGroup(t, testApp)

	editOpts := &schema.EditGroupInfoOpts{
		ID:   group.ID,
		Name: "Updated Group Name",
		Type: "work",
	}
	err := groupService.EditGroup(editOpts)

	assert.Nil(t, err)

	updatedGroup, getErr := groupService.GetGroupByID(group.ID)
	assert.Nil(t, getErr)
	assert.Equal(t, "Updated Group Name", updatedGroup.Name)
	assert.Equal(t, "work", updatedGroup.Type)
}

func TestAddGroupMembers_OK(t *testing.T) {
	testApp, groupService := setupTestGroup(t)
	group, _ := createGroup(t, testApp)
	member1 := createGroupUser(t, testApp)
	member2 := createGroupUser(t, testApp)
	addMembersOpts := &schema.AddGroupMembersOpts{
		ID:      group.ID,
		UserIDs: []uint{member1.ID, member2.ID},
	}
	err := groupService.AddGroupMembers(addMembersOpts)

	assert.Nil(t, err)
}

func TestRemoveGroupMember_OK(t *testing.T) {
	testApp, groupService := setupTestGroup(t)
	group, _ := createGroup(t, testApp)
	member := createGroupUser(t, testApp)
	addMembersOpts := &schema.AddGroupMembersOpts{
		ID:      group.ID,
		UserIDs: []uint{member.ID},
	}
	_ = groupService.AddGroupMembers(addMembersOpts)

	removeOpts := &schema.RemoveGroupMemberOpts{
		ID:     group.ID,
		UserID: member.ID,
	}
	err := groupService.RemoveGroupMember(removeOpts)

	assert.Nil(t, err)
}

func TestRemoveGroupMember_NotFound(t *testing.T) {
	testApp, groupService := setupTestGroup(t)
	group, _ := createGroup(t, testApp)

	removeOpts := &schema.RemoveGroupMemberOpts{
		ID:     group.ID,
		UserID: 9999,
	}
	err := groupService.RemoveGroupMember(removeOpts)
	assert.Nil(t, err)
}

func RandString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandPhoneNumber() string {
	return fmt.Sprintf("9%09d", rand.Intn(1e9))
}

func RandEmail() string {
	return fmt.Sprintf("%s@example.com", RandString(8))
}
