package database

import (
	"database/sql"
	"fmt"
	"rest-article/config"
	"rest-article/log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)

var logger = log.NewLogger().WithField("module", "database")

func CreateDatabase() (*sql.DB, error) {

	param := "parseTime=true&multiStatements=true"
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		config.App().Database.User,
		config.App().Database.Pass,
		config.App().Database.Host,
		config.App().Database.Schema,
		param)

	db, err := sql.Open(config.App().Database.Type, dataSourceName)
	if err != nil {
		logger.Errorf("error opening connection to db because: %v", err)
		return nil, err
	}

	//if err := migrateDatabase(db); err != nil {
	//	return db, err
	//}

	return db, nil
}

func migrateDatabase(db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	//dir, err := os.Getwd()
	//if err != nil {
	//	log.Fatal(err)
	//}

	migration, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", config.App().Database.MigrationPath),
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	migration.Log.Printf("Applying database migrations")
	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	version, _, err := migration.Version()
	if err != nil {
		return err
	}

	migration.Log.Printf("Active database version: %d", version)

	return nil
}
