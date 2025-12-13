package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zeltbrennt/go-api/internal/server"
	"github.com/zeltbrennt/go-api/internal/store"
)

const port = 8090

func main() {
	store := store.NewMockStore()
	// store, err := store.NewMongoStore("mongodb://root:example@mongodb:27017", "myDB", "myCollection")
	// TODO: this should fail, if mongo is not available!
	//if err != nil {
	//	log.Fatal("failed to get database connection")
	//}
	server := server.NewTaskServer(store)
	log.Println("Server startet and listening on port", port)
	err := http.ListenAndServe(fmt.Sprintf("localhost: %d", port), server.Routes())
	if err != nil {
		log.Fatal("server crashed...")
	}
}
