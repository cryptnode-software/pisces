package gorm

import (
	"errors"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	database *gorm.DB = nil
)

func NewDatabase(dsn string) (db *gorm.DB, err error) {

	if db, err = gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: false,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		DSN:                       dsn,
		DefaultStringSize:         256,
	})); err != nil {
		return nil, err
	}
	database = db
	return db, nil

}

func GetDatabase() (*gorm.DB, error) {
	if database == nil {
		return nil, errors.New("no database was found in memory please initilize one first")
	}
	return database, nil
}
