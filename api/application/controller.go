package application

import (
	"encoding/json"
	"net/http"
)

type BaseController struct{}

func NewBaseController() *BaseController {
	return &BaseController{}
}

func (baseController *BaseController) json(res http.ResponseWriter, payload interface{}, statusCode int) {
	response, _ := json.Marshal(payload)
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	res.Write(response)
}

func (baseController *BaseController) Json(res http.ResponseWriter, payload interface{}, statusCode int) {
	baseController.json(res, payload, statusCode)
}

func (baseController *BaseController) JsonError(res http.ResponseWriter, msg string, statusCode int) {

	response := map[string]interface{}{
		"status":  statusCode,
		"message": msg,
	}

	baseController.json(res, response, statusCode)
}
