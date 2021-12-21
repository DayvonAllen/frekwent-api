package repository

import (
	"context"
	"fmt"
	"freq/database"
	"freq/helper"
	"freq/models"
	"freq/util"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepoImpl struct {
	user *models.User
}

func (a AuthRepoImpl) Login(username string, password string, ip string) (*models.User, string, error) {
	var login util.Authentication

	conn := database.MongoConn

	fmt.Println(username)
	fmt.Println(helper.IsEmail(username))

	if helper.IsEmail(username) {
		err := conn.AdminCollection.FindOne(context.TODO(), bson.D{{"email",
			username}}).Decode(a.user)

		if err != nil {
			return nil, "", fmt.Errorf("error finding by email")
		}

		fmt.Println(a.user)

	} else {
		err := conn.AdminCollection.FindOne(context.TODO(), bson.D{{"username",
			username}}).Decode(a.user)

		if err != nil {
			return nil, "", fmt.Errorf("error finding by username")
		}
	}

	fmt.Println(a.user)

	err := bcrypt.CompareHashAndPassword([]byte(a.user.Password), []byte(password))

	if err != nil {
		return nil, "", fmt.Errorf("error comparing password")
	}

	token, err := login.GenerateJWT(*a.user)

	if err != nil {
		return nil, "", fmt.Errorf("error generating token")
	}

	ipAddress := new(models.LoginIP)

	ipAddress.IpAddress = ip

	go func() {
		err := LoginIpRepoImpl{}.Create(ipAddress)
		if err != nil {
			return
		}
		return
	}()

	return a.user, token, nil
}

func NewAuthRepoImpl() AuthRepoImpl {
	var authRepoImpl AuthRepoImpl

	return authRepoImpl
}
