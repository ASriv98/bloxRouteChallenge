package data

import "github.com/google/uuid"

type Message struct {
	ID     string      `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

func NewMessage(method string, item interface{}) Message {
	return Message{
		ID:     uuid.New().String(),
		Method: method,
		Params: item,
	}
}
