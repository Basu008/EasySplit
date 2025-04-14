package app

import (
	"fmt"
	"os"
)

func InitService(a *App) {
	fmt.Println("\nSetting up Services...")
	var err error
	a.User, err = InitUser(&UserImplOpts{
		App: a,
		DB:  a.Postgres.DB,
	})
	if err != nil {
		os.Exit(1)
		return
	}
	a.Friend, err = InitFriend(&FriendImplOpts{
		App: a,
		DB:  a.Postgres.DB,
	})
	if err != nil {
		os.Exit(1)
		return
	}
	a.Expense, err = InitExpense(&ExpenseImplOpts{
		App: a,
		DB:  a.Postgres.DB,
	})
	if err != nil {
		os.Exit(1)
		return
	}
	a.Group, err = InitGroup(&GroupImplOpts{
		App: a,
		DB:  a.Postgres.DB,
	})
	if err != nil {
		os.Exit(1)
		return
	}
	a.Payment, err = InitPayment(&PaymentImplOpts{
		App: a,
		DB:  a.Postgres.DB,
	})
	if err != nil {
		os.Exit(1)
		return
	}
}
