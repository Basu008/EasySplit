package app

import (
	"testing"

	"github.com/Basu008/EasySplit.git/model"
	"github.com/Basu008/EasySplit.git/schema"
	"github.com/stretchr/testify/assert"
)

func setupTestPayment(t *testing.T) (*App, Payment) {
	testApp := NewTestApp(getTestConfig())
	t.Cleanup(func() {
		cleanUpDB(testApp.Postgres.DB)
	})

	paymentService, serviceErr := InitPayment(&PaymentImplOpts{
		App: testApp,
		DB:  testApp.Postgres.DB,
	})
	assert.Nil(t, serviceErr)
	testApp.Payment = paymentService
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
	return testApp, paymentService
}

func createPaymentSetup(t *testing.T, testApp *App) (*model.Expense, *model.User, *model.User) {
	// Step 1: Create a group with members
	group, members := createGroupWithMembers(t, testApp, testApp.Group, 2)

	// Step 2: Create an expense inside that group
	memberShares := []schema.MemberIDWithShare{
		{ID: members[0].ID, Share: 50},
		{ID: members[1].ID, Share: 50},
	}

	createExpenseOpts := &schema.CreateExpense{
		GroupID:           group.ID,
		CreatedBy:         members[0].ID,
		Amount:            200,
		Description:       "Group Dinner",
		ExpenseShareType:  model.Equal,
		MemberIDWithShare: memberShares,
	}
	err := testApp.Expense.CreateExpense(createExpenseOpts)
	assert.Nil(t, err)

	// Step 3: Get the expense ID
	expenses, _ := testApp.Expense.GetExpenses(group.ID)
	assert.GreaterOrEqual(t, len(expenses), 1)

	expense := expenses[0]

	return &expense, members[0], members[1]
}

func TestCreatePayment_OK(t *testing.T) {
	testApp, paymentService := setupTestPayment(t)

	expense, payer, payee := createPaymentSetup(t, testApp)

	createPaymentOpts := &schema.CreatePaymentOpts{
		PayerID:   payer.ID,
		PayeeID:   payee.ID,
		ExpenseID: expense.ID,
		Amount:    100,
		Mode:      "upi",
	}

	err := paymentService.CreatePayment(createPaymentOpts)

	assert.Nil(t, err)

	// Optional: Check if ExpenseShare isSettled
	expenseShare, shareErr := testApp.Expense.GetExpenseShare(expense.ID, payer.ID)
	assert.Nil(t, shareErr)
	assert.True(t, expenseShare.IsSettled)
}
