package main

import (
	"database/sql"

	msg "personaerpgcompanion/pkg/models/botmsg"
	db "personaerpgcompanion/pkg/models/mysql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ParseUserMessage(botmsg *tgbotapi.Message, botmsgReply *tgbotapi.MessageConfig, dbConnect *sql.DB) {

	botmsgBuffer := tgbotapi.NewMessage(botmsg.Chat.ID, "")
	//botmsgReply.ChatID = botmsg.Chat.ID

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
			break
		case "a":
			messageToReply = db.IdentifyArmor(botmsg.CommandArguments(), dbConnect)
			break
		case "t":
			msg.TestMessage(&botmsgBuffer)
			messageToReply = "test message"
			break
		default:
			messageToReply = msg.CommandNotFoundMessage()
			break
		}
	}

	botmsgBuffer.Text = messageToReply
	*botmsgReply = botmsgBuffer
}
