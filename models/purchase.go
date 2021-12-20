package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Purchase struct {
	Id                       primitive.ObjectID `bson:"_id" json:"_"`
	FirstName                string             `bson:"firstName" json:"firstName" validate:"max=60,min=1"`
	LastName                 string             `bson:"lastName" json:"lastName" validate:"max=60,min=1"`
	Email                    string             `bson:"email" json:"email" validate:"max=60,min=1"`
	StreetAddress            string             `bson:"streetAddress" json:"streetAddress" validate:"max=60,min=1"`
	OptionalAddress          string             `bson:"optionalAddress" json:"optionalAddress" validate:"max=60"`
	City                     string             `bson:"city" json:"city" validate:"max=60,min=1"`
	State                    string             `bson:"state" json:"state" validate:"max=60,min=1"`
	ZipCode                  string             `bson:"zipCode" json:"zipCode" validate:"max=60,min=1"`
	PurchasedItems           *[]Product         `bson:"purchasedItems" json:"purchasedItems"`
	FinalPrice               int16              `bson:"finalPrice" json:"finalPrice"`
	CouponUsed               Coupon             `bson:"couponUsed" json:"couponUsed"`
	CreditCardNumber         string             `bson:"creditCardNumber" json:"creditCardNumber" validate:"max=60,min=1"`
	CreditCardExpirationDate string             `bson:"creditCardExpirationDate" json:"creditCardExpirationDate" validate:"max=60,min=1"`
	CreditCardSecurityCode   string             `bson:"creditCardSecurityCode" json:"creditCardSecurityCode" validate:"max=60,min=1"`
	PurchaseConfirmationId   string             `bson:"purchaseConfirmationId" json:"purchaseConfirmationId"`
	TrackingId               string             `bson:"trackingId" json:"trackingId" validate:"max=150,min=1"`
	Tax                      int16              `bson:"tax" json:"tax"`
	Shipped                  bool               `bson:"shipped" json:"shipped"`
	Delivered                bool               `bson:"delivered" json:"delivered"`
	InfoEmailOptIn           bool               `bson:"infoEmailOptIn" json:"infoEmailOptIn"`
	CreatedAt                time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt                time.Time          `bson:"updatedAt" json:"-"`
}
