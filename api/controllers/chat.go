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

type ChatController struct {
	*application.BaseController
	chatService *services.ChatService
}

func NewChatController() *ChatController {
	return &ChatController{
		application.NewBaseController(),
		services.NewChatService(),
	}
}

func (chatController *ChatController) Create(w http.ResponseWriter, r *http.Request) {

	chatRequest := models.NewChat()
	appToken := chi.URLParam(r, "token")
	if err := json.NewDecoder(r.Body).Decode(&chatRequest); err != nil {
		application.GetLogger().Error(err.Error())
		chatController.JsonError(w, utils.ErrorInvalidRequestPayload, http.StatusBadRequest)
		return
	}


	response, err := chatController.chatService.Create(appToken,chatRequest)
	if err != nil {
		chatController.JsonError(w, utils.ErrorCreateChat, http.StatusBadRequest)
		return
	}

	chatController.Json(w, response, http.StatusOK)
}

func (chatController *ChatController) List(w http.ResponseWriter, r *http.Request) {

	//appToken := chi.URLParam(r, "token")
	//
	//list, err := chatController.chatService.List(r, appToken)
	//if err != nil {
	//	chatController.JsonError(w, "Failed to get chat list", http.StatusBadRequest)
	//	return
	//}
	//
	//chatController.Json(w, list, http.StatusOK)
}

func (chatController *ChatController) Get(w http.ResponseWriter, r *http.Request) {

	//appToken := chi.URLParam(r, "token")
	//chatNumber := chi.URLParam(r, "chatNumber")
	//number, _ := strconv.Atoi(chatNumber)
	//// validate app token
	//
	//chat, err := chatController.chatService.GetChatByTokenAndNumber(appToken, number)
	//if err != nil {
	//	chatController.JsonError(w, "Failed to get chat list", http.StatusBadRequest)
	//	return
	//}
	//
	//chatController.Json(w, chat, http.StatusOK)
}
