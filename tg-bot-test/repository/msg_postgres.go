package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"tg-bot-for-ts/domain"
)

type MsgPostgres struct {
	db *sqlx.DB
}

func (r *MsgPostgres) CreateSignature(signature *domain.Signature) (string, error) {
	var id string
	query := fmt.Sprintf("INSERT INTO %s (publicKey, body, login, timeOfSignature) VALUES ($1,$2,$3,$4) RETURNING ID", signatureTable)
	row := r.db.QueryRow(query, signature.PublicKey, signature.Body, signature.Login, signature.TimeOfSignature)
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (r *MsgPostgres) CreateIntegrationMessage(msg *domain.Message) (string, error) {
	var id string
	query := fmt.Sprintf("INSERT INTO %s (person, login, request) VALUES ($1,$2,$3) RETURNING ID", integrationTable)
	row := r.db.QueryRow(query, msg.Person, msg.Login, msg.Request)
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (r *MsgPostgres) CreateTestMessage(msg *domain.Message) (string, error) {
	var id string
	query := fmt.Sprintf("INSERT INTO %s (person, login, request) VALUES ($1,$2,$3) RETURNING ID", testTable)
	row := r.db.QueryRow(query, msg.Person, msg.Login, msg.Request)
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (r *MsgPostgres) CreateTrafficMessage(msg *domain.Message) (string, error) {
	var id string
	query := fmt.Sprintf("INSERT INTO %s (person, login, request) VALUES ($1,$2,$3) RETURNING ID", trafficTable)
	row := r.db.QueryRow(query, msg.Person, msg.Login, msg.Request)
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (r *MsgPostgres) CreateOtherMessage(msg *domain.Message) (string, error) {
	var id string
	query := fmt.Sprintf("INSERT INTO %s (person, login, request) VALUES ($1,$2,$3) RETURNING ID", otherTable)
	row := r.db.QueryRow(query, msg.Person, msg.Login, msg.Request)
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func NewMsgPostgres(db *sqlx.DB) *MsgPostgres {
	return &MsgPostgres{db: db}
}

/*func (r *MsgPostgres) CreateMessage(msg *domain.Message) (string, error) {
	var id string
	query := fmt.Sprintf("INSERT INTO %s (person, login, request) VALUES ($1,$2,$3) RETURNING ID", messageTable)
	row := r.db.QueryRow(query, msg.Person, msg.Login, msg.Request)
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}*/

/*func (r *MsgPostgres) GetAllMessage() (*[]domain.Message, error) {
	var list []domain.Message
	query := fmt.Sprintf("SELECT * FROM %s", messageTable)
	err := r.db.Select(&list, query)
	return &list, err
}*/
