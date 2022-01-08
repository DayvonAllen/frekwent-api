package main

import (
	"fmt"
	"freq/database"
	"freq/events"
	"freq/models"
	"freq/router"
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"os/signal"
	"time"
)

func init() {
	//_, err = conn.PurchaseCollection.DeleteMany(context.TODO(), bson.M{})
	//if err != nil {
	//	return
	//}
	//
	//_, err = conn.CustomerCollection.DeleteMany(context.TODO(), bson.M{})
	//if err != nil {
	//	return
	//}
	//_, err = conn.EmailCollection.DeleteMany(context.TODO(), bson.M{})
	//if err != nil {
	//	return
	//}
}

func main() {
	time.Sleep(10 * time.Second)
	conn := database.Sess

	user := new(models.User)
	err := conn.DB(database.DB).C(database.ADMIN).Find(bson.D{{"email", "admin@admin.com"}}).One(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			admin := new(models.User)
			admin.Email = "admin@admin.com"
			admin.Username = "admin"
			admin.Password = "password"
			admin.CreatedAt = time.Now()
			admin.UpdatedAt = time.Now()

			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
			admin.Password = string(hashedPassword)

			err = conn.DB(database.DB).C(database.ADMIN).Insert(admin)
			if err != nil {
				return
			}
		}
	}

	go events.CreateConsumer()

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
