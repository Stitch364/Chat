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
	File        file
}

var Apis = new(apis)
