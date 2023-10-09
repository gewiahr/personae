package main

import (
	"database/sql"
	"fmt"

	crypt "personaerpgcompanion/pkg"
	msg "personaerpgcompanion/pkg/models/botmsg"
	db "personaerpgcompanion/pkg/models/mysql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ParseUserMessage(botmsg *tgbotapi.Message, botmsgReply *tgbotapi.MessageConfig, dbConnect *sql.DB) int {

	status := 0
	botmsgBuffer := tgbotapi.NewMessage(botmsg.Chat.ID, "")
	//botmsgReply.ChatID = botmsg.Chat.ID

	messageToReply := ""

	if !(botmsg.IsCommand()) {
		messageToReply = msg.CommandNotFoundMessage()
	} else {
		switch botmsg.Command() {
		case "start":
			messageToReply = msg.WelcomeMessage()
		case "help":
			messageToReply = msg.HelpMessage()
		case "w":
			messageToReply = db.IdentifyWeapon(botmsg.CommandArguments(), dbConnect)
		case "a":
			messageToReply = db.SearchForArmor(botmsg.CommandArguments(), dbConnect)
			//messageToReply = db.IdentifyArmor(botmsg.CommandArguments(), dbConnect)
		case "t":
			msg.TestMessage(&botmsgBuffer)
			messageToReply = "test message"
		case "q":
			if botmsg.From.ID == crypt.ManageID() {
				status = 1
				messageToReply = fmt.Sprintf("Bot terminated with status [%d]", status)
			}
		default:
			messageToReply = msg.CommandNotFoundMessage()
		}
	}

	botmsgBuffer.Text = messageToReply
	*botmsgReply = botmsgBuffer

	return status
}
