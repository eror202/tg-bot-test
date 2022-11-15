package dto

type MessageDto struct {
	FullName string `json:"name"`
	Login    string `json:"login"`
	Request  string `json:"request"`
}
