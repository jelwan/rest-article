package database

import (
	"database/sql"
	"fmt"
	"rest-article/config"
	"rest-article/log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)

var logger = log.NewLogger().WithField("module", "database")

func CreateDatabase() (*sql.DB, error) {

	param := "parseTime=true&multiStatements=true"
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		config.App().Database.User,
		config.App().Database.Pass,
		config.App().Database.Host,
		config.App().Database.Port,
		config.App().Database.Schema,
		param)

	db, err := sql.Open(config.App().Database.Type, dataSourceName)
	if err != nil {
		logger.Errorf("error opening connection to db because: %v", err)
		return nil, err
	}

	return db, nil
}
