package app

import (
	"github.com/Basu008/EasySplit.git/model"
	"gorm.io/gorm"
)

type Expense interface {
	MigrateExpense() error
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
	// if err := ei.MigrateExpense(); err != nil {
	// 	log.Print(err)
	// 	return nil, errors.New("unable to migrate Expense")
	// }
	return &ei, nil
}

func (ei *ExpenseImpl) MigrateExpense() error {
	err := ei.DB.AutoMigrate(&model.Expense{})
	return err
}
