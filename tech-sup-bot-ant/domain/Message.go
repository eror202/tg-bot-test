package domain

type Message struct {
	ID      string `json:"-" db:"id"`
	Person  string `json:"person" binding:"required"`
	Login   string `json:"login" binding:"required"`
	Request string `json:"request" binding:"required"`
}

func NewMessage(person string, login string, request string) *Message {
	return &Message{Person: person, Login: login, Request: request}
}
