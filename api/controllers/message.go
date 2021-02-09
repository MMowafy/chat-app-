package controllers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"instabug-task/api/application"
	"instabug-task/api/models"
	"instabug-task/api/services"
	"instabug-task/api/utils"
	"net/http"
)

type MessageController struct {
	*application.BaseController
	MessageService *services.MessageService
}

func NewMessageController() *MessageController {
	return &MessageController{
		application.NewBaseController(),
		services.NewMessageService(),
	}
}
func (messageController *MessageController) Create(w http.ResponseWriter, r *http.Request) {

	appToken := chi.URLParam(r, "token")
	chatNumber := chi.URLParam(r, "chatNumber")
	messageRequest := models.NewMessage()
	if err := json.NewDecoder(r.Body).Decode(&messageRequest); err != nil {
		application.GetLogger().Error(err.Error())
		messageController.JsonError(w, utils.ErrorInvalidRequestPayload, http.StatusBadRequest)
		return
	}

	response, err := messageController.MessageService.Create(appToken,chatNumber,messageRequest)
	if err != nil {
		messageController.JsonError(w, utils.ErrorCreateMessage, http.StatusBadRequest)
		return
	}

	messageController.Json(w, response, http.StatusOK)
}


func (messageController *MessageController) List(w http.ResponseWriter, r *http.Request) {

	//appToken := chi.URLParam(r, "token")
	//chatNumber := chi.URLParam(r, "chatNumber")
	//number, _ := strconv.Atoi(chatNumber)
	//
	//list, err := messageController.MessageService.ListMessagesByChatNumber(r,appToken, number)
	//if err != nil {
	//	messageController.JsonError(w, "Failed to get Message list", http.StatusBadRequest)
	//	return
	//}
	//
	//messageController.Json(w, list, http.StatusOK)
}

