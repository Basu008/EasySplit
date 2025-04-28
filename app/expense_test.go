package app

import (
	"testing"

	"github.com/Basu008/EasySplit.git/model"
	"github.com/Basu008/EasySplit.git/schema"
	"github.com/stretchr/testify/assert"
)

func setupTestExpense(t *testing.T) (*App, Expense) {
	testApp := NewTestApp(getTestConfig())
	t.Cleanup(func() {
		cleanUpDB(testApp.Postgres.DB)
	})

	expenseService, serviceErr := InitExpense(&ExpenseImplOpts{
		App: testApp,
		DB:  testApp.Postgres.DB,
	})
	assert.Nil(t, serviceErr)
	testApp.Expense = expenseService

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
	return testApp, expenseService
}

func createExpenseSetup(t *testing.T, testApp *App, expenseService Expense) (group *schema.GroupResponse, members []*model.User) {
	// Create group and members first
	group, members = createGroupWithMembers(t, testApp, testApp.Group, 2)

	// Prepare expense creation
	memberShares := []schema.MemberIDWithShare{}
	for _, m := range members {
		memberShares = append(memberShares, schema.MemberIDWithShare{
			ID:    m.ID,
			Share: 50, // 50% each
		})
	}

	createExpenseOpts := &schema.CreateExpense{
		GroupID:           group.ID,
		CreatedBy:         members[0].ID,
		Amount:            100,
		Description:       "Dinner Expense",
		ExpenseShareType:  model.Equal,
		MemberIDWithShare: memberShares,
	}

	err := expenseService.CreateExpense(createExpenseOpts)
	assert.Nil(t, err)

	return group, members
}

func createGroupWithMembers(t *testing.T, testApp *App, groupService Group, numMembers int) (*schema.GroupResponse, []*model.User) {
	// Step 1: Create the group owner
	owner := createTestUser(t, testApp)

	// Step 2: Create the member users
	var members []*model.User
	for i := 0; i < numMembers; i++ {
		member := createTestUser(t, testApp)
		members = append(members, member)
	}

	// Step 3: Collect user IDs
	userIDs := []uint{}
	for _, m := range members {
		userIDs = append(userIDs, m.ID)
	}

	// Step 4: Create the group
	createOpts := &schema.CreateGroupOpts{
		Name:    "Group with Members",
		OwnerID: owner.ID,
		Type:    "friends",
		UserIDs: userIDs,
	}
	err := groupService.CreateGroup(createOpts)
	assert.Nil(t, err)

	// Step 5: Fetch the group to get ID
	groups, fetchErr := groupService.GetGroups(owner.ID, 0)
	assert.Nil(t, fetchErr)
	assert.NotEmpty(t, groups)

	return &groups[0], members
}

func TestCreateExpense_OK(t *testing.T) {
	testApp, expenseService := setupTestExpense(t)

	createExpenseSetup(t, testApp, expenseService)
}

func TestGetExpenses_OK(t *testing.T) {
	testApp, expenseService := setupTestExpense(t)

	group, _ := createExpenseSetup(t, testApp, expenseService)

	expenses, err := expenseService.GetExpenses(group.ID)

	assert.Nil(t, err)
	assert.NotNil(t, expenses)
	assert.GreaterOrEqual(t, len(expenses), 1)
}

func TestGetExpense_OK(t *testing.T) {
	testApp, expenseService := setupTestExpense(t)

	group, _ := createExpenseSetup(t, testApp, expenseService)

	expenses, _ := expenseService.GetExpenses(group.ID)
	assert.GreaterOrEqual(t, len(expenses), 1)

	expense := expenses[0]

	fullExpense, err := expenseService.GetExpense(expense.ID)

	assert.Nil(t, err)
	assert.NotNil(t, fullExpense)
	assert.Equal(t, expense.ID, fullExpense.ID)
	assert.GreaterOrEqual(t, len(fullExpense.MembersShare), 1)
}

func TestDeleteExpense_OK(t *testing.T) {
	testApp, expenseService := setupTestExpense(t)

	group, _ := createExpenseSetup(t, testApp, expenseService)

	expenses, _ := expenseService.GetExpenses(group.ID)
	expense := expenses[0]

	ok := expenseService.DeleteExpense(expense.ID)

	assert.True(t, ok)

	// Confirm it's deleted
	_, err := expenseService.GetExpense(expense.ID)
	assert.NotNil(t, err)
}

func TestGetExpenseShare_OK(t *testing.T) {
	testApp, expenseService := setupTestExpense(t)

	group, members := createExpenseSetup(t, testApp, expenseService)

	expenses, _ := expenseService.GetExpenses(group.ID)
	expense := expenses[0]

	share, err := expenseService.GetExpenseShare(expense.ID, members[0].ID)

	assert.Nil(t, err)
	assert.NotNil(t, share)
	assert.Equal(t, expense.ID, share.ExpenseID)
	assert.Equal(t, members[0].ID, share.UserID)
}
