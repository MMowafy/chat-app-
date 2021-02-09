package responses

type CreateChatResponse struct {
	ChatNumber int `json:"chatNumber"`
}

func NewCreateChatResponse() *CreateChatResponse {
	return &CreateChatResponse{}
}
