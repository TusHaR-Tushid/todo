package main

import (
	"fmt"
	"go-todo/database"
	"go-todo/server"
	"log"
)

func main() {
	err := database.ConnectAndMigrate("localhost", "5435", "todo", "postgres", "1234", database.SSLModeDisable)
	if err != nil {
		log.Printf("%v", err)
		panic(err)
	}
	fmt.Println("connected")
	srv := server.SetupRoutes()
	err = srv.Run(":8080")
	if err != nil {
		log.Printf("could not run the server %v", err)
		panic(err)
	}
}
