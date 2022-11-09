package service

import (
	"tg-bot-for-ts/domain"
	"tg-bot-for-ts/repository"
)

type MessageService struct {
	repo repository.MessageRepository
}

func (s *MessageService) CreateIntegrationMessage(message *domain.Message) (string, error) {
	return s.CreateIntegrationMessage(message)
}

func (s *MessageService) CreateTestMessage(message *domain.Message) (string, error) {
	return s.CreateTestMessage(message)
}

func (s *MessageService) CreateTrafficMessage(message *domain.Message) (string, error) {
	return s.CreateTrafficMessage(message)
}

func (s *MessageService) CreateOtherMessage(message *domain.Message) (string, error) {
	return s.CreateOtherMessage(message)
}

func (s *MessageService) CreateSignature(signature *domain.Signature) (string,error) {
	return s.CreateSignature(signature)
}
func NewMessageService(repo repository.MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}
