package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Product struct {
	Id          primitive.ObjectID `bson:"_id" json:"_"`
	Name        string             `bson:"name" json:"name" validate:"max=60,min=1"`
	Images      *[]string          `bson:"images" json:"images"`
	Price       uint16             `bson:"price" json:"price"`
	Quantity    uint16             `bson:"quantity" json:"quantity"`
	Description string             `bson:"description" json:"description" validate:"max=450,min=1"`
	Ingredients *[]string          `bson:"ingredients" json:"ingredients"`
	CreatedAt   time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"-"`
}
