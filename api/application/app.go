package application

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	elasticsearch7 "github.com/olivere/elastic"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"instabug-task/api/models"
	"instabug-task/api/utils"
	"log"
	"net/http"
	"os"
)

var app *application

var chatMapping = `{
  "settings": {
    "analysis": {
      "analyzer": {
        "gramAnalyzer": {
          "tokenizer": "gramTokenizer",
          "filter": [
            "lowercase"
          ]
        }
      },
      "tokenizer": {
        "gramTokenizer": {
          "type": "edge_ngram",
          "min_gram": 2,
          "max_gram": 10
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "message": {
        "type": "keyword",
        "fields": {
          "match": {
            "type": "text",
            "analyzer": "gramAnalyzer",
            "search_analyzer": "gramAnalyzer"
          }
        }
      }
    }
  }
}
`

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

func (app *application) MigrateAndSeedDB() {

	conn, _ := GetMysqlConnectionByName("appdb")
	GetLogger().Info("Migration started ")

	db := conn.AutoMigrate(&models.Chat{})
	if db.Error != nil {
		GetLogger().Fatal("Error while migrating CoffeeProduct table => ", db.Error.Error())
	}

	db = conn.AutoMigrate(&models.Message{})
	if db.Error != nil {
		GetLogger().Fatal("Error while migrating CoffeeProduct table => ", db.Error.Error())
	}

	db = conn.Model(&models.Message{}).AddForeignKey("chatId", "chats(id)", "RESTRICT", "RESTRICT")
	if db.Error != nil {
		GetLogger().Error("Error while creating chat foreign key in messages => ", db.Error.Error())
	}

	GetLogger().Info("Migration Ended ")
}

func (app *application) SetupQueues() {
	conn, _ := GetRabbitmqConnectionByName("rabbitmq")
	app.setupChatExchangeAndQueues(conn)
	app.setupMessageExchangeAndQueues(conn)
}

func (app *application) setupChatExchangeAndQueues(conn *amqp.Connection) {
	ch, err := conn.Channel()

	if err != nil {
		GetLogger().Fatal(err.Error())
	}
	defer ch.Close()

	// setup exchanges

	err = app.setupExchange(ch, utils.ChatExchange, "direct")

	if err != nil {
		GetLogger().Fatal(err.Error())
	}

	//setup queues

	q, err := app.setupQueue(ch, utils.ChatQueue)

	if err != nil {
		GetLogger().Fatal(err.Error())
	}

	err = ch.QueueBind(
		q.Name,                       // queue name
		string(utils.ChatRoutingKey), // routing key
		string(utils.ChatExchange),   // exchange
		false,
		nil)

	if err != nil {
		GetLogger().Fatal(err.Error())
	}
	GetLogger().Info("Setup for chat queue and exchange done successfully")
}

func (app *application) setupMessageExchangeAndQueues(conn *amqp.Connection) {
	ch, err := conn.Channel()

	if err != nil {
		GetLogger().Fatal(err.Error())
	}
	defer ch.Close()

	// setup exchanges

	err = app.setupExchange(ch, utils.MessageExchange, "direct")

	if err != nil {
		GetLogger().Fatal(err.Error())
	}

	//setup queues

	q, err := app.setupQueue(ch, utils.MessageQueue)

	if err != nil {
		GetLogger().Fatal(err.Error())
	}

	err = ch.QueueBind(
		q.Name,                          // queue name
		string(utils.MessageRoutingKey), // routing key
		string(utils.MessageExchange),   // exchange
		false,
		nil)

	if err != nil {
		GetLogger().Fatal(err.Error())
	}
	GetLogger().Info("Setup for message queue and exchange done successfully")
}

func (app *application) setupExchange(ch *amqp.Channel, exchangeName utils.ExchangeType, exchangeType string) error {
	err := ch.ExchangeDeclare(
		string(exchangeName), // name
		exchangeType,         // type
		true,                 // durable
		false,                // auto-deleted
		false,                // internal
		false,                // no-wait
		nil,                  // arguments
	)

	return err
}

func (app *application) setupQueue(ch *amqp.Channel, queueName utils.QueueType) (amqp.Queue, error) {

	args := amqp.Table{
		// further enhancement to define dlx
	}

	q, err := ch.QueueDeclare(
		string(queueName), // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		args,              // arguments
	)

	return q, err
}

func (app *application) SetupIndexes() {
	conn, err := app.GetDbConnectionByName("elastic_search")
	if err != nil {
		log.Panic(err)
	}
	client := conn.Connection.(*elasticsearch7.Client)
	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("message_index").Do(context.Background())
	if err != nil {
		// Handle error
		GetLogger().Panic(err)
	}
	//if index does not exist, create a new one with the specified mapping
	if !exists {
		createIndex, err := client.CreateIndex("message_index").BodyString(chatMapping).Do(context.Background())
		if err != nil {
			panic(err)
		}
		if !createIndex.Acknowledged {
			log.Println(createIndex)
		} else {
			log.Println("successfully created index")
		}
	} else {
		log.Println("Index already exist")
	}
}

func (app *application) SetRoutes(routes []Route) *application {

	for _, route := range routes {
		app.router.MethodFunc(route.Method, route.Pattern, route.HandlerFunc)
	}
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
