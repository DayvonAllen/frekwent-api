package main

import (
	"fmt"
	"freq/database"
	"freq/router"
	"log"
	"os"
	"os/signal"
)

func init() {
	database.ConnectToDB()
}

func main() {
	app := router.Setup()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		_ = <-c
		fmt.Println("Shutting down...")
		_ = app.Shutdown()
	}()

	if err := app.Listen(":8080"); err != nil {
		log.Panic(err)
	}
}
