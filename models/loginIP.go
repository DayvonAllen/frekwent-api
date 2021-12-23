package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LoginIP struct {
	Id         primitive.ObjectID `bson:"_id" json:"id"`
	IpAddress  string             `bson:"ipAddress" json:"ipAddress"`
	AccessedAt time.Time          `bson:"accessedAt" json:"accessedAt"`
	CreatedAt  time.Time          `bson:"createdAt" json:"-"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type LoginIpList struct {
	LoginIps         *[]LoginIP `json:"loginIps"`
	NumberOfLoginIps int64      `json:"numberOfLoginIps"`
	CurrentPage      int        `json:"currentPage"`
	NumberOfPages    int        `json:"numberOfPages"`
}
