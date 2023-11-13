package types

import (
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Entity interface {
	LoadSQLData(row *sql.Row) (Entity, error)
	GetQualities() []string
}

func GetEntityByType(entityTypeName string) Entity {

	entityType := *new(Entity)

	switch entityTypeName {
	case "weapon":
		entityType = *new(Weapon)
	case "armor":
		entityType = *new(Armor)
	case "quality":
		entityType = *new(Quality)
	}

	return entityType

}

type Weapon struct {
	Name  string
	TP    string
	Skill string

	RNG   string
	DMG   int
	DLS   int
	Hand1 string
	Hand2 string

	Rarity int
	Price  int
	Curr   string

	Qualities  string
	Additional string

	Source string
	Pic    string
}

func (w Weapon) LoadSQLData(row *sql.Row) (Entity, error) {

	err := row.Scan(
		&w.Name,
		&w.TP,
		&w.Skill,
		&w.RNG,
		&w.DMG,
		&w.DLS,
		&w.Hand1,
		&w.Hand2,
		&w.Rarity,
		&w.Price,
		&w.Curr,
		&w.Qualities,
		&w.Additional,
		&w.Source,
		&w.Pic)

	return w, err

}

func (w Weapon) GetQualities() []string {
	if w.Qualities == "-" {
		return nil
	}

	qA := strings.Split(w.Qualities, ", ")
	qAA := [][]string{}
	for _, q := range qA {
		qs := strings.Split(q, " ")
		qAA = append(qAA, qs)
	}

	for i, q := range qAA {
		qA[i] = q[0]
	}

	return qA
}

type Armor struct {
	Name string
	TP   string

	Phys  int
	Super int

	Rarity int
	Price  int
	Curr   string

	Qualities  string
	Additional string

	Source string
	Pic    string
}

func (a Armor) LoadSQLData(row *sql.Row) (Entity, error) {

	err := row.Scan(
		&a.Name,
		&a.TP,
		&a.Phys,
		&a.Super,
		&a.Rarity,
		&a.Price,
		&a.Curr,
		&a.Qualities,
		&a.Additional,
		&a.Source,
		&a.Pic)

	return a, err

}

func (a Armor) GetQualities() []string {
	if a.Qualities == "-" {
		return nil
	}

	qA := strings.Split(a.Qualities, ", ")
	qAA := [][]string{}
	for _, q := range qA {
		qs := strings.Split(q, " ")
		qAA = append(qAA, qs)
	}

	for i, q := range qAA {
		qA[i] = q[0]
	}

	return qA
}

type Quality struct {
	Name    string
	General string

	Weapon bool
	Armor  bool

	Source string
}

func (q Quality) LoadSQLData(row *sql.Row) (Entity, error) {

	err := row.Scan(
		&q.Name,
		&q.General,
		&q.Weapon,
		&q.Armor,
		&q.Source)

	return q, err

}

func (q Quality) GetQualities() []string {
	return nil
}

var Command = map[string]string{
	"w": "weapon",
	"a": "armor",
	"q": "quality",
	"i": "items",
}
