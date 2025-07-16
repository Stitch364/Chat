package routs

type routers struct {
	User        user
	Email       email
	Account     account
	Application application
	Message     message
	Setting     setting
	Chat        ws
	Group       group
}

var Routers = new(routers)
