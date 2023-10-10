package db

import (
	"database/sql"
	"fmt"
	"strings"

	crypt "personaerpgcompanion/pkg"
	msg "personaerpgcompanion/pkg/models/botmsg"

	. "personaerpgcompanion/pkg/models"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var dbSettings = crypt.DBSettings("default")

func OpenDB(dbName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dbSettings+dbName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func IdentifyWeapon(weaponName string, dbConnect *sql.DB) string {
	stats := Weapon{}
	weaponRow := new(sql.Row)
	weaponName = fmt.Sprintf("'%s'", weaponName)
	message := ""

	queryString := fmt.Sprintf("SELECT * FROM weapons WHERE %s = %s;", "name", weaponName)
	weaponRow = dbConnect.QueryRow(queryString)

	err := weaponRow.Scan(
		&stats.Name,
		&stats.TP,
		&stats.Skill,
		&stats.RNG,
		&stats.DMG,
		&stats.DLS,
		&stats.Hand1,
		&stats.Hand2,
		&stats.Rarity,
		&stats.Price,
		&stats.Curr,
		&stats.Qualities,
		&stats.Additional,
		&stats.Source,
		&stats.Pic)
	if err != nil {
		message = "Нет такого оружия!"
	} else {
		message = msg.ComposeWeaponMessage(stats)
	}

	return message
}

func IdentifyArmor(armorName string, dbConnect *sql.DB) string {

	stats := Armor{}
	armorRow := new(sql.Row)
	armorName = strings.Replace(armorName, "'", "''", -1)
	message := ""

	queryString := fmt.Sprintf("SELECT * FROM armor WHERE %s = '%s';", "name", armorName)
	armorRow = dbConnect.QueryRow(queryString)

	err := armorRow.Scan(
		&stats.Name,
		&stats.TP,
		&stats.Phys,
		&stats.Super,
		&stats.Rarity,
		&stats.Price,
		&stats.Curr,
		&stats.Qualities,
		&stats.Additional,
		&stats.Source,
		&stats.Pic)
	if err != nil {
		message = "Нет такой одежды или брони!"
	} else {
		message = msg.ComposeArmorMessage(stats)
	}

	return message
}

func SearchForWeapon(msgConf *tgbotapi.MessageConfig, weaponName string, dbConnect *sql.DB) string {

	var weaponArray []string
	var message string
	var buffer string
	weaponName = strings.Replace(weaponName, "'", "''", -1)

	queryString := fmt.Sprintf("SELECT name FROM weapons WHERE name LIKE '%%%s%%';", weaponName)
	weaponData, _ := dbConnect.Query(queryString)
	for weaponData.Next() {
		_ = weaponData.Scan(&buffer)
		weaponArray = append(weaponArray, buffer)
	}

	weaponArrayLength := len(weaponArray)
	if weaponArrayLength == 0 {
		message = "Нет такого оружия!"
	} else if weaponArrayLength == 1 {
		message = IdentifyWeapon(weaponArray[0], dbConnect)
	} else if weaponArrayLength > 1 {
		message = ShowWeaponList(weaponArray, msgConf)
	}

	return message
}

func SearchForArmor(msgConf *tgbotapi.MessageConfig, armorName string, dbConnect *sql.DB) string {

	var armorArray []string
	var message string
	var buffer string
	armorName = strings.Replace(armorName, "'", "''", -1)

	queryString := fmt.Sprintf("SELECT name FROM armor WHERE name LIKE '%%%s%%';", armorName)
	armorData, _ := dbConnect.Query(queryString)
	for armorData.Next() {
		_ = armorData.Scan(&buffer)
		armorArray = append(armorArray, buffer)
	}

	armorArrayLength := len(armorArray)
	if armorArrayLength == 0 {
		message = "Нет такой одежды или брони!"
	} else if armorArrayLength == 1 {
		message = IdentifyArmor(armorArray[0], dbConnect)
	} else if armorArrayLength > 1 {
		message = ShowArmorList(armorArray, msgConf)
	}

	return message
}

func ShowWeaponMenu(msgConf *tgbotapi.MessageConfig, dbConnect *sql.DB) string {

	var buffer string
	var odd = true
	inlineMenu := tgbotapi.NewInlineKeyboardMarkup()
	inlineButtonRow := tgbotapi.NewInlineKeyboardRow()

	queryString := fmt.Sprintf("SELECT DISTINCT tp FROM weapons")
	tpData, _ := dbConnect.Query(queryString)
	for tpData.Next() {
		_ = tpData.Scan(&buffer)
		if odd {
			inlineButtonRow = tgbotapi.NewInlineKeyboardRow()
			inlineButtonRow = append(inlineButtonRow, tgbotapi.NewInlineKeyboardButtonData(buffer, "cb_weaponmenu_"+buffer))
		} else {
			inlineButtonRow = append(inlineButtonRow, tgbotapi.NewInlineKeyboardButtonData(buffer, "cb_weaponmenu_"+buffer))
			inlineMenu.InlineKeyboard = append(inlineMenu.InlineKeyboard, inlineButtonRow)
		}
		odd = !odd
	}
	if !odd {
		inlineMenu.InlineKeyboard = append(inlineMenu.InlineKeyboard, inlineButtonRow)
	}

	msgConf.BaseChat.ReplyMarkup = inlineMenu
	return "Выберите категорию оружия из списка ниже:"

}

func ShowArmorMenu(msgConf *tgbotapi.MessageConfig, dbConnect *sql.DB) string {

	var buffer string
	var odd = true
	inlineMenu := tgbotapi.NewInlineKeyboardMarkup()
	inlineButtonRow := tgbotapi.NewInlineKeyboardRow()

	queryString := fmt.Sprintf("SELECT DISTINCT tp FROM armor")
	tpData, _ := dbConnect.Query(queryString)
	for tpData.Next() {
		_ = tpData.Scan(&buffer)
		if odd {
			inlineButtonRow = tgbotapi.NewInlineKeyboardRow()
			inlineButtonRow = append(inlineButtonRow, tgbotapi.NewInlineKeyboardButtonData(buffer, "cb_armormenu_"+buffer))
		} else {
			inlineButtonRow = append(inlineButtonRow, tgbotapi.NewInlineKeyboardButtonData(buffer, "cb_armormenu_"+buffer))
			inlineMenu.InlineKeyboard = append(inlineMenu.InlineKeyboard, inlineButtonRow)
		}
		odd = !odd
	}
	if !odd {
		inlineMenu.InlineKeyboard = append(inlineMenu.InlineKeyboard, inlineButtonRow)
	}

	msgConf.BaseChat.ReplyMarkup = inlineMenu
	return "Выберите категорию брони из списка ниже:"

}

func ShowWeaponCategory(tp string, msgConf *tgbotapi.MessageConfig, dbConnect *sql.DB) string {

	var buffer string
	inlineMenu := tgbotapi.NewInlineKeyboardMarkup()

	queryString := fmt.Sprintf("SELECT name FROM weapons WHERE tp = '%s'", tp)
	tpData, _ := dbConnect.Query(queryString)
	for tpData.Next() {
		_ = tpData.Scan(&buffer)
		inlineMenu.InlineKeyboard = append(inlineMenu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(buffer, "cb_weaponsearch_"+buffer)))
	}

	msgConf.BaseChat.ReplyMarkup = inlineMenu
	return "Выберите оружие из списка ниже:"

}

func ShowArmorCategory(tp string, msgConf *tgbotapi.MessageConfig, dbConnect *sql.DB) string {

	var buffer string
	inlineMenu := tgbotapi.NewInlineKeyboardMarkup()

	queryString := fmt.Sprintf("SELECT name FROM armor WHERE tp = '%s'", tp)
	tpData, _ := dbConnect.Query(queryString)
	for tpData.Next() {
		_ = tpData.Scan(&buffer)
		inlineMenu.InlineKeyboard = append(inlineMenu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(buffer, "cb_armorsearch_"+buffer)))
	}

	msgConf.BaseChat.ReplyMarkup = inlineMenu
	return "Выберите броню из списка ниже:"

}

func ShowWeaponList(weaponArray []string, msgConf *tgbotapi.MessageConfig) string {

	inlineMenu := tgbotapi.NewInlineKeyboardMarkup()

	for i := 0; i < len(weaponArray); i++ {
		inlineMenu.InlineKeyboard = append(inlineMenu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(weaponArray[i], "cb_weaponsearch_"+weaponArray[i])))
	}

	msgConf.BaseChat.ReplyMarkup = inlineMenu
	return "Выберите оружие из списка ниже:"

}

func ShowArmorList(armorArray []string, msgConf *tgbotapi.MessageConfig) string {

	inlineMenu := tgbotapi.NewInlineKeyboardMarkup()

	for i := 0; i < len(armorArray); i++ {
		inlineMenu.InlineKeyboard = append(inlineMenu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(armorArray[i], "cb_armorsearch_"+armorArray[i])))
	}

	msgConf.BaseChat.ReplyMarkup = inlineMenu
	return "Выберите броню из списка ниже:"

}
