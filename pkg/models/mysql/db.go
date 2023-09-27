package db

import (
	"database/sql"
	"fmt"
	crypt "personaerpgcompanion/pkg"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type weapon struct {
	name       string
	tp         string
	skill      string
	rng        string
	dmg        int
	dls        int
	hand1      string
	hand2      string
	rarity     int
	price      int
	curr       string
	qualities  string
	additional string
	source     string
	pic        string
}

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
	weaponStats := weapon{}
	weaponRow := new(sql.Row)
	weaponName = fmt.Sprintf("'%s'", weaponName)
	message := ""

	queryString := fmt.Sprintf("SELECT * FROM weapons WHERE %s = %s;", "name", weaponName)
	weaponRow = dbConnect.QueryRow(queryString)

	err := weaponRow.Scan(
		&weaponStats.name,
		&weaponStats.tp,
		&weaponStats.skill,
		&weaponStats.rng,
		&weaponStats.dmg,
		&weaponStats.dls,
		&weaponStats.hand1,
		&weaponStats.hand2,
		&weaponStats.rarity,
		&weaponStats.price,
		&weaponStats.curr,
		&weaponStats.qualities,
		&weaponStats.additional,
		&weaponStats.source,
		&weaponStats.pic)
	if err != nil {
		message = "Нет такого оружия!"
	} else {
		message = ComposeWeaponMessage(weaponStats)
	}

	return message
}

func ComposeWeaponMessage(weaponStats weapon) string {
	message := ""
	// Name group
	message = weaponStats.name

	message += "\n"
	for i := 0; i < len(weaponStats.name); i++ {
		message += "="
	}

	// Main stats
	message += "\nУрон: " + strconv.Itoa(weaponStats.dmg)
	message += "\nСмертельность: " + strconv.Itoa(weaponStats.dls)

	message += "\n"

	// Skill
	message += "\nНавык: " + weaponStats.skill
	// Grip
	if weaponStats.hand1 != "X" {
		if weaponStats.hand1 == "O" {
			message += "\nОдноручное"
		} else {
			message += "\nВ одной руке: " + weaponStats.hand1
		}
	}
	if weaponStats.hand2 != "X" {
		if weaponStats.hand2 == "O" {
			message += "\nДвуручное"
		} else {
			message += "\nВ двух руках: " + weaponStats.hand2
		}
	}

	message += "\n"

	// Rarity
	message += "\nРедкость: " + strconv.Itoa(weaponStats.rarity)
	// Price
	message += "\nЦена: " + strconv.Itoa(weaponStats.price)
	switch weaponStats.curr {
	case "z":
		message += " зени (медь)"
		break
	case "b":
		message += " бу (серебро)"
		break
	case "k":
		message += " коку (золото)"
		break
	}

	message += "\n"

	// Qualities
	if weaponStats.qualities != "-" {
		message += "\nСвойства: " + weaponStats.qualities
	}
	if weaponStats.additional != "-" {
		message += "\nДополнительно: " + weaponStats.additional
	}

	// Picture
	if len(weaponStats.pic) > 0 {
		message += "\n" + weaponStats.pic
	}

	return message
}
