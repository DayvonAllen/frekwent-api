package repository

import (
	"context"
	"fmt"
	"freq/database"
	"freq/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

type LoginIpRepoImpl struct {
	loginIp  models.LoginIP
	loginIps []models.LoginIP
}

func (l LoginIpRepoImpl) Create(ip *models.LoginIP) error {
	_, err := l.FindByIp(ip.IpAddress)

	if err != nil {
		conn := database.ConnectToDB()

		ip.Id = primitive.NewObjectID()
		ip.AccessedAt = time.Now()
		ip.CreatedAt = time.Now()
		ip.UpdatedAt = time.Now()

		_, err = conn.LoginIPCollection.InsertOne(context.TODO(), ip)

		if err != nil {
			return fmt.Errorf("error processing data")
		}
	}

	err = l.UpdateLoginIp(ip)

	if err != nil {
		return err
	}

	return nil
}

func (l LoginIpRepoImpl) FindAll(page string, newLoginQuery bool) (*[]models.LoginIP, error) {
	conn := database.ConnectToDB()

	findOptions := options.FindOptions{}
	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}

	findOptions.SetSkip((int64(pageNumber) - 1) * int64(perPage))
	findOptions.SetLimit(int64(perPage))

	if newLoginQuery {
		findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	cur, err := conn.LoginIPCollection.Find(context.TODO(), bson.M{}, &findOptions)

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &l.loginIps); err != nil {
		panic(err)
	}

	// Close the cursor once finished
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			panic(fmt.Errorf("error processing data"))
		}
	}(cur, context.TODO())

	return &l.loginIps, nil
}

func (l LoginIpRepoImpl) FindByIp(ip string) (*models.LoginIP, error) {
	conn := database.ConnectToDB()

	err := conn.LoginIPCollection.FindOne(context.TODO(), bson.D{{"ipAddress", ip}}).Decode(&l.loginIp)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, fmt.Errorf("error processing data")
	}

	return &l.loginIp, nil
}

func (l LoginIpRepoImpl) UpdateLoginIp(ip *models.LoginIP) error {
	conn := database.ConnectToDB()

	ip.AccessedAt = time.Now()
	ip.UpdatedAt = time.Now()

	_, err := conn.LoginIPCollection.UpdateByID(context.TODO(), ip.Id, ip)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return err
		}
		return fmt.Errorf("error processing data")
	}

	return nil
}

func NewLoginIpRepoImpl() LoginIpRepoImpl {
	var loginIpRepoImpl LoginIpRepoImpl

	return loginIpRepoImpl
}
