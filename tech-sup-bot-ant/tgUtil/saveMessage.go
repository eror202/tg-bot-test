package tgUtil

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"tg-bot-for-ts/domain"
	"tg-bot-for-ts/repository"
	"tg-bot-for-ts/service"
	"time"
)

func SaveOtherMessage(newUpdate tgbotapi.Update, repo *repository.Repository) string {

	uuid, err := service.MessageCRUD(repo).CreateOtherMessage(domain.NewMessage(newUpdate.Message.Chat.FirstName+newUpdate.Message.Chat.LastName,
		newUpdate.Message.Chat.UserName, newUpdate.Message.Text))
	if err != nil {
		logrus.Panic(err)
	}
	return uuid
}

func SaveIntegrationMessage(newUpdate tgbotapi.Update, repo *repository.Repository) string {

	uuid, err := service.MessageCRUD(repo).CreateIntegrationMessage(domain.NewMessage(newUpdate.Message.Chat.FirstName+" "+newUpdate.Message.Chat.LastName,
		newUpdate.Message.Chat.UserName, newUpdate.Message.Text))
	if err != nil {
		logrus.Panic("Сообщение не сохранилось")
	}
	return uuid
}

func SaveTestMessage(newUpdate tgbotapi.Update, repo *repository.Repository) string {

	uuid, err := service.MessageCRUD(repo).CreateTestMessage(domain.NewMessage(newUpdate.Message.Chat.FirstName+" "+newUpdate.Message.Chat.LastName,
		newUpdate.Message.Chat.UserName, newUpdate.Message.Text))
	if err != nil {
		logrus.Panic("Сообщение не сохранилось")
	}
	return uuid
}

func SaveTrafficMessage(newUpdate tgbotapi.Update, repo *repository.Repository) string {

	uuid, err := service.MessageCRUD(repo).CreateTrafficMessage(domain.NewMessage(newUpdate.Message.Chat.FirstName+" "+newUpdate.Message.Chat.LastName,
		newUpdate.Message.Chat.UserName, newUpdate.Message.Text))
	if err != nil {
		logrus.Panic("Сообщение не сохранилось")
	}
	return uuid
}

func CreateSignature(newUpdate tgbotapi.Update, repo *repository.Repository, message, key string) string {
	today := time.Now()
	_, err := service.MessageCRUD(repo).CreateSignature(domain.NewSignature(key, message, newUpdate.Message.Chat.UserName, today.String()))
	if err != nil {
		logrus.Error("Сигнатура не сохранилась в бд")
	}
	return ""
}
