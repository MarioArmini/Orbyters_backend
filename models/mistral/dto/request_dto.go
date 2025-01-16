package dto

type RequestDto struct {
	UserId         uint   `json:"userId"`
	Inputs         string `json:"inputs"`
	ConversationId *uint  `json:"conversationId"`
}
