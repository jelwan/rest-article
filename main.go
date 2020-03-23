package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"rest-article/app"
	"rest-article/config"
	"rest-article/database"
	"rest-article/log"
)

var logger = log.NewLogger().WithField("package", "main")

func init() {
	logger.Infof("Starting server on port ['%d']", config.App().Server.Port.Http)
}

func main() {
	ctx := context.Background()

	db, err := database.CreateDatabase()
	if err != nil {
		logger.Fatalf("Database connection failed: %s", err.Error())
	}

	defer db.Close()

	api := app.NewApp(
		mux.NewRouter().StrictSlash(true),
		db,
		ctx)

	api.SetupRouter()

	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.App().Server.Port.Http), api.Router); err != nil {
		logger.Errorf("error starting server because: %v", err)
		os.Exit(1)
	}
}
