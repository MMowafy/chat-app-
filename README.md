## How to use

```
## To build and start app

docker-compose up --build 

- This command will pull, build image and run container for the following with dependency on each other (will take some time for services to be up and running):  
    - RabbitMQ
    - Redis
    - ElasticSearch
    - MySQL
    - API service (golang APP) for creating messages and chats API's and setup queues
    - CLI (golang App) to listen and consume from message queues  
    - mychatapp (ruby app) for other CRUD operations for Applications, Chats, Messages and migration of DB

- Enhancment to be done:
    - Should be dockerized in a better way 
    - Write unit and integration test
    - Handle rejected messages to a dead letter exchange with retry count
    - Better code and architecture in ruby app 
    - Better dependency injection and use more interfaces in insfrastructure

- What is missing:
    - I added the appropriate mapping and creation of indexes of elasticSearch when go app is up
      but need to add an api to do the search with this query to get exact match first then partial match according to score given
        GET /messages_index/_search
            {
                "query": {
                    "multi_match": {
                        "query": "water",
                        "type": "cross_fields",
                        "tie_breaker": 1,
                        "fields": [
                            "messages.match^10",
                            "messages.exactMatch^50"
                                ]
                    }
                }
            }

    - Add to redis messages count and chats count by publishing to a queue after creation




## hit API

Create App

curl --location --request POST 'http://0.0.0.0:3000/apps' \
--header 'Content-Type: application/json' \
--header 'Content-Type: text/plain' \
--data-raw '{
	"name" : "testapp"
	
}'

List Apps

curl --location --request GET 'http://0.0.0.0:3000/apps?page=1&&limit=10'

Create Chat

curl --location --request POST 'http://localhost:8080/applications/:token/chats' \
--header 'Content-Type: application/json' \
--header 'Content-Type: text/plain' \
--data-raw '{
	
}'

Create Message

curl --location --request POST 'http://localhost:8080/applications/:token/chats/:chatNumber/messages' \
--header 'Content-Type: application/json' \
--header 'Content-Type: text/plain' \
--data-raw '{
	"message": "this is meessage"
}'

List Chats by app token 

curl --location --request GET 'http://0.0.0.0:3000/apps/:token/chats?page=1&&limit=10'

List Messages by app token and chatNumber

curl --location --request GET 'http://0.0.0.0:3000/apps/:token/chats/:number/messages?page=1&&limit=10'

```

