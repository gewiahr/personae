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

var Name = crypt.DBName("test")
var Settings = crypt.DBSettings("default")
var Connect *sql.DB

func OpenDB(dbName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", Settings+dbName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func IdentifyWeapon(weaponName string, dbConnect *sql.DB) tgbotapi.MessageConfig {

	stats := Weapon{}
	weaponRow := new(sql.Row)
	weaponName = strings.Replace(weaponName, "'", "''", -1)
	botmsg := tgbotapi.NewMessage(0, "")

	queryString := fmt.Sprintf("SELECT * FROM weapon WHERE %s = '%s';", "name", weaponName)
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
		botmsg.Text = "Нет такого оружия!"
	} else {
		botmsg.Text = msg.ComposeWeaponMessage(stats)
		if stats.Qualities != "-" {
			botmsg.BaseChat.ReplyMarkup = AppendQualities(stats.Qualities)
		}
	}

	return botmsg
}

func IdentifyArmor(armorName string, dbConnect *sql.DB) tgbotapi.MessageConfig {

	stats := Armor{}
	armorRow := new(sql.Row)
	armorName = strings.Replace(armorName, "'", "''", -1)
	botmsg := tgbotapi.NewMessage(0, "")

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
		botmsg.Text = "Нет такой одежды или брони!"
	} else {
		botmsg.Text = msg.ComposeArmorMessage(stats)
		if stats.Qualities != "-" {
			botmsg.BaseChat.ReplyMarkup = AppendQualities(stats.Qualities)
		}
	}

	return botmsg
}

func IdentifyQuality(qualityName string, dbConnect *sql.DB) string {

	stats := Quality{}
	qualityRow := new(sql.Row)
	qualityName = strings.Replace(qualityName, "'", "''", -1)
	message := ""

	queryString := fmt.Sprintf("SELECT * FROM quality WHERE %s = '%s';", "name", qualityName)
	qualityRow = dbConnect.QueryRow(queryString)

	err := qualityRow.Scan(
		&stats.Name,
		&stats.General,
		&stats.Weapon,
		&stats.Armor,
		&stats.Source)
	if err != nil {
		message = "Добавление свойств в процессе!"
	} else {
		message = msg.ComposeQualityMessage(stats)
	}

	return message
}

func IdentifyEntity(entityType string, entityName string) (interface{}, error) { //, dbConnect *sql.DB

	var err error
	stats := new(interface{})

	entityRow := new(sql.Row)
	entityName = strings.Replace(entityName, "'", "''", -1)

	queryString := fmt.Sprintf("SELECT * FROM "+entityType+" WHERE %s = '%s';", "name", entityName)
	entityRow = Connect.QueryRow(queryString)

	switch entityType {
	case "weapon":
		*stats, err = ScanWeaponRow(entityRow)
	case "armor":
		*stats, err = ScanArmorRow(entityRow)
	case "quality":
		*stats, err = ScanQualityRow(entityRow)
	}

	// if err != nil {
	// 	message = msg.ErrorIdentify(entityType)
	// } else {
	// 	switch entityType {
	// 	case "weapon":
	// 		*stats, err = ScanWeaponRow(entityRow)
	// 	case "armor":
	// 		*stats, err = ScanArmorRow(entityRow)
	// 	case "quality":
	// 		*stats, err = ScanQualityRow(entityRow)
	// 	}
	// 	message = msg.ComposeWeaponMessage(stats)
	// }

	return stats, err
}

func ScanWeaponRow(weaponRow *sql.Row) (Weapon, error) {
	var stats Weapon

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

	return stats, err
}

func ScanArmorRow(armorRow *sql.Row) (Armor, error) {
	var stats Armor

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

	return stats, err
}

func ScanQualityRow(qualityRow *sql.Row) (Quality, error) {
	var stats Quality

	err := qualityRow.Scan(
		&stats.Name,
		&stats.General,
		&stats.Weapon,
		&stats.Armor,
		&stats.Source)

	return stats, err
}

func SearchForWeapon(msgConf *tgbotapi.MessageConfig, weaponName string, dbConnect *sql.DB) string {

	var weaponArray []string
	var message string
	var buffer string
	weaponName = strings.Replace(weaponName, "'", "''", -1)

	queryString := fmt.Sprintf("SELECT name FROM weapon WHERE name LIKE '%%%s%%';", weaponName)
	weaponData, _ := dbConnect.Query(queryString)
	for weaponData.Next() {
		_ = weaponData.Scan(&buffer)
		weaponArray = append(weaponArray, buffer)
	}

	weaponArrayLength := len(weaponArray)
	if weaponArrayLength == 0 {
		message = "Нет такого оружия!"
	} else if weaponArrayLength == 1 {
		*msgConf = IdentifyWeapon(weaponArray[0], dbConnect)
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
		*msgConf = IdentifyArmor(armorArray[0], dbConnect)
	} else if armorArrayLength > 1 {
		message = ShowArmorList(armorArray, msgConf)
	}

	return message
}

func SearchForQualities(msgConf *tgbotapi.MessageConfig, qualityName string, dbConnect *sql.DB) string {

	var qualityArray []string
	var message string
	var buffer string
	qualityName = strings.Replace(qualityName, "'", "''", -1)

	queryString := fmt.Sprintf("SELECT name FROM quality WHERE name LIKE '%%%s%%';", qualityName)
	armorData, _ := dbConnect.Query(queryString)
	for armorData.Next() {
		_ = armorData.Scan(&buffer)
		qualityArray = append(qualityArray, buffer)
	}

	qualityArrayLength := len(qualityArray)
	if qualityArrayLength == 0 {
		message = "Нет такого свойства!"
	} else if qualityArrayLength == 1 {
		message = IdentifyQuality(qualityArray[0], dbConnect)
	} else if qualityArrayLength > 1 {
		message = ShowQualityList(qualityArray, msgConf)
	}

	return message

}

func ShowWeaponMenu(msgConf *tgbotapi.MessageConfig, dbConnect *sql.DB) string {

	var buffer string
	var odd = true
	inlineMenu := tgbotapi.NewInlineKeyboardMarkup()
	inlineButtonRow := tgbotapi.NewInlineKeyboardRow()

	queryString := fmt.Sprintf("SELECT DISTINCT tp FROM weapon")
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

	queryString := fmt.Sprintf("SELECT name FROM weapon WHERE tp = '%s'", tp)
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

func ShowQualityList(qualityArray []string, msgConf *tgbotapi.MessageConfig) string {

	inlineMenu := tgbotapi.NewInlineKeyboardMarkup()

	for i := 0; i < len(qualityArray); i++ {
		inlineMenu.InlineKeyboard = append(inlineMenu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(qualityArray[i], "cb_qualitysearch_"+qualityArray[i])))
	}

	msgConf.BaseChat.ReplyMarkup = inlineMenu
	return "Выберите свойство из списка ниже:"

}

func AppendQualities(qualities string) tgbotapi.InlineKeyboardMarkup {

	var odd = true
	var buffer string
	inlineMenu := tgbotapi.NewInlineKeyboardMarkup()
	inlineButtonRow := tgbotapi.NewInlineKeyboardRow()

	qualityArray := strings.Split(qualities, ", ")
	for i := 0; len(qualityArray) > i; i++ {
		buffer, _, _ = strings.Cut(qualityArray[i], ",")
		qualityArray[i] = buffer
	}

	for i := 0; i < len(qualityArray); i++ {
		if odd {
			inlineButtonRow = tgbotapi.NewInlineKeyboardRow()
			inlineButtonRow = append(inlineButtonRow, tgbotapi.NewInlineKeyboardButtonData(qualityArray[i], "cb_qualitysearch_"+qualityArray[i]))
		} else {
			inlineButtonRow = append(inlineButtonRow, tgbotapi.NewInlineKeyboardButtonData(qualityArray[i], "cb_qualitysearch_"+qualityArray[i]))
			inlineMenu.InlineKeyboard = append(inlineMenu.InlineKeyboard, inlineButtonRow)
		}
		odd = !odd
	}
	if !odd {
		inlineMenu.InlineKeyboard = append(inlineMenu.InlineKeyboard, inlineButtonRow)
	}

	return inlineMenu

}
