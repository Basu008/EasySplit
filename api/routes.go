package api

func (a *API) InitRoutes() {

	//To check whether the system is up and running
	a.Router.Root.Handle("/health-check", a.requestHandler(a.healthCheck)).Methods("GET")

	//User
	a.Router.APIRoot.Handle("/user/login", a.requestHandler(a.loginUser)).Methods("POST")
	a.Router.APIRoot.Handle("/user/confirm-otp", a.requestHandler(a.confirmOTP)).Methods("POST")
	a.Router.APIRoot.Handle("/user", a.requestWithAuthHandler(a.updateUser)).Methods("PATCH")
	a.Router.APIRoot.Handle("/user", a.requestHandler(a.getUser)).Methods("GET")

	//Friend
	a.Router.APIRoot.Handle("/friend", a.requestWithAuthHandler(a.sendFriendRequest)).Methods("POST")
	a.Router.APIRoot.Handle("/friend", a.requestWithAuthHandler(a.updateFriendRequest)).Methods("PATCH")
	a.Router.APIRoot.Handle("/friends", a.requestWithAuthHandler(a.allFriends)).Methods("GET")
	a.Router.APIRoot.Handle("/friend/{id}", a.requestWithAuthHandler(a.getFriendStatus)).Methods("GET")

	//Group
	a.Router.APIRoot.Handle("/group", a.requestWithAuthHandler(a.createGroup)).Methods("POST")
	a.Router.APIRoot.Handle("/group/{id}", a.requestWithAuthHandler(a.getGroupByID)).Methods("GET")
	a.Router.APIRoot.Handle("/groups", a.requestWithAuthHandler(a.getGroups)).Methods("GET")
	a.Router.APIRoot.Handle("/group/{id}", a.requestWithAuthHandler(a.editGroup)).Methods("PATCH")

	a.Router.APIRoot.Handle("/group/members", a.requestWithAuthHandler(a.addGroupMembers)).Methods("POST")
	a.Router.APIRoot.Handle("/group/member", a.requestWithAuthHandler(a.removeGroupMember)).Methods("DELETE")

	//Expense
	a.Router.APIRoot.Handle("/expense", a.requestWithAuthHandler(a.createExpense)).Methods("POST")
	a.Router.APIRoot.Handle("/expenses/{group_id}", a.requestWithAuthHandler(a.getExpenses)).Methods("GET")
	a.Router.APIRoot.Handle("/expense/{expense_id}", a.requestWithAuthHandler(a.getExpense)).Methods("GET")
	a.Router.APIRoot.Handle("/expense/{expense_id}", a.requestWithAuthHandler(a.deleteExpense)).Methods("DELETE")

	//Payment
	a.Router.APIRoot.Handle("/payment", a.requestWithAuthHandler(a.createPayment)).Methods("POST")
}
