package app

import (
	"net/http"

	"github.com/Basu008/EasySplit.git/model"
	"github.com/Basu008/EasySplit.git/schema"
	"gorm.io/gorm"
)

type Payment interface {
	CreatePayment(opts *schema.CreatePaymentOpts) *model.Error

	MigratePayment() error
}

type PaymentImplOpts struct {
	App *App
	DB  *gorm.DB
}

type PaymentImpl struct {
	App *App
	DB  *gorm.DB
}

func InitPayment(opts *PaymentImplOpts) (Payment, error) {
	pi := PaymentImpl{
		App: opts.App,
		DB:  opts.DB,
	}
	err := pi.MigratePayment()
	return &pi, err
}

func (pi *PaymentImpl) CreatePayment(opts *schema.CreatePaymentOpts) *model.Error {
	if err := pi.DB.Transaction(func(tx *gorm.DB) error {
		payment := model.Payment{
			PayerID:   opts.PayerID,
			PayeeID:   opts.PayeeID,
			ExpenseID: opts.ExpenseID,
			Amount:    opts.Amount,
			Mode:      opts.Mode,
		}
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}
		expenseShare := model.ExpenseShare{
			ExpenseID: opts.ExpenseID,
		}
		whereQuery := model.ExpenseShare{
			ExpenseID: opts.ExpenseID,
			UserID:    opts.PayerID,
		}
		updates := map[string]any{
			"is_settled": true,
		}
		if err := tx.Model(&expenseShare).Where(&whereQuery).Updates(updates).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return &model.Error{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return nil
}

func (pi *PaymentImpl) MigratePayment() error {
	return pi.DB.AutoMigrate(&model.Payment{})
}
