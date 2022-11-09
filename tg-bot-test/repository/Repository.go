package repository

import (
	"github.com/jmoiron/sqlx"
	"tg-bot-for-ts/domain"
)

type MessageRepository interface {
	CreateIntegrationMessage(message *domain.Message) (string, error)
	CreateTestMessage(message *domain.Message) (string, error)
	CreateTrafficMessage(message *domain.Message) (string, error)
	CreateOtherMessage(message *domain.Message) (string, error)
	CreateSignature(signature *domain.Signature) (string, error)
}

type Repository struct {
	MessageRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		MessageRepository: NewMsgPostgres(db),
	}
}
