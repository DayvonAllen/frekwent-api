package repository

import (
	"context"
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
	"sync"
	"time"
)

type CustomerRepoImpl struct {
	customer  models.Customer
	customers []models.Customer
}

func (c CustomerRepoImpl) Create(customer *models.Customer) error {
	conn := database.MongoConn

	key := config.Config("KEY")

	encrypt := helper.Encryption{Key: []byte(key)}

	var wg sync.WaitGroup
	wg.Add(5)

	customer.Id = primitive.NewObjectID()
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(customer.StreetAddress)

		if err != nil {
			panic(err)
		}

		customer.StreetAddress = pi
	}()

	go func() {
		defer wg.Done()

		if len(customer.OptionalAddress) > 0 {
			pi, err := encrypt.Encrypt(customer.OptionalAddress)

			if err != nil {
				panic(err)
			}

			customer.OptionalAddress = pi
		}
	}()

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(customer.City)

		if err != nil {
			panic(err)
		}

		customer.City = pi
	}()

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(customer.State)

		if err != nil {
			panic(err)
		}

		customer.State = pi
	}()

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(customer.ZipCode)

		if err != nil {
			panic(err)
		}

		customer.ZipCode = pi
	}()

	wg.Wait()

	_, err := conn.CustomerCollection.InsertOne(context.TODO(), customer)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (c CustomerRepoImpl) FindAll(page string, newCustomerQuery bool) (*[]models.Customer, error) {
	conn := database.MongoConn

	findOptions := options.FindOptions{}
	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}

	findOptions.SetSkip((int64(pageNumber) - 1) * int64(perPage))
	findOptions.SetLimit(int64(perPage))

	if newCustomerQuery {
		findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	cur, err := conn.CustomerCollection.Find(context.TODO(), bson.M{}, &findOptions)

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), c.customers); err != nil {
		panic(err)
	}

	// Close the cursor once finished
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			panic(fmt.Errorf("error processing data"))
		}
	}(cur, context.TODO())

	var wg sync.WaitGroup
	key := config.Config("KEY")

	encrypt := helper.Encryption{Key: []byte(key)}

	for _, customer := range c.customers {
		wg.Add(5)
		go func() {
			defer wg.Done()
			pi, err := encrypt.Decrypt(customer.StreetAddress)

			if err != nil {
				panic(err)
			}

			customer.StreetAddress = pi
		}()

		go func() {
			defer wg.Done()

			if len(customer.OptionalAddress) > 0 {
				pi, err := encrypt.Decrypt(customer.OptionalAddress)

				if err != nil {
					panic(err)
				}

				customer.OptionalAddress = pi
			}
		}()

		go func() {
			defer wg.Done()
			pi, err := encrypt.Decrypt(customer.City)

			if err != nil {
				panic(err)
			}

			customer.City = pi
		}()

		go func() {
			defer wg.Done()
			pi, err := encrypt.Decrypt(customer.State)

			if err != nil {
				panic(err)
			}

			customer.State = pi
		}()

		go func() {
			defer wg.Done()
			pi, err := encrypt.Decrypt(customer.ZipCode)

			if err != nil {
				panic(err)
			}

			customer.ZipCode = pi
		}()

		wg.Wait()
	}

	return &c.customers, nil
}

func (c CustomerRepoImpl) FindAllByFullName(firstName string, lastName string, page string, newLoginQuery bool) (*[]models.Customer, error) {
	conn := database.MongoConn

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

	cur, err := conn.CustomerCollection.Find(context.TODO(), bson.D{{"firstName", firstName}, {"lastName",
		lastName}}, &findOptions)

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), c.customers); err != nil {
		panic(err)
	}

	// Close the cursor once finished
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			panic(fmt.Errorf("error processing data"))
		}
	}(cur, context.TODO())

	return &c.customers, nil
}

func NewCustomerRepoImpl() CustomerRepoImpl {
	var customerRepoImpl CustomerRepoImpl

	return customerRepoImpl
}
