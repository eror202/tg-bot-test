package service

import (
	"tg-bot-for-ts/domain"
	"tg-bot-for-ts/repository"
)

type MessageCRUD interface {
	CreateIntegrationMessage(message *domain.Message) (string, error)
	CreateTestMessage(message *domain.Message) (string, error)
	CreateTrafficMessage(message *domain.Message) (string, error)
	CreateOtherMessage(message *domain.Message) (string, error)
	CreateSignature(signature *domain.Signature) (string, error)
}

type Service struct {
	MessageCRUD
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		MessageCRUD: NewMessageService(repo.MessageRepository),
	}
}
