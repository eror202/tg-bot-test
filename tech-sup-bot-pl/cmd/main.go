package main

//Bot для клиентов Antarex

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"strings"
	tg_bot_for_ts "tg-bot-for-ts/faq"
	"tg-bot-for-ts/repository"
	"tg-bot-for-ts/service"
	"tg-bot-for-ts/tgUtil"
)

func init() {

}

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("Ошибка в конфиг файле : %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("env файл не найден или нет параметров в нем: %s", err.Error())
	}
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("БД не инициализирована : %s", err.Error())
	}
	repo := repository.NewRepository(db)
	_ = service.NewService(repo)
	/*handlers := handler.NewHandler(services)

	srv := new(tg_bot_for_ts.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("Сервер не запустился, ошибка: %s", err.Error())
	}*/

	bot, err := tgbotapi.NewBotAPI(viper.GetString("tg.bot-token"))

	if err != nil {
		logrus.Panic(err)
	}
	bot.Debug = true

	logrus.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		if update.Message != nil {
			switch update.Message.Command() {
			case "start":
				reply := "Добро пожаловать!)" +
					"\nЯ - бот-ассисент, который поможет вам при интеграции нажей платежной системы." +
					"\nНачните работу с раздела \"FAQ\", чтобы найти часто задаваемые вопросы во время интеграции." +
					"\n\nВ разделе \"Генерация подписи\" вы можете сгенерировать signature для работы с нашем API" +
					"\n\nМы ценим каждого нашего клиента. Наша команда всегда рада вашим идеям и предложениям по улучшению нашего сервиса." +
					"Оставив фидбэк в разделе \"Обратная связь\" вы помогаете стать нам лучше. Спасибо за доверие! Легкой интеграции, коллеги!"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				msg.ReplyMarkup = StartKeyBoard

				tgUtil.SendBotMessage(msg, bot)

			case "commands":
				replay := "/main - вернуться в главное меню"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, replay)
				tgUtil.SendBotMessage(msg, bot)

			case "main":
				replay := "Вы в главном меню, выберете на экране интересующий вас раздел \n"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, replay)
				msg.ReplyMarkup = StartKeyBoard
				tgUtil.SendBotMessage(msg, bot)

			default:
				replay := "Сообщение не распознано, для вызова списка команд нажмите: \n" +
					"/commands"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, replay)
				tgUtil.SendBotMessage(msg, bot)
			}
		} else if update.CallbackQuery != nil {

			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				logrus.Error(err)
			}
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			switch update.CallbackQuery.Data {

			case "Выберите тип проблемы":
				reply := "Выберите тип проблемы, с которой вы хотите обратиться"
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
				msg.ReplyMarkup = KeyBoard
				tgUtil.SendBotMessage(msg, bot)

			case "Генерация подписи":
				reply := "В данном разделе вы можете провести генерация сигнатуры для отправки запросов на сервер!)" +
					"\n\nПодробно о сигнатуре вы можете узнать в разделе \"FAQ -> Signature\"" +
					"\n\nВ поле ввода сообщений передайте входные данные в формате:" +
					"\n\nApi-Key{BodyRequest},Secret-Key" +
					"\n\nОбращаем ваше внимание, данные передаются без знаков разделителей между значениями. RequestBody и Secret-Key разделяет запятая без пробела."
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
				msg.ReplyMarkup = toMainTheme

				tgUtil.SendBotMessage(msg, bot)
				for newUpdate := range updates {
					if newUpdate.CallbackQuery != nil {

						callback := tgbotapi.NewCallback(newUpdate.CallbackQuery.ID, newUpdate.CallbackQuery.Data)
						if _, err := bot.Request(callback); err != nil {
							logrus.Error(err)
						}
						switch newUpdate.CallbackQuery.Data {
						case "Меню":
							logrus.Printf(newUpdate.CallbackQuery.Data)
							replay := "Вы в главном меню, выберете на экране интересующий вас вопрос \n" +
								"\n /commands - список команд"
							msg := tgbotapi.NewMessage(newUpdate.CallbackQuery.Message.Chat.ID, replay)
							msg.ReplyMarkup = StartKeyBoard
							tgUtil.SendBotMessage(msg, bot)
						}
					} else if newUpdate.Message == nil {
						continue
					} else {
						txt := newUpdate.Message.Text
						log.Printf(txt)
						msg := tgbotapi.NewMessage(newUpdate.Message.Chat.ID, "Ваша подпись: \n \n"+CreateSignature(txt, newUpdate, repo))
						msg.ReplyMarkup = toMainTheme
						tgUtil.SendBotMessage(msg, bot)
						break
					}
					break
				}

			case "Обратная связь":
				reply := "Коллеги, нет предела совершенству! \n\nДанный раздел предназначен для обратной связи от наших дорогих клиентов!" +
					"\n\nПожалуйста, опишите одним сообщением ваши пожелания по улучшению нашего сервиса."
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
				msg.ReplyMarkup = toMainTheme
				tgUtil.SendBotMessage(msg, bot)

				for newUpdate := range updates {
					if newUpdate.CallbackQuery != nil {

						callback := tgbotapi.NewCallback(newUpdate.CallbackQuery.ID, newUpdate.CallbackQuery.Data)
						if _, err := bot.Request(callback); err != nil {
							logrus.Error(err)
						}
						switch newUpdate.CallbackQuery.Data {
						case "Меню":
							logrus.Printf(newUpdate.CallbackQuery.Data)
							replay := "Вы в главном меню, выберете на экране интересующий вас вопрос \n" +
								"\n /commands - список команд"
							msg := tgbotapi.NewMessage(newUpdate.CallbackQuery.Message.Chat.ID, replay)
							msg.ReplyMarkup = StartKeyBoard
							tgUtil.SendBotMessage(msg, bot)
						}
					} else if newUpdate.Message == nil {
						continue
					} else {
						txt := "ХР\n" + newUpdate.Message.Chat.FirstName + " " + newUpdate.Message.Chat.LastName + "\n" + "@" + newUpdate.Message.Chat.UserName + "\n" + newUpdate.Message.Text
						replay := "Благодарим вас за обратную связь!" +
							"\n\nВместе с вами мы делаем сервис еще лучше и удобнее для вас, наших дорогих клиентов!)"
						msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, replay)
						msg.ReplyMarkup = toMainTheme
						tgUtil.SendBotMessage(msg, bot)
						msg = tgbotapi.NewMessageToChannel("1661385575", txt)
						tgUtil.SendBotMessage(msg, bot)
					}
					break
				}
			case "FAQ":
				replay := "Данный раздел сориентирует вас по часто задаваемым вопросам во время интеграции нашей платежной системы. " +
					"\nВыберите интересующий вас раздел:"
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, replay)
				msg.ReplyMarkup = FAQKeyBoard
				tgUtil.SendBotMessage(msg, bot)

			case "Регистрация":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, tg_bot_for_ts.Registration)
				msg.ReplyMarkup = FAQKeyBoard
				tgUtil.SendBotMessage(msg, bot)

			case "Signature":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, tg_bot_for_ts.SignatureQuestion)
				msg.ReplyMarkup = FAQKeyBoard
				tgUtil.SendBotMessage(msg, bot)

			case "API":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, tg_bot_for_ts.APIQuestion)
				msg.ReplyMarkup = FAQKeyBoard
				tgUtil.SendBotMessage(msg, bot)

			case "IT":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, tg_bot_for_ts.ITQuestion)
				msg.ReplyMarkup = FAQKeyBoard
				tgUtil.SendBotMessage(msg, bot)

			case "Заявки":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, tg_bot_for_ts.RequestQuestion)
				msg.ReplyMarkup = FAQKeyBoard
				tgUtil.SendBotMessage(msg, bot)

			case "Меню":
				replay := "Вы в главном меню, выберете на экране интересующий вас раздел \n" +
					"\n /commands - список команд"
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, replay)
				msg.ReplyMarkup = StartKeyBoard
				tgUtil.SendBotMessage(msg, bot)

			default:
				replay := "Некорректный запрос"
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, replay)
				tgUtil.SendBotMessage(msg, bot)
			}
		} else {
			reply := "Нераспознанная команда, для вызовка списка команд нажмите: \n" +
				"/commands"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			tgUtil.SendBotMessage(msg, bot)
		}
	}
}

func CreateSignature(inputData string, newUpdate tgbotapi.Update, repo *repository.Repository) string {
	if strings.Contains(inputData, ",") {
		str := strings.Split(inputData, ",")
		message, key := str[0], str[1]
		signature := hmac.New(sha512.New, []byte(key))
		signature.Write([]byte(message))
		return hex.EncodeToString(signature.Sum(nil))
	} else {
		return "Ошибка в заполнении входных данных"
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
