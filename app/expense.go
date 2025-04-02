package app

import (
	"errors"
	"log"
	"net/http"

	"github.com/Basu008/EasySplit.git/model"
	"github.com/Basu008/EasySplit.git/schema"
	"gorm.io/gorm"
)

type Expense interface {
	CreateExpense(opt *schema.CreateExpense) *model.Error
	GetExpenses(groupID uint) ([]model.Expense, *model.Error)
	GetExpense(expenseID uint) (*model.ExpenseWithShares, *model.Error)
	DeleteExpense(id uint) bool

	MigrateExpense() error
	MigrateExpenseShare() error
}

type ExpenseImplOpts struct {
	App *App
	DB  *gorm.DB
}

type ExpenseImpl struct {
	App *App
	DB  *gorm.DB
}

func InitExpense(opts *ExpenseImplOpts) (Expense, error) {
	ei := ExpenseImpl{
		App: opts.App,
		DB:  opts.DB,
	}
	if err := ei.MigrateExpense(); err != nil {
		log.Print(err)
		return nil, errors.New("unable to migrate Expense")
	}
	if err := ei.MigrateExpenseShare(); err != nil {
		log.Print(err)
		return nil, errors.New("unable to migrate Expense")
	}
	return &ei, nil
}

func (ei *ExpenseImpl) CreateExpense(opts *schema.CreateExpense) *model.Error {
	tx := ei.DB.Begin()
	expense := model.Expense{
		GroupID:     opts.GroupID,
		CreatedBy:   opts.CreatedBy,
		Amount:      opts.Amount,
		Description: opts.Description,
	}
	if err := tx.Create(&expense).Error; err != nil {
		tx.Rollback()
		return &model.Error{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	expenseShares := calculateShares(expense.ID, opts)
	if err := tx.CreateInBatches(expenseShares, len(expenseShares)).Error; err != nil {
		tx.Rollback()
		return &model.Error{
			Err:  err,
			Code: http.StatusInternalServerError,
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

func (ei *ExpenseImpl) GetExpenses(groupID uint) ([]model.Expense, *model.Error) {
	var expenses []model.Expense
	err := ei.DB.Find(&expenses).Where("group_id = ?", groupID).Error
	if err != nil {
		return nil, &model.Error{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return expenses, nil
}

func (ei *ExpenseImpl) GetExpense(id uint) (*model.ExpenseWithShares, *model.Error) {
	expense := model.Expense{
		ID: id,
	}
	if err := ei.DB.First(&expense).Error; err != nil {
		return nil, &model.Error{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}
	var expenseShare []model.ExpenseShare
	if err := ei.DB.Find(&expenseShare).Error; err != nil {
		return nil, &model.Error{
			Err:  err,
			Code: http.StatusBadRequest,
		}
	}
	resp := model.ExpenseWithShares{
		ID:          expense.ID,
		TotalAmount: expense.Amount,
		Description: expense.Description,
		CreatedBy:   expense.CreatedBy,
	}
	shares := []model.MembersShare{}
	for _, es := range expenseShare {
		ms := model.MembersShare{
			ID:     es.UserID,
			Amount: es.Amount,
		}
		shares = append(shares, ms)
	}
	resp.MembersShare = shares
	return &resp, nil
}

func (ei *ExpenseImpl) DeleteExpense(id uint) bool {
	tx := ei.DB.Begin()
	expense := model.Expense{
		ID: id,
	}
	var expenseShare []model.ExpenseShare
	if err := tx.Where("expense_id = ?", id).Delete(&expenseShare).Error; err != nil {
		tx.Rollback()
		return false
	}
	if err := tx.Delete(&expense).Error; err != nil {
		tx.Rollback()
		return false
	}
	if err := tx.Commit().Error; err != nil {
		return false
	}
	return true
}

func (ei *ExpenseImpl) MigrateExpense() error {
	err := ei.DB.AutoMigrate(&model.Expense{})
	return err
}

func (ei *ExpenseImpl) MigrateExpenseShare() error {
	err := ei.DB.AutoMigrate(&model.ExpenseShare{})
	return err
}

func calculateShares(expenseID uint, opts *schema.CreateExpense) []model.ExpenseShare {
	switch opts.ExpenseShareType {
	case model.Equal:
		return calculateEqualShare(expenseID, opts)
	case model.Percent:
		return calculatePercentShare(expenseID, opts)
	default:
		return calculateCustomShare(expenseID, opts)
	}
}

func calculatePercentShare(expenseID uint, opts *schema.CreateExpense) []model.ExpenseShare {
	expenseShares := []model.ExpenseShare{}
	amount := opts.Amount
	totalShare := amount
	membersWithShares := opts.MemberIDWithShare
	if opts.UserShare != 0 {
		userShare := (amount * opts.UserShare) / 100
		amount -= userShare
		totalShare = amount
	}
	for i, memberWithShare := range membersWithShares {
		var value float64
		if i == len(membersWithShares)-1 {
			value = totalShare
		} else {
			value = (amount * memberWithShare.Share) / 100
		}
		expenseShares = append(expenseShares, model.ExpenseShare{
			ExpenseID: expenseID,
			UserID:    memberWithShare.ID,
			Amount:    value,
		})
		totalShare -= value
	}
	return expenseShares
}

func calculateCustomShare(expenseID uint, opts *schema.CreateExpense) []model.ExpenseShare {
	expenseShares := []model.ExpenseShare{}
	membersWithShares := opts.MemberIDWithShare
	for _, memberWithShare := range membersWithShares {
		expenseShares = append(expenseShares, model.ExpenseShare{
			ExpenseID: expenseID,
			UserID:    memberWithShare.ID,
			Amount:    memberWithShare.Share,
		})
	}
	return expenseShares
}

func calculateEqualShare(expenseID uint, opts *schema.CreateExpense) []model.ExpenseShare {
	expenseShares := []model.ExpenseShare{}
	amount := opts.Amount
	totalShare := amount
	membersWithShares := opts.MemberIDWithShare
	if opts.UserShare != 0 {
		amount -= opts.UserShare
		totalShare = amount
	}
	for i, mememberWithShare := range membersWithShares {
		var value float64
		if i == len(membersWithShares)-1 {
			value = totalShare
		} else {
			value = amount / float64(len(membersWithShares))
		}
		expenseShares = append(expenseShares, model.ExpenseShare{
			ExpenseID: expenseID,
			UserID:    mememberWithShare.ID,
			Amount:    value,
		})
		totalShare -= value
	}
	return expenseShares
}
