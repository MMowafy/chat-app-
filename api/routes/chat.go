package routes

import (
	"instabug-task/api/application"
	"instabug-task/api/controllers"
	"net/http"
)

const ChatResourceName = ApplicationResourceName + "/{token}/chats"

var chatRoutes = []application.Route{
	{http.MethodPost, ChatResourceName,
		func(writer http.ResponseWriter, request *http.Request) {
			controllers.NewChatController().Create(writer, request)
		},
	},
	{http.MethodGet, ChatResourceName,
		func(writer http.ResponseWriter, request *http.Request) {
			controllers.NewChatController().List(writer, request)
		},
	},
	{http.MethodGet, ChatResourceName + "/{chatNumber}",
		func(writer http.ResponseWriter, request *http.Request) {
			controllers.NewChatController().Get(writer, request)
		},
	},
}
