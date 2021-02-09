package application

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

var app *application

type application struct {
	router        *chi.Mux
	server        *AppServer              `json:"server"`
	db            []DbConfig              `json:"db"`
	DbConnections map[string]DbConnection `json:"DbConnections"`
	Config        *viper.Viper
	logger        *Logger
}

func NewApplication() *application {
	app = &application{}
	app.setupLogger().
		readConfig().
		setConfig().
		setupDB().
		setupRouter()
	return app
}

func (app *application) readConfig() *application {
	envconf, errenv := getEnvConfig()

	v := viper.New()
	v.SetConfigType("json")

	if errenv == nil {
		readErr := v.ReadConfig(envconf)

		if readErr != nil {
			app.logger.Fatal(fmt.Sprintf("Couldn't read config from OS (APP_CONFIG) .. with error %s", readErr.Error()))
		}

	} else {

		v.AddConfigPath("./")
		v.AddConfigPath("./../")
		readErr := v.ReadInConfig()
		if readErr != nil {
			app.logger.Fatal(fmt.Sprintf("Couldn't read config from config.json File .. make sure that the file exists and it has valid json .. with error %s", readErr.Error()))
		}
		v.SetConfigName("config")
	}

	app.Config = v
	return app

}

func (app *application) setConfig() *application {

	app.server = NewAppServer()
	app.server.Host = GetConfig().GetString("app.host")
	app.server.Port = GetConfig().GetString("app.port")
	app.server.Cors = GetConfig().GetBool("app.cors")
	app.server.HttpLogs = GetConfig().GetBool("app.http_logs")
	app.server.DbLogMode = GetConfig().GetBool("app.db_logs")

	db := GetConfig().Get("db")

	var dbConfig map[string]DbConfig
	var finalDbConfig []DbConfig
	jsonResponse, err := json.Marshal(db)
	if err != nil {
		app.logger.Fatal("failed to marshal db config with err ==> ", err.Error())
	}

	err = json.Unmarshal(jsonResponse, &dbConfig)
	if err != nil {
		app.logger.Fatal("failed to unmarshal db config with err ==> ", err.Error())
	}

	for _, config := range dbConfig {
		finalDbConfig = append(finalDbConfig, config)
	}

	app.db = finalDbConfig
	return app
}

func (app *application) setupLogger() *application {
	app.logger = NewLogger()
	return app
}

func (app *application) setupDB() *application {
	dbConnections := make(map[string]DbConnection)

	for _, dbConfig := range app.db {
		var dbConnection DbConnection
		dbConnection.Config = dbConfig
		dbConnection.Name = dbConfig.Name
		conn, err := setDbConnection(dbConfig)
		if err != nil {
			app.logger.Fatal(err)
		} else {
			dbConnection.Connection = conn
			dbConnections[dbConfig.Name] = dbConnection
		}

	}
	app.DbConnections = dbConnections
	return app
}

func (app *application) setupRouter() *application {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.NoCache)

	if app.server.Cors {
		app.logger.Debug("Cors Enabled")
		// Basic CORS
		// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
		requestCors := cors.New(cors.Options{
			// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "OPTIONS", "DELETE", "PUT", "PATCH"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		})
		r.Use(requestCors.Handler)
	}
	app.router = r
	return app
}

func (app *application) StartServer() {

	app.logger.Infof("Server started http://localhost:%s", app.server.Port)
	err := http.ListenAndServe(":"+app.server.Port, app.router)
	if err != nil {
		app.logger.Fatal(err)
	}
}

func getEnvConfig() (*bytes.Reader, error) {
	envConf := os.Getenv("APP_CONFIG")
	if envConf == "" {
		return nil, errors.New("no Env Variable set")
	}
	envConfByte := []byte(envConf)
	r := bytes.NewReader(envConfByte)

	return r, nil
}

func GetConfig() *viper.Viper {
	return app.Config
}

func GetLogger() *Logger {
	return app.logger
}
