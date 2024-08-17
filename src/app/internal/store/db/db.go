package db

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"goth/internal/store"
	"os"
)

func open(dsn string) (*gorm.DB, error) {

	// make the temp directory if it doesn't exist
	err := os.MkdirAll("/tmp", 0755)
	if err != nil {
		return nil, err
	}

	schemaName := os.Getenv("DB_SCHEMA_NAME")

	return gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   fmt.Sprint(schemaName + "."),
			SingularTable: false,
		},
	})
}

// TODO: make a data structure to go through each table created in store.go
func createTables(db *gorm.DB) error {
	var err error
	if !db.Migrator().HasTable(&store.User{}) {
		err = db.Migrator().CreateTable(&store.User{})
	}
	if !db.Migrator().HasTable(&store.Session{}) {
		err = db.Migrator().CreateTable(&store.Session{})
	}
	if err != nil {
		return err
	}
	return nil
}

func MustOpen(dsn string) *gorm.DB {

	if dsn == "" {
		panic(errors.New("No DSN found."))
	}

	db, err := open(dsn)
	if err != nil {
		panic(err)
	}

	err = createTables(db)
	if err != nil {
		panic(err)
	}

	return db
}
