package main

import (
	"log"

	crypt "personaerpgcompanion/pkg"
	db "personaerpgcompanion/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	// Bot
	bot, err := tgbotapi.NewBotAPI(crypt.TGKey())
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Database
	dbName := crypt.DBName("test")
	dbConnect, err := db.OpenDB(dbName)
	if err != nil {
		log.Panic(err)
	}
	defer dbConnect.Close()

	// Update
	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s - %d] %s", update.Message.From.UserName, update.Message.From.ID, update.Message.Text)
			messageToReply := ""
			//	"Команда не опознана!\n\"w НазваниеОружия\" чтобы найти оружие"

			if !(update.Message.IsCommand()) {

				messageToReply = "Команда не опознана!\nВоспользуйтесь командой /help чтобы узнать о командах бота"

			} else {
				switch update.Message.Command() {
				case "start":
					messageToReply = "Вас приветствует бот Personae — ваш компаньон в мире НРИ!"
					messageToReply += "\n\nСейчас бот находится на стадии разработки, а функционал непрерывно дорабатывается. "
					messageToReply += "Чтобы предложить свои идеи или обозначить ошибке в работе бота, свяжитесь с разработчиком @gewiahr."
					messageToReply += "\n\nДля того чтобы узнать текущий функционал введите команду /help."
					break
				case "help":
					messageToReply = "⚔️ Для поиска оружия используйте команду /w и название оружия. Пока что бот понимает только оригинальные названия оружия на английском языке"
					break
				case "w":
					messageToReply = db.IdentifyWeapon(update.Message.CommandArguments(), dbConnect)
					break
				default:
					messageToReply = "Команда не опознана!\nВоспользуйтесь командой /help чтобы узнать о командах бота"
					break
				}
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageToReply)
			msg.ParseMode = "HTML"
			//photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FilePath("aaa.png"))
			//msg. Photo
			//entity := new(tgbotapi.MessageEntity)
			//entity.Type = "bot_command"
			//entity.Length = 3
			//entity.Offset = 2
			//msg.Entities = append(msg.Entities, *entity)
			//msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}
