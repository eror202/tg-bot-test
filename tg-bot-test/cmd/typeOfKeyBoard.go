package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var KeyBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Интеграция", "Интеграция"),
		tgbotapi.NewInlineKeyboardButtonData("Тесты", "Тесты")),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Трафик", "Трафик"),
		tgbotapi.NewInlineKeyboardButtonData("Другое", "Другое")),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Главное меню", "Меню")))

var PostKeyBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отправить запрос", "Запрос отправлен")))

var StartKeyBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Генерация подписи", "Генерация подписи"),
		tgbotapi.NewInlineKeyboardButtonData("Обратиться в тех поддержку", "Выберите тип проблемы")),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("FAQ", "FAQ")))

var toMainTheme = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Вернуться в меню", "Меню")))

var FAQKeyBoard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Регистрация", "Регистрация"),
		tgbotapi.NewInlineKeyboardButtonData("API", "API")),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Signature", "Signature"),
		tgbotapi.NewInlineKeyboardButtonData("IT вопросы", "IT")),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Заявки", "Заявки")),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Вернуться в меню", "Меню")))

var toBackOrMainTheme = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Вернуться в меню", "Меню"),
		tgbotapi.NewInlineKeyboardButtonData("Назад", "Назад")))