package logic

type logics struct {
	User        user
	Email       email
	Account     account
	Application application
	Message     message
	Setting     setting
}

var Logics = new(logics)
