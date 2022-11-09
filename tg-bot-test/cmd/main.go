package main

//Bot для клиентов Antarex

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
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
		Password: os.Getenv("DB_PASSWORD"),
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

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))

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
				reply := "Добрый день, я бот - ассисент, который хочет вам помочь, " +
					"выберете на экране интересующий вас вопрос. \n У бота есть команды, " +
					"нажмите ниже на команду commands для выводы списка: \n /commands - список команд"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				msg.ReplyMarkup = StartKeyBoard

				tgUtil.SendBotMessage(msg, bot)

			case "commands":
				replay := "/main - вернуться на стартовую локацию \n" +
					"/FAQ - частые вопросы и ответы"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, replay)
				tgUtil.SendBotMessage(msg, bot)

			case "main":
				replay := "Вы в главном меню, выберете на экране интересующий вас вопрос \n" +
					"\n /commands - список команд"
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
				logrus.Panic(err)
			}
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
			switch update.CallbackQuery.Data {

			case "Выберите тип проблемы":
				reply := "Выберите тип проблемы, с которой вы хотите обратиться"
				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
				msg.ReplyMarkup = KeyBoard
				tgUtil.SendBotMessage(msg, bot)

			case "Генерация подписи":
				reply := "Отправьте message и secret-key, которые разделены запятой и без пробелов между запятой, и  подпись сгенерируется"
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
				msg.ReplyMarkup = toMainTheme

				tgUtil.SendBotMessage(msg, bot)
				for newUpdate := range updates {
					if newUpdate.CallbackQuery != nil {

						callback := tgbotapi.NewCallback(newUpdate.CallbackQuery.ID, newUpdate.CallbackQuery.Data)
						if _, err := bot.Request(callback); err != nil {
							logrus.Panic(err)
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
						msg := tgbotapi.NewMessage(newUpdate.Message.Chat.ID, "Ваша подпись: \n \n"+CreateSignature(txt, newUpdate, repo)+
							"\n\nЕсли у вас остались вопросы, то перейдите в главное меню")
						msg.ReplyMarkup = toMainTheme
						tgUtil.SendBotMessage(msg, bot)
						break
					}
					break
				}

			case "Интеграция":
				reply := "Вы в разделе \"Интеграция\".\nОпишите вашу проблему"
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
				msg.ReplyMarkup = toMainTheme
				tgUtil.SendBotMessage(msg, bot)

				for newUpdate := range updates {
					if newUpdate.CallbackQuery != nil {

						callback := tgbotapi.NewCallback(newUpdate.CallbackQuery.ID, newUpdate.CallbackQuery.Data)
						if _, err := bot.Request(callback); err != nil {
							logrus.Panic(err)
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
					} else {
						txt := newUpdate.Message.Chat.FirstName + " " + newUpdate.Message.Chat.LastName + "\n" + "@" + newUpdate.Message.Chat.UserName + "\nВопрос:\"Интеграция\"" + "\n" + newUpdate.Message.Text
						replay := "Ваш запрос зарегестрирован и ему присвоен uuid:"
						msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, replay+"\n"+tgUtil.SaveIntegrationMessage(newUpdate, repo))
						msg.ReplyMarkup = toMainTheme
						tgUtil.SendBotMessage(msg, bot)
						msg.ReplyMarkup = toMainTheme
						msg = tgbotapi.NewMessageToChannel("1661385575", txt)
						tgUtil.SendBotMessage(msg, bot)
					}
					break
				}
			case "Тесты":
				reply := "Вы в разделе \"Тесты\".\nОпишите вашу проблему"
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
				msg.ReplyMarkup = toMainTheme
				tgUtil.SendBotMessage(msg, bot)

				for newUpdate := range updates {

					if newUpdate.CallbackQuery != nil {

						callback := tgbotapi.NewCallback(newUpdate.CallbackQuery.ID, newUpdate.CallbackQuery.Data)
						if _, err := bot.Request(callback); err != nil {
							logrus.Panic(err)
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
					} else {
						txt := newUpdate.Message.Chat.FirstName + " " + newUpdate.Message.Chat.LastName + "\n" + "@" + newUpdate.Message.Chat.UserName + "\nВопрос:\"Тесты\"" + "\n" + newUpdate.Message.Text
						replay := "Ваш запрос зарегестрирован и ему присвоен uuid:"
						msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, replay+"\n"+tgUtil.SaveTestMessage(newUpdate, repo))
						msg.ReplyMarkup = toMainTheme
						tgUtil.SendBotMessage(msg, bot)
						msg = tgbotapi.NewMessageToChannel("1661385575", txt)
						tgUtil.SendBotMessage(msg, bot)
					}
					break
				}
			case "Трафик":
				reply := "Вы в разделе \"Трафик\".\nОпишите вашу проблему"
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
				msg.ReplyMarkup = toMainTheme
				tgUtil.SendBotMessage(msg, bot)

				for newUpdate := range updates {

					if newUpdate.CallbackQuery != nil {

						callback := tgbotapi.NewCallback(newUpdate.CallbackQuery.ID, newUpdate.CallbackQuery.Data)
						if _, err := bot.Request(callback); err != nil {
							logrus.Panic(err)
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
					} else {
						txt := newUpdate.Message.Chat.FirstName + " " + newUpdate.Message.Chat.LastName + "\n" + "@" + newUpdate.Message.Chat.UserName + "\nВопрос:\"Трафик\"" + "\n" + newUpdate.Message.Text
						replay := "Ваш запрос зарегестрирован и ему присвоен uuid:"
						msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, replay+"\n"+tgUtil.SaveTrafficMessage(newUpdate, repo))
						msg.ReplyMarkup = toMainTheme
						tgUtil.SendBotMessage(msg, bot)
						msg = tgbotapi.NewMessageToChannel("1661385575", txt)
						tgUtil.SendBotMessage(msg, bot)
					}
					break
				}
			case "Другое":
				reply := "Вы в разделе \"Другое\" \nОпишите вашу проблему"
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, reply)
				msg.ReplyMarkup = toMainTheme

				tgUtil.SendBotMessage(msg, bot)

				for newUpdate := range updates {

					if newUpdate.CallbackQuery != nil {

						callback := tgbotapi.NewCallback(newUpdate.CallbackQuery.ID, newUpdate.CallbackQuery.Data)
						if _, err := bot.Request(callback); err != nil {
							logrus.Panic(err)
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
					} else {
						txt := newUpdate.Message.Chat.FirstName + " " + newUpdate.Message.Chat.LastName + "\n" + "@" + newUpdate.Message.Chat.UserName + "\nВопрос:\"Другое\"" + "\n" + newUpdate.Message.Text
						replay := "Ваш запрос зарегестрирован и ему присвоен uuid:"
						msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, replay+"\n"+tgUtil.SaveOtherMessage(newUpdate, repo))
						msg.ReplyMarkup = toMainTheme
						tgUtil.SendBotMessage(msg, bot)
						msg = tgbotapi.NewMessageToChannel("1661385575", txt)
						tgUtil.SendBotMessage(msg, bot)
					}
					break
				}

			case "FAQ":
				replay := "В данном разделе находятся основные вопросы и ответы к ним!" +
					"\nВыберите интересующий Вас раздел"
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, replay)
				msg.ReplyMarkup = FAQKeyBoard
				tgUtil.SendBotMessage(msg, bot)

			case "Регистрация":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, tg_bot_for_ts.Registration)
				msg.ReplyMarkup = toBackOrMainTheme
				tgUtil.SendBotMessage(msg, bot)

			case "Signature":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, tg_bot_for_ts.SignatureQuestion)
				msg.ReplyMarkup = toBackOrMainTheme
				tgUtil.SendBotMessage(msg, bot)

			case "API":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, tg_bot_for_ts.APIQuestion)
				msg.ReplyMarkup = toBackOrMainTheme
				tgUtil.SendBotMessage(msg, bot)

			case "IT":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, tg_bot_for_ts.ITQuestion)
				msg.ReplyMarkup = toBackOrMainTheme
				tgUtil.SendBotMessage(msg, bot)

			case "Заявки":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, tg_bot_for_ts.RequestQuestion)
				msg.ReplyMarkup = toBackOrMainTheme
				tgUtil.SendBotMessage(msg, bot)

			case "Меню":
				replay := "Вы в главном меню, выберете на экране интересующий вас вопрос \n" +
					"\n /commands - список команд"
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, replay)
				msg.ReplyMarkup = StartKeyBoard
				tgUtil.SendBotMessage(msg, bot)

			case "Назад":
				replay := "В данном разделе находятся основные вопросы и ответы к ним!" +
					"\nВыберите интересующий Вас раздел"
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, replay)
				msg.ReplyMarkup = FAQKeyBoard
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
		//tgUtil.CreateSignature(newUpdate, repo, message, key)
		mac := hmac.New(sha512.New, []byte(key))
		mac.Write([]byte(message))
		hash := sha512.New()
		hash.Write(mac.Sum(nil))
		return base64.StdEncoding.EncodeToString(hash.Sum(nil))
	} else {
		return "Ошибка в заполнении входных данных"
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
