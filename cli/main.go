package main

import (
	_ "github.com/lib/pq"
	"instabug-task/cli/application"
	"instabug-task/cli/subscribers"
)

func main() {
	app := application.NewApplication()
	go subscribers.SubscribeToCreateChat()
	go subscribers.SubscribeToCreateMessage()
	app.StartServer()
}


