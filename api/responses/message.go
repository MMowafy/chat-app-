package responses

type CreateMessageResponse struct {
	MessageNumber int `json:"messageNumber"`
}

func NewCreateMessageResponse() *CreateMessageResponse {
	return &CreateMessageResponse{}
}
