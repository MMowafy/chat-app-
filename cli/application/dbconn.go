package application

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	elasticsearch7 "github.com/olivere/elastic"
	amqp "github.com/streadway/amqp"
	"log"
	"strconv"
)

type DbConfig struct {
	Name     string `json:"name"`
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
	DSN      string `json:"dsn"`
}

type DbConnection struct {
	Config     DbConfig
	Name       string
	Connection interface{}
}

func (app *application) GetDbConnectionByName(connName string) (DbConnection, error) {

	if val, ok := app.DbConnections[connName]; ok {
		return val, nil
	}

	return DbConnection{}, errors.New("connection is not defined, connection name: " + connName)
}

func GetMysqlConnectionByName(connName string) (*gorm.DB, error) {

	dbConn, err := app.GetDbConnectionByName(connName)

	if err != nil {
		GetLogger().Error(err.Error())
		return nil, err
	}

	if dbConn.Connection == nil {
		errMsg := "connection error, connection name: " + connName
		GetLogger().Error(errMsg)
		return nil, errors.New(errMsg)
	}

	return dbConn.Connection.(*gorm.DB), nil
}

func GetRedisConnectionByName(connectionName string) (redis.Cmdable, error) {
	conn, err := app.GetDbConnectionByName(connectionName)

	if err != nil {
		return nil, err
	}

	if conn.Connection == nil {
		return nil, errors.New("Can not find opened redis connection")
	}

	return conn.Connection.(redis.Cmdable), nil
}

func GetRabbitmqConnectionByName(connectionName string) (*amqp.Connection, error) {
	dbConn, err := app.GetDbConnectionByName(connectionName)
	if dbConn.Connection == nil {
		GetLogger().Fatal(err.Error())
		return nil, err
	}
	conn := dbConn.Connection.(*amqp.Connection)
	return conn, nil
}

func GetElasticSearchClient() *elasticsearch7.Client {
	conn, err := app.GetDbConnectionByName("elastic_search")
	if err != nil {
		log.Panic(err)
	}
	return conn.Connection.(*elasticsearch7.Client)
}

func setDbConnection(dbConfig DbConfig) (interface{}, error) {
	if dbConfig.Driver == "mysql" {
		return setMysqlDbConnection(dbConfig)
	} else if dbConfig.Driver == "redis" {
		return setRedisConnection(dbConfig)
	} else if dbConfig.Driver == "rabbitmq" {
		return setRabbitMqConnection(dbConfig)
	} else if dbConfig.Driver == "elastic" {
		return setElasticSearchConnection(dbConfig)
	}
	return nil, errors.New("Can't handle connection to driver " + dbConfig.Driver)
}

func setMysqlDbConnection(dbConfig DbConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DbName)
	app.logger.Infof("Connecting to mysql db at %s", dsn)

	db, err := gorm.Open("mysql", dsn)

	if err != nil {
		app.logger.Fatal("Can't connect to mysql db error ", err.Error())
		return nil, err
	}
	app.logger.Info("Connected to mysql db successfully at ", dsn)

	return db, nil
}

func setRedisConnection(dbConfig DbConfig) (redis.Cmdable, error) {
	dbName, _ := strconv.Atoi(dbConfig.DbName)

	app.logger.Info("Connecting to Redis db at %", dbConfig.DSN)
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     dbConfig.DSN,
		Password: "",
		DB:       dbName,
	})

	_, err := RedisClient.Ping(context.TODO()).Result()
	if err != nil {
		app.logger.Fatalf("Can't connect to redis server %s error %s", dbConfig.DSN, err)
		return nil, err
	}
	app.logger.Infof("Connected to Redis successfully at %s", dbConfig.DSN)

	return RedisClient, nil
}

//connect to rabbitMq
func setRabbitMqConnection(dbConfig DbConfig) (*amqp.Connection, error) {

	app.logger.Info("Connecting to RabbitMq db at %", dbConfig.DSN)

	conn, err := amqp.Dial(dbConfig.DSN)
	if err != nil {
		app.logger.Fatalf("can't connect to RabbitMq %s error %s", dbConfig.DSN, err)
		return nil, err
	}
	app.logger.Info("Connected to RabbitMq successfully at %", dbConfig.DSN)
	return conn, nil
}

func setElasticSearchConnection(dbConfig DbConfig) (*elasticsearch7.Client, error) {

	app.logger.Info("Connecting to elasticsearch at %", dbConfig.DSN)

	client, err := elasticsearch7.NewClient(
		elasticsearch7.SetURL(dbConfig.DSN),
		elasticsearch7.SetSniff(false),
		elasticsearch7.SetHealthcheck(false))
	if err != nil {
		app.logger.Fatalf("can't connect to ElasticSearch %s error %s", dbConfig.DSN, err)
		return nil, err
	}

	app.logger.Info("Connected to Elasticsearch successfully at %", dbConfig.DSN)
	return client, nil
}
