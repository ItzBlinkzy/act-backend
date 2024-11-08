package model

type RequestDataChatBot struct {
	Message string `json:"message"`
	Context string `json:"context"`
}

type ChatGPTResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
