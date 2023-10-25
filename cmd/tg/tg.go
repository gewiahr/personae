package main

import (
	"log"

	crypt "personaerpgcompanion/pkg"
	db "personaerpgcompanion/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	state := 0

	// Bot
	bot, err := tgbotapi.NewBotAPI(crypt.TGKey())
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Database
	db.Connect, err = db.OpenDB(db.Name)
	if err != nil {
		log.Panic(err)
	}
	defer db.Connect.Close()

	// Update
	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s - %d] %s", update.Message.From.UserName, update.Message.From.ID, update.Message.Text)

			botmsgConf := new(tgbotapi.MessageConfig)
			state = ParseUserMessage(update.Message, botmsgConf, db.Connect)

			botmsgConf.BaseChat.ChatID = update.Message.From.ID
			bot.Send(botmsgConf)

			if state == 1 {
				break
			}
		} else if update.CallbackQuery != nil {
			log.Printf("[%s - %d] %s", update.CallbackQuery.From.UserName, update.CallbackQuery.From.ID, update.CallbackQuery.Data)

			botmsgConf := new(tgbotapi.MessageConfig)
			ParseCallback(update.CallbackQuery, botmsgConf, db.Connect)

			botmsgConf.BaseChat.ChatID = update.CallbackQuery.From.ID
			bot.Send(botmsgConf)
		}
	}

	log.Printf("Bot terminated with status [%d]", state)
}
