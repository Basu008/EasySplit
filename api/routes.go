package api

func (a *API) InitRoutes() {

	//User
	a.Router.APIRoot.Handle("/user/login", a.requestHandler(a.loginUser)).Methods("POST")
	a.Router.APIRoot.Handle("/user/confirm-otp", a.requestHandler(a.confirmOTP)).Methods("POST")
	a.Router.APIRoot.Handle("/user", a.requestWithAuthHandler(a.updateUser)).Methods("PATCH")
	a.Router.APIRoot.Handle("/user", a.requestHandler(a.getUser)).Methods("GET")

	//Friend
	a.Router.APIRoot.Handle("/friend", a.requestWithAuthHandler(a.sendFriendRequest)).Methods("POST")
	a.Router.APIRoot.Handle("/friend", a.requestWithAuthHandler(a.updateFriendRequest)).Methods("PATCH")
	a.Router.APIRoot.Handle("/friends", a.requestWithAuthHandler(a.allFriends)).Methods("GET")
	a.Router.APIRoot.Handle("/friend/{{id}}", a.requestWithAuthHandler(a.getFriendStatus)).Methods("GET")

	//To check whether the system is up and running
	a.Router.Root.Handle("/health-check", a.requestHandler(a.healthCheck)).Methods("GET")
}
