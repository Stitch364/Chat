package api

import "chat/controller/api/chat"

type apis struct {
	User        user
	Email       email
	Account     account
	Application application
	Message     message
	Setting     setting
	Chat        chat.Group
	Group       group
}

var Apis = new(apis)
