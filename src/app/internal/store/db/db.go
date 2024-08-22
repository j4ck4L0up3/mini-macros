package db

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"goth/internal/store"
	"os"
	"reflect"
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

func createTables(db *gorm.DB, tables *store.Tables) error {

	var err error
	values := reflect.ValueOf(tables).Elem()

	for i := 0; i < values.NumField(); i++ {
		field := values.Field(i)

		if !db.Migrator().HasTable(field.Addr().Interface()) {
			err = db.Migrator().CreateTable(field.Addr().Interface())
		}
		if err != nil {
			return err
		}
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

	tables := &store.Tables{
		User:               &store.User{},
		Admin:              &store.Admin{},
		Macro:              &store.Macro{},
		Session:            &store.Session{},
		PasswordResetToken: &store.PasswordResetToken{},
	}

	err = createTables(db, tables)
	if err != nil {
		panic(err)
	}

	return db
}
