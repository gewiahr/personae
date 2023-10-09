package db

import (
	"database/sql"
	"fmt"

	crypt "personaerpgcompanion/pkg"
	msg "personaerpgcompanion/pkg/models/botmsg"

	. "personaerpgcompanion/pkg/models"

	_ "github.com/go-sql-driver/mysql"
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
	armorName = fmt.Sprintf("'%s'", armorName)
	message := ""

	queryString := fmt.Sprintf("SELECT * FROM armor WHERE %s = %s;", "name", armorName)
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

func SearchForArmor(armorName string, dbConnect *sql.DB) string {

	var armorArray []string
	var message string
	var buffer string

	dataQuery := fmt.Sprintf("SELECT name FROM armor WHERE name LIKE '%%%s%%';", armorName)
	armorData, _ := dbConnect.Query(dataQuery)
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
		for i := 0; i < armorArrayLength; i++ {
			message += fmt.Sprintf("%s\n", armorArray[i])
		}
	}

	return message
}
