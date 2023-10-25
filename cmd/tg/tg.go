package main

import (
	crypt "personaerpgcompanion/pkg"
	db "personaerpgcompanion/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	exitState := 0

	// Bot initialization
	bot, err := tba.NewBotAPI(crypt.TGKey())
	LogPanic(err, "Bot initialization")

	// Bot listen to updates
	u := tba.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	// Database connection
	db.Connect, err = db.OpenDB(db.Name)
	LogPanic(err, "Database connection")
	defer db.Connect.Close()

	// Listen for updates
	for update := range updates {
		if update.Message != nil {
			LogMessage(update.Message)

			botReply := new(tba.MessageConfig)
			exitState = ParseUserMessage(update.Message, botReply, db.Connect)

			botReply.BaseChat.ChatID = update.Message.From.ID
			bot.Send(botReply)

		} else if update.CallbackQuery != nil {
			LogCallback(update.CallbackQuery)

			botReply := new(tba.MessageConfig)
			ParseCallback(update.CallbackQuery, botReply, db.Connect)

			botReply.BaseChat.ChatID = update.CallbackQuery.From.ID
			bot.Send(botReply)
		}

		if exitState == 1 {
			break
		}
	}

	LogExit(exitState)
}
