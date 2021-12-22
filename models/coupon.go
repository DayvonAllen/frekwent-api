package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Coupon struct {
	Id             primitive.ObjectID `bson:"_id" json:"id"`
	Percentage     uint8              `bson:"percentage" json:"percentage"`
	Code           string             `bson:"code" json:"code"`
	ExpirationDate time.Time          `bson:"expirationDate" json:"expirationDate"`
	CreatedAt      time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt      time.Time          `bson:"updatedAt" json:"-"`
}
