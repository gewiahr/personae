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

			botmsgConf := new(tgbotapi.MessageConfig)
			ParseUserMessage(update.Message, botmsgConf, dbConnect)

			bot.Send(botmsgConf)
		}
	}
}
