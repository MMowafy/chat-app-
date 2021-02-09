package routes

import (
	"instabug-task/api/application"
	"instabug-task/api/controllers"
	"net/http"
)

const MessageResourceName = ChatResourceName + "/{chatNumber}/messages"

var messageRoutes = []application.Route{
	{http.MethodPost, MessageResourceName,
		func(writer http.ResponseWriter, request *http.Request) {
			controllers.NewMessageController().Create(writer, request)
		},
	},
	{http.MethodGet, MessageResourceName,
		func(writer http.ResponseWriter, request *http.Request) {
			controllers.NewMessageController().List(writer, request)
		},
	},
}
