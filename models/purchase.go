package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Purchase struct {
	Id                     primitive.ObjectID `bson:"_id" json:"id"`
	FirstName              string             `bson:"firstName" json:"firstName" validate:"max=60,min=1"`
	LastName               string             `bson:"lastName" json:"lastName" validate:"max=60,min=1"`
	Email                  string             `bson:"email" json:"email" validate:"max=60,min=1"`
	StreetAddress          string             `bson:"streetAddress" json:"streetAddress" validate:"max=60,min=1"`
	OptionalAddress        string             `bson:"optionalAddress" json:"optionalAddress" validate:"max=60"`
	City                   string             `bson:"city" json:"city" validate:"max=60,min=1"`
	State                  string             `bson:"state" json:"state" validate:"max=60,min=1"`
	ZipCode                string             `bson:"zipCode" json:"zipCode" validate:"max=60,min=1"`
	PurchasedItems         *[]Product         `bson:"purchasedItems" json:"purchasedItems"`
	FinalPrice             int16              `bson:"finalPrice" json:"finalPrice"`
	CouponUsed             Coupon             `bson:"couponUsed" json:"couponUsed"`
	PurchaseConfirmationId string             `bson:"purchaseConfirmationId" json:"purchaseConfirmationId"`
	PurchaseIntent         string             `bson:"purchaseIntent" json:"purchaseIntent"`
	TrackingId             string             `bson:"trackingId" json:"trackingId"`
	//Tax                    int16              `bson:"tax" json:"tax"`
	Shipped        bool      `bson:"shipped" json:"shipped"`
	Delivered      bool      `bson:"delivered" json:"delivered"`
	InfoEmailOptIn bool      `bson:"infoEmailOptIn" json:"infoEmailOptIn"`
	Refunded       bool      `bson:"refunded" json:"refunded"`
	CreatedAt      time.Time `bson:"createdAt" json:"-"`
	UpdatedAt      time.Time `bson:"updatedAt" json:"-"`
}

type PurchaseAddressDTO struct {
	Id              primitive.ObjectID `json:"id"`
	StreetAddress   string             `json:"streetAddress" validate:"max=60,min=1"`
	OptionalAddress string             `json:"optionalAddress" validate:"max=60"`
	City            string             `json:"city" validate:"max=60,min=1"`
	State           string             `json:"state" validate:"max=60,min=1"`
	ZipCode         string             `json:"zipCode" validate:"max=60,min=1"`
}

type PurchaseDeliveredDTO struct {
	Id        primitive.ObjectID `json:"id"`
	Delivered bool               `json:"delivered"`
}

type PurchaseShippedDTO struct {
	Id         primitive.ObjectID `json:"id"`
	Shipped    bool               `json:"shipped"`
	TrackingId string             `json:"trackingId" validate:"max=150,min=1"`
}

type PurchaseTrackingDTO struct {
	Id         primitive.ObjectID `json:"id"`
	TrackingId string             `json:"trackingId" validate:"max=150,min=1"`
}

type PurchaseList struct {
	Purchases         *[]Purchase `json:"purchases"`
	NumberOfPurchases int64       `json:"numberOfPurchases"`
	CurrentPage       int         `json:"currentPage"`
	NumberOfPages     int         `json:"numberOfPages"`
}
