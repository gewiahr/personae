package main

import (
	"database/sql"
	msg "personaerpgcompanion/pkg/models/botmsg"
	db "personaerpgcompanion/pkg/models/mysql"

	. "personaerpgcompanion/pkg/models"

	tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ReplyToInfoCommand(botReply *tba.MessageConfig, cc string, arguments string, dbConnect *sql.DB) {

	message := ""
	entityType := Command[cc]

	if arguments == "" {
		categories := db.GetMenuCategories(entityType, dbConnect)
		if categories == nil {
			entities := db.SearchForEntities(entityType, arguments, dbConnect)
			botReply.BaseChat.ReplyMarkup = ConstructButtonMenu(entities, entityType, "exact")
			message = "Найдено несколько вариантов:"
		} else {
			botReply.BaseChat.ReplyMarkup = ConstructButtonMenu(categories, Command[cc], "category")
			message = "Выберите категорию из списка ниже:"
		}
	} else {
		entities := db.SearchForEntities(entityType, arguments, dbConnect)
		entitiesAmmount := len(entities)
		if entitiesAmmount == 0 {
			message = "Не найдено!"
		} else if entitiesAmmount == 1 {
			message = ReplyWithSingleEntity(botReply, arguments, entityType, dbConnect)
		} else if entitiesAmmount > 1 {
			botReply.BaseChat.ReplyMarkup = ConstructButtonMenu(entities, entityType, "exact")
			message = "Найдено несколько вариантов:"
		}
	}

	botReply.Text = message

}

func ReplyWithSingleEntity(botReply *tba.MessageConfig, entityName string, entityType string, dbConnect *sql.DB) string {

	message := ""

	entity := db.IdentifyEntity(entityName, entityType, dbConnect)
	if entity == nil {
		message = "Не найдено!"
	} else {
		message = msg.ComposeEntityMessage(entity, entityType)
		qualities := entity.GetQualities()
		if qualities != nil {
			botReply.BaseChat.ReplyMarkup = ConstructButtonMenu(qualities, "quality", "exact")
		}
	}

	return message

}

func ConstructButtonMenu(list []string, listName string, listCallback string) tba.InlineKeyboardMarkup {

	var odd = true
	inlineMenu := tba.NewInlineKeyboardMarkup()
	inlineButtonRow := tba.NewInlineKeyboardRow()

	for _, entity := range list {
		if odd {
			inlineButtonRow = tba.NewInlineKeyboardRow()
			inlineButtonRow = append(inlineButtonRow, tba.NewInlineKeyboardButtonData(entity, "cb_"+listCallback+"_"+listName+"_"+entity))
		} else {
			inlineButtonRow = append(inlineButtonRow, tba.NewInlineKeyboardButtonData(entity, "cb_"+listCallback+"_"+listName+"_"+entity))
			inlineMenu.InlineKeyboard = append(inlineMenu.InlineKeyboard, inlineButtonRow)
		}
		odd = !odd
	}
	if !odd {
		inlineMenu.InlineKeyboard = append(inlineMenu.InlineKeyboard, inlineButtonRow)
	}

	return inlineMenu

}
