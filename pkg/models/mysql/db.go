package db

import (
	"database/sql"
	"fmt"
	"strings"

	crypt "personaerpgcompanion/pkg"

	. "personaerpgcompanion/pkg/models"

	_ "github.com/go-sql-driver/mysql"
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

func IdentifyEntity(entityName string, entityType string, dbConnect *sql.DB) Entity {

	entity := GetEntityByType(entityType)
	entityName = strings.Replace(entityName, "'", "''", -1)

	queryString := fmt.Sprintf("SELECT * FROM "+entityType+" WHERE %s = '%s';", "name", entityName)
	entityRow := dbConnect.QueryRow(queryString)

	entity, err := entity.LoadSQLData(entityRow)
	if err != nil {
		entity = nil
	}

	return entity

}

func GetMenuCategories(dbTable string, dbConnect *sql.DB) []string {

	// *** replace with a method to separate the entities with and without types
	if dbTable == "quality" {
		return nil
	}

	var buffer string
	var categories []string

	queryString := fmt.Sprintf("SELECT DISTINCT cat FROM " + dbTable)
	tpData, err := dbConnect.Query(queryString)
	if err != nil {
		categories = nil
	} else {
		for tpData.Next() {
			_ = tpData.Scan(&buffer)
			categories = append(categories, buffer)
		}
	}

	return categories

}

func SearchForEntities(entityType string, entityName string, dbConnect *sql.DB) []string {

	var entityArray []string
	var buffer string
	entityName = strings.Replace(entityName, "'", "''", -1)

	queryString := fmt.Sprintf("SELECT name FROM %s WHERE name LIKE '%%%s%%';", entityType, entityName)
	weaponData, _ := dbConnect.Query(queryString)
	for weaponData.Next() {
		_ = weaponData.Scan(&buffer)
		entityArray = append(entityArray, buffer)
	}

	return entityArray

}

func SearchForEntitiesByCategory(entityType string, entityCategory string, dbConnect *sql.DB) []string {

	var entityArray []string
	var buffer string
	entityCategory = strings.Replace(entityCategory, "'", "''", -1)

	queryString := fmt.Sprintf("SELECT name FROM %s WHERE cat = '%s';", entityType, entityCategory)
	weaponData, _ := dbConnect.Query(queryString)
	for weaponData.Next() {
		_ = weaponData.Scan(&buffer)
		entityArray = append(entityArray, buffer)
	}

	return entityArray

}
