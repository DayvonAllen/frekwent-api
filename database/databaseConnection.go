package database

import (
	"context"
	"freq/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Connection struct {
	*mongo.Client
	ProductCollection  *mongo.Collection
	AdminCollection    *mongo.Collection
	EmailCollection    *mongo.Collection
	CouponCollection   *mongo.Collection
	LoginIPCollection  *mongo.Collection
	CustomerCollection *mongo.Collection
	PurchaseCollection *mongo.Collection
	*mongo.Database
}

func ConnectToDB() *Connection {
	p := config.Config("DB_PORT")
	n := config.Config("DB_NAME")
	h := config.Config("DB_HOST")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(n+h+p))

	if err != nil {
		panic(err)
	}
	// create database
	db := client.Database("frekwent")

	// create collection
	productCollection := db.Collection("products")
	adminCollection := db.Collection("admin")
	emailCollection := db.Collection("emails")
	couponCollection := db.Collection("coupons")
	loginIPCollection := db.Collection("loginIPs")
	customerCollection := db.Collection("customers")
	purchaseCollection := db.Collection("purchase")

	dbConnection := &Connection{
		client,
		productCollection,
		adminCollection,
		emailCollection,
		couponCollection,
		loginIPCollection,
		customerCollection,
		purchaseCollection,
		db,
	}

	return dbConnection
}
