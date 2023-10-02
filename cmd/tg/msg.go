package main

import (
	"database/sql"

	msg "personaerpgcompanion/pkg/models/botmsg"
	db "personaerpgcompanion/pkg/models/mysql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ParseUserMessage(botmsg *tgbotapi.Message, dbConnect *sql.DB) string {

	messageToReply := ""

	if !(botmsg.IsCommand()) {
		messageToReply = msg.CommandNotFoundMessage()
	} else {
		switch botmsg.Command() {
		case "start":
			messageToReply = msg.WelcomeMessage()
			break
		case "help":
			messageToReply = msg.HelpMessage()
			break
		case "w":
			messageToReply = db.IdentifyWeapon(botmsg.CommandArguments(), dbConnect)
		case "a":
			messageToReply = db.IdentifyArmor(botmsg.CommandArguments(), dbConnect)
			break
		default:
			messageToReply = msg.CommandNotFoundMessage()
			break
		}
	}

	return messageToReply
}
