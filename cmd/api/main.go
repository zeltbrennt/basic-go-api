package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/zeltbrennt/go-api/internal/server"
	"github.com/zeltbrennt/go-api/internal/store"
)

const port = 8090

// TaskServer godoc
//
//	@titlte Task server
//	@verison 1
//	@description Task Server API
//	@basePath /api/v1
//	@host localhost:8090
func main() {
	store := store.NewMockStore()
	// store, err := store.NewMongoStore("mongodb://root:example@mongodb:27017", "myDB", "myCollection")
	// TODO: this should fail, if mongo is not available!
	//if err != nil {
	//	log.Fatal("failed to get database connection")
	//}
	server := server.NewTaskServer(store)
	logHandler := slog.NewJSONHandler(os.Stderr, nil) // change this to a fileWriter to keep the logs
	logger := slog.New(logHandler)

	logger.Info(fmt.Sprintf("Server startet and listening on port %d", port))
	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", port), server.Routes())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
