package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Customer struct {
	Id              primitive.ObjectID `bson:"_id" json:"id"`
	FirstName       string             `bson:"firstName" json:"firstName"`
	LastName        string             `bson:"lastName" json:"lastName"`
	Email           string             `bson:"email" json:"email"`
	StreetAddress   string             `bson:"streetAddress" json:"streetAddress"`
	OptionalAddress string             `bson:"optionalAddress" json:"optionalAddress"`
	City            string             `bson:"city" json:"city"`
	State           string             `bson:"state" json:"state"`
	ZipCode         string             `bson:"zipCode" json:"zipCode"`
	PurchasedItems  *[]Product         `bson:"purchasedItems" json:"purchasedItems"`
	InfoEmailOptIn  bool               `bson:"infoEmailOptIn" json:"infoEmailOptIn"`
	CreatedAt       time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"-"`
}

type CustomerList struct {
	Customers         *[]Customer `json:"customers"`
	NumberOfCustomers int64       `json:"numberOfCustomers"`
	CurrentPage       int         `json:"currentPage"`
	NumberOfPages     int         `json:"numberOfPages"`
}
