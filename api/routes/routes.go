package routes

import "instabug-task/api/application"

func GetRoutes() []application.Route {
	var appRoutes []application.Route
	appRoutes = append(appRoutes, chatRoutes...)
	appRoutes = append(appRoutes, messageRoutes...)
	return appRoutes
}
