package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Purchase struct {
	Id                       primitive.ObjectID `bson:"_id" json:"_"`
	FirstName                string             `bson:"firstName" json:"firstName"`
	LastName                 string             `bson:"lastName" json:"lastName"`
	Email                    string             `bson:"email" json:"email"`
	StreetAddress            string             `bson:"streetAddress" json:"streetAddress"`
	OptionalAddress          string             `bson:"optionalAddress" json:"optionalAddress"`
	City                     string             `bson:"city" json:"city"`
	State                    string             `bson:"state" json:"state"`
	ZipCode                  string             `bson:"zipCode" json:"zipCode"`
	PurchasedItems           *[]Product         `bson:"purchasedItems" json:"purchasedItems"`
	FinalPrice               int16              `bson:"finalPrice" json:"finalPrice"`
	CouponUsed               Coupon             `bson:"couponUsed" json:"couponUsed"`
	CreditCardNumber         string             `bson:"creditCardNumber" json:"creditCardNumber"`
	CreditCardExpirationDate string             `bson:"creditCardExpirationDate" json:"creditCardExpirationDate"`
	CreditCardSecurityCode   string             `bson:"creditCardSecurityCode" json:"creditCardSecurityCode"`
	PurchaseConfirmationId   string             `bson:"purchaseConfirmationId" json:"purchaseConfirmationId"`
	TrackingId               string             `bson:"trackingId" json:"trackingId"`
	Tax                      int16              `bson:"tax" json:"tax"`
	Shipped                  bool               `bson:"shipped" json:"shipped"`
	Delivered                bool               `bson:"delivered" json:"delivered"`
	CreatedAt                time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt                time.Time          `bson:"updatedAt" json:"-"`
}
