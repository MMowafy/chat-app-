package main

import (
	_ "github.com/lib/pq"
	"instabug-task/api/application"
	"instabug-task/api/routes"
)

func main() {
	app := application.NewApplication()
	//app.MigrateAndSeedDB()
	app.SetupQueues()
	app.SetupIndexes()
	app.SetRoutes(routes.GetRoutes())
	app.StartServer()
}
