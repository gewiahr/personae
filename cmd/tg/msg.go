package main

import (
	"database/sql"
	"strings"

	crypt "personaerpgcompanion/pkg"
	_ "personaerpgcompanion/pkg/models"
	msg "personaerpgcompanion/pkg/models/botmsg"
	db "personaerpgcompanion/pkg/models/mysql"

	tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ParseUserMessage(usermsg *tba.Message, botReply *tba.MessageConfig, dbConnect *sql.DB) int {

	status := 0
	botmsgBuffer := tba.NewMessage(usermsg.Chat.ID, "")
	messageToReply := ""

	if usermsg.IsCommand() {
		arguments := usermsg.CommandArguments()
		command := usermsg.Command()
		switch command {
		case "start":
			messageToReply = msg.WelcomeMessage()
		case "help":
			messageToReply = msg.HelpMessage()
		case "w":
			ReplyToInfoCommand(&botmsgBuffer, command, arguments, dbConnect) // messageToReply = ReplyToWeaponCommand(&botmsgBuffer, arguments, dbConnect)
		case "a":
			ReplyToInfoCommand(&botmsgBuffer, command, arguments, dbConnect)
		case "q":
			ReplyToInfoCommand(&botmsgBuffer, command, arguments, dbConnect) //messageToReply = db.SearchForQualities(&botmsgBuffer, arguments, dbConnect)
		case "e":
			//intf := new(interface{})
			//a := new(Armor)
			//*intf, _ = db.IdentifyEntity(Command["w"], arguments)
			//botmsgBuffer
		case "z":
			if usermsg.From.ID == crypt.ManageID() {
				status = 1
			}
		default:
			messageToReply = msg.CommandNotFoundMessage()
		}
	} else {
		messageToReply = msg.CommandNotFoundMessage()
	}

	if messageToReply != "" {
		botmsgBuffer.Text = messageToReply
	}
	//botmsgBuffer.BaseChat.ChatID = botmsg.Chat.ID
	*botReply = botmsgBuffer

	return status
}

func ParseCallback(callback *tba.CallbackQuery, botReply *tba.MessageConfig, dbConnect *sql.DB) {

	cbData := strings.Split(callback.Data, "_")
	if cbData[0] != "cb" {
		LogCallbackError(callback, "sterror")
		return
	}

	message := ""

	switch cbData[1] {
	case "exact":
		message = ReplyWithSingleEntity(botReply, cbData[3], cbData[2], dbConnect)
	case "category":
		entities := db.SearchForEntitiesByCategory(cbData[2], cbData[3], dbConnect)
		entitiesAmmount := len(entities)
		if entitiesAmmount == 0 {
			message = "Не найдено!"
		} else if entitiesAmmount == 1 {
			message = ReplyWithSingleEntity(botReply, entities[0], cbData[2], dbConnect)
		} else if entitiesAmmount > 1 {
			botReply.BaseChat.ReplyMarkup = ConstructButtonMenu(entities, cbData[2], "exact")
			message = "Найдено несколько вариантов:"
		}

	case "list":
		entities := db.SearchForEntities(cbData[2], cbData[3], dbConnect)
		entitiesAmmount := len(entities)
		if entitiesAmmount == 0 {
			message = "Не найдено!"
		} else if entitiesAmmount == 1 {
			message = ReplyWithSingleEntity(botReply, entities[0], cbData[2], dbConnect)
		} else if entitiesAmmount > 1 {
			botReply.BaseChat.ReplyMarkup = ConstructButtonMenu(entities, cbData[2], "exact")
			message = "Найдено несколько вариантов:"
		}
	}

	botReply.Text = message

}
