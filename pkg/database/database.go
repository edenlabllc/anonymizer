package database

import (
	"ehealth-migration/pkg/gormlogger"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"reflect"
)

type Database struct {
	DB *gorm.DB
}

func (d *Database) Migrate(tableName string, m interface{}) {
	logEntry := fmt.Sprintf("Auto Migrating %s...", reflect.TypeOf(m))
	// Migrate the schema
	db := d.DB.Table(tableName).AutoMigrate(m)
	if db != nil && db.Error != nil {
		//We have an error
		log.Fatal().Msg(fmt.Sprintf("%s %s with error %s", logEntry, "Failed", db.Error))
	}
	log.Info().Msg(fmt.Sprintf("%s %s", logEntry, "Success"))
}

//GetDbClient initializing new database client
func GetDbClient(name, user, pass, host, port string) (*Database, error) {
	db, err := gorm.Open("postgres", fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, pass, name, host, port))

	db.SetLogger(&gormlogger.GormLogger{})

	db.LogMode(true)
	if err != nil {
		return nil, err
	}

	return &Database{
		db,
	}, nil
}
