package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MailMember struct {
	Id              primitive.ObjectID `bson:"_id" json:"id"`
	MemberFirstName string             `bson:"memberFirstName" json:"memberFirstName"`
	MemberLastName  string             `bson:"memberLastName" json:"memberLastName"`
	MemberEmail     string             `bson:"memberEmail" json:"memberEmail"`
	CreatedAt       time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}
