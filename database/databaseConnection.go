package database

import (
	"freq/config"
	"github.com/globalsign/mgo"
	"log"
	"time"
)

var ADMIN = "admin"
var PRODUCTS = "products"
var EMAILS = "emails"
var COUPONS = "coupons"
var IPS = "loginIPs"
var CUSTOMERS = "customers"
var PURCHASES = "purchases"
var MAIL_MEMBERS = "mailMembers"
var DB = "frekwent"
var Sess = ConnectToDB()

func ConnectToDB() *mgo.Session {
	u := config.Config("DB_URL")

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:     []string{u},
		Timeout:   60 * time.Second,
		PoolLimit: 20,
		Database:  "Frekwent",
	}

	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}

	mongoSession.SetMode(mgo.Monotonic, true)

	return mongoSession
}
