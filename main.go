package main

import (
	"context"
	"fmt"
	"freq/database"
	"freq/models"
	"freq/router"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"os/signal"
	"time"
)

func init() {
	database.ConnectToDB()

	conn := database.MongoConn

	user := new(models.User)
	err := conn.AdminCollection.FindOne(context.TODO(), bson.D{{"email", "admin@admin.com"}}).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			admin := new(models.User)
			admin.Email = "admin@admin.com"
			admin.Username = "admin"
			admin.Password = "password"
			admin.Id = primitive.NewObjectID()
			admin.CreatedAt = time.Now()
			admin.UpdatedAt = time.Now()

			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
			admin.Password = string(hashedPassword)

			_, err = conn.AdminCollection.InsertOne(context.TODO(), admin)
			if err != nil {
				return
			}
		}
	}
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
