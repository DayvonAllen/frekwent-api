package repository

import (
	"context"
	"errors"
	"fmt"
	"freq/config"
	"freq/database"
	"freq/helper"
	"freq/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

type EmailRepoImpl struct {
	email     models.Email
	emails    []models.Email
	emailList models.EmailList
}

func (e EmailRepoImpl) Create(email *models.Email) error {
	conn := database.ConnectToDB()

	email.Id = primitive.NewObjectID()
	email.CreatedAt = time.Now()
	email.UpdatedAt = time.Now()

	_, err := conn.EmailCollection.InsertOne(context.TODO(), email)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (e EmailRepoImpl) FindAll(page string, newEmailQuery bool) (*models.EmailList, error) {
	conn := database.ConnectToDB()

	findOptions := options.FindOptions{}
	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}

	findOptions.SetSkip((int64(pageNumber) - 1) * int64(perPage))
	findOptions.SetLimit(int64(perPage))

	if newEmailQuery {
		findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	cur, err := conn.EmailCollection.Find(context.TODO(), bson.M{}, &findOptions)

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &e.emails); err != nil {
		panic(err)
	}

	if e.emails == nil {
		return nil, errors.New("no emails in the database")
	}

	// Close the cursor once finished
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			panic(fmt.Errorf("error processing data"))
		}
	}(cur, context.TODO())

	key := config.Config("KEY")

	encrypt := helper.Encryption{Key: []byte(key)}

	decryptedEmail := make([]models.Email, 0, len(e.emails))

	for _, email := range e.emails {
		pi, err := encrypt.Decrypt(email.To)

		if err != nil {
			panic(err)
		}

		email.To = pi

		decryptedEmail = append(decryptedEmail, email)
	}

	count, err := conn.EmailCollection.CountDocuments(context.TODO(), bson.M{})

	if err != nil {
		panic(err)
	}

	e.emailList.NumberOfEmails = count

	if e.emailList.NumberOfEmails < 10 {
		e.emailList.NumberOfPages = 1
	} else {
		e.emailList.NumberOfPages = int(count/10) + 1
	}

	e.emailList.Emails = &decryptedEmail
	e.emailList.CurrentPage = pageNumber

	return &e.emailList, nil
}

func (e EmailRepoImpl) FindAllByEmail(page string, newLoginQuery bool, email string) (*models.EmailList, error) {
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

	cur, err := conn.EmailCollection.Find(context.TODO(), bson.D{{"email", email}}, &findOptions)

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &e.emails); err != nil {
		panic(err)
	}

	if e.emails == nil {
		return nil, errors.New(fmt.Sprintf("No emails found for the email address: %s", email))
	}

	// Close the cursor once finished
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			panic(fmt.Errorf("error processing data"))
		}
	}(cur, context.TODO())

	count, err := conn.EmailCollection.CountDocuments(context.TODO(), bson.M{})

	if err != nil {
		panic(err)
	}

	e.emailList.NumberOfEmails = count

	if e.emailList.NumberOfEmails < 10 {
		e.emailList.NumberOfPages = 1
	} else {
		e.emailList.NumberOfPages = int(count/10) + 1
	}

	e.emailList.Emails = &e.emails
	e.emailList.CurrentPage = pageNumber

	return &e.emailList, nil
}

func NewEmailRepoImpl() EmailRepoImpl {
	var emailRepoImpl EmailRepoImpl

	return emailRepoImpl
}
