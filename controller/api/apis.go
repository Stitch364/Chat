package api

type apis struct {
	User        user
	Email       email
	Account     account
	Application application
	Message     message
}

var Apis = new(apis)
