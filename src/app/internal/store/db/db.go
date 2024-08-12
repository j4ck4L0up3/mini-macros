package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"goth/internal/store"
	"os"
)

func open(dsn string) (*gorm.DB, error) {

	// make the temp directory if it doesn't exist
	err := os.MkdirAll("/tmp", 0755)
	if err != nil {
		return nil, err
	}

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func MustOpen(dsn string) *gorm.DB {

	// FIXME: this will have to be changed
	if dsn == "" {
		dsn = "test.db"
	}

	db, err := open(dsn)
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&store.User{}, &store.Session{})

	if err != nil {
		panic(err)
	}

	return db
}
