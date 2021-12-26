package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Status int32
type EmailType int32

const (
	Success Status = iota
	Pending
	Failed
)

const (
	CouponPromotion EmailType = iota
	Purchased
	CustomerInteraction
)

type Email struct {
	Id            primitive.ObjectID `bson:"_id" json:"id"`
	Type          EmailType          `bson:"type" json:"type"`
	CustomerEmail string             `bson:"customerEmail" json:"customerEmail"`
	From          string             `bson:"from" json:"from"`
	Content       string             `bson:"content" json:"content"`
	Subject       string             `bson:"subject" json:"subject"`
	Status        Status             `bson:"status" json:"status"`
	CreatedAt     time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type EmailDto struct {
	Email      string `json:"email"`
	Content    string `json:"content"`
	Subject    string `json:"subject"`
	CouponCode string `json:"couponCode"`
}

type EmailList struct {
	Emails         *[]Email `json:"emails"`
	NumberOfEmails int64    `json:"numberOfEmails"`
	CurrentPage    int      `json:"currentPage"`
	NumberOfPages  int      `json:"numberOfPages"`
}
