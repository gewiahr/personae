package main

import (
	"database/sql"
	"fmt"
	"strings"

	crypt "personaerpgcompanion/pkg"
	. "personaerpgcompanion/pkg/models"
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
		arguments := botmsg.CommandArguments()
		switch botmsg.Command() {
		case "start":
			messageToReply = msg.WelcomeMessage()
		case "help":
			messageToReply = msg.HelpMessage()
		case "w":
			if arguments == "" {
				messageToReply = db.ShowWeaponMenu(&botmsgBuffer, dbConnect)
			} else {
				messageToReply = db.SearchForWeapon(&botmsgBuffer, arguments, dbConnect)
			}
		case "a":
			if arguments == "" {
				messageToReply = db.ShowArmorMenu(&botmsgBuffer, dbConnect)
			} else {
				messageToReply = db.SearchForArmor(&botmsgBuffer, arguments, dbConnect)
			}
		case "q":
			messageToReply = db.SearchForQualities(&botmsgBuffer, arguments, dbConnect)
		case "t":
			msg.TestMessage(&botmsgBuffer)
			messageToReply = "test message"
		case "e":
			intf := new(interface{})
			//a := new(Armor)
			*intf, _ = db.IdentifyEntity(WeaponStr, arguments)
			//botmsgBuffer
		case "z":
			if botmsg.From.ID == crypt.ManageID() {
				status = 1
				messageToReply = fmt.Sprintf("Bot terminated with status [%d]", status)
			}
		default:
			messageToReply = msg.CommandNotFoundMessage()
		}
	}

	if messageToReply != "" {
		botmsgBuffer.Text = messageToReply
	}
	//botmsgBuffer.BaseChat.ChatID = botmsg.Chat.ID
	*botmsgReply = botmsgBuffer

	return status
}

func ParseCallback(callback *tgbotapi.CallbackQuery, botmsgReply *tgbotapi.MessageConfig, dbConnect *sql.DB) {

	//botmsgBuffer := tgbotapi.NewMessage(callback.Message.Chat.ID, "")
	cbData := callback.Data

	if strings.Contains(cbData, "cb_armormenu_") {
		tp := strings.TrimLeft(cbData, "cb_armormenu_")
		botmsgReply.Text = db.ShowArmorCategory(tp, botmsgReply, dbConnect)
	} else if strings.Contains(cbData, "cb_armorsearch_") {
		armor := strings.TrimLeft(cbData, "cb_armorsearch_")
		*botmsgReply = db.IdentifyArmor(armor, dbConnect)
		//botmsgReply.Text = db.SearchForArmor(botmsgReply, armor, dbConnect)
	} else if strings.Contains(cbData, "cb_weaponmenu_") {
		tp := strings.TrimLeft(cbData, "cb_weaponmenu_")
		botmsgReply.Text = db.ShowWeaponCategory(tp, botmsgReply, dbConnect)
	} else if strings.Contains(cbData, "cb_weaponsearch_") {
		weapon := strings.TrimLeft(cbData, "cb_weaponsearch_")
		*botmsgReply = db.IdentifyWeapon(weapon, dbConnect)
		//botmsgReply.Text = db.SearchForWeapon(botmsgReply, weapon, dbConnect)
	} else if strings.Contains(cbData, "cb_qualitysearch_") {
		quality := strings.TrimLeft(cbData, "cb_qualitysearch_")
		botmsgReply.Text = db.IdentifyQuality(quality, dbConnect)
		//botmsgReply.Text = db.SearchForWeapon(botmsgReply, weapon, dbConnect)
	}

}
