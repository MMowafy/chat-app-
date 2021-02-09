package infrastructure

import (
	"context"
	elasticsearch7 "github.com/olivere/elastic"
	"instabug-task/cli/application"
	"instabug-task/cli/models"
)

type MessageIndexService struct {
	client *elasticsearch7.Client
}

func NewMessageIndexService() *MessageIndexService {
	return &MessageIndexService{
		application.GetElasticSearchClient(),
	}
}

func (messageIndexService *MessageIndexService) IndexDoc(message *models.Message) {
	_, err := messageIndexService.client.Index().Index("message_index").Type("_doc").BodyJson(message).Do(context.Background())
	if err != nil {
		application.GetLogger().Error("Failed to create doc in message index with err => ", err.Error())
	}
}
