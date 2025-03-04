package app

import "os"

func InitService(a *App) {
	var err error
	a.User, err = InitUser(&UserImplOpts{
		App: a,
		DB:  a.Postgres.DB,
	})
	if err != nil {
		os.Exit(1)
		return
	}
}
