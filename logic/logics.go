package logic

type logics struct {
	User        user
	Email       email
	Account     account
	Application application
	Message     message
	Setting     setting
	Group       group
}

var Logics = new(logics)
