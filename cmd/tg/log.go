package main

import (
	"log"

	tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func LogPanic(err error, message string) bool {
	panicked := false
	if err != nil {
		log.Panic(err, message)
		panicked = true
	}

	return panicked
}

func LogMessage(message *tba.Message) {
	log.Printf("[%s - %d] %s", message.From.UserName, message.From.ID, message.Text)
}

func LogCallback(callback *tba.CallbackQuery) {
	log.Printf("[%s - %d] %s", callback.From.UserName, callback.From.ID, callback.Data)
}

func LogExit(exitState int) {
	log.Printf("Bot terminated with status [%d]", exitState)
}
