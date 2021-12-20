package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Email struct {
	Id        primitive.ObjectID `bson:"_id" json:"_"`
	Type      string             `bson:"type" json:"type"`
	To        string             `bson:"to" json:"to"`
	From      string             `bson:"from" json:"from"`
	Content   string             `bson:"content" json:"content"`
	Subject   string             `bson:"subject" json:"subject"`
	CreatedAt time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
