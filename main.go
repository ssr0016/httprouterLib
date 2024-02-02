package main

import (
	"log"

	"github.com/ssr0016/librarySystem/api"
	db "github.com/ssr0016/librarySystem/model"
)

func main() {

	store, err := db.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":3000", store)
	server.Run()
}
