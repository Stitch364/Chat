package routs

type routers struct {
	User        user
	Email       email
	Account     account
	Application application
	Message     message
	Setting     setting
}

var Routers = new(routers)
