package repository

import (
	"context"
	"fmt"
	"freq/config"
	"freq/database"
	"freq/helper"
	"freq/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"sync"
)

type CustomerRepoImpl struct {
	customer     models.Customer
	customers    []models.Customer
	customerList models.CustomerList
}

func (c CustomerRepoImpl) Create(customer *models.Customer) error {
	conn := database.ConnectToDB()

	err := conn.CustomerCollection.FindOne(context.TODO(), bson.D{
		{"firstName", customer.FirstName},
		{"lastName", customer.LastName},
		{"email", customer.Email},
	}).Decode(&c.customer)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			_, err = conn.CustomerCollection.InsertOne(context.TODO(), customer)

			if err != nil {
				return fmt.Errorf("error processing data")
			}

			return nil
		}
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (c CustomerRepoImpl) FindAll(page string, newCustomerQuery bool) (*models.CustomerList, error) {
	conn := database.ConnectToDB()

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

	if err = cur.All(context.TODO(), &c.customers); err != nil {
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

	decryptedCustomers := make([]models.Customer, 0, len(c.customers))

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
		decryptedCustomers = append(decryptedCustomers, customer)
	}

	count, err := conn.CustomerCollection.CountDocuments(context.TODO(), bson.M{})

	if err != nil {
		panic(err)
	}

	c.customerList.NumberOfCustomers = count

	if c.customerList.NumberOfCustomers < 10 {
		c.customerList.NumberOfPages = 1
	} else {
		c.customerList.NumberOfPages = int(count/10) + 1
	}

	c.customerList.Customers = &decryptedCustomers
	c.customerList.CurrentPage = pageNumber

	return &c.customerList, nil
}

func (c CustomerRepoImpl) FindAllByFullName(firstName string, lastName string, page string, newLoginQuery bool) (*models.CustomerList, error) {
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

	cur, err := conn.CustomerCollection.Find(context.TODO(), bson.D{{"firstName", firstName}, {"lastName",
		lastName}}, &findOptions)

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &c.customers); err != nil {
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

	decryptedCustomers := make([]models.Customer, 0, len(c.customers))

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
		decryptedCustomers = append(decryptedCustomers, customer)
	}

	count, err := conn.CustomerCollection.CountDocuments(context.TODO(), bson.M{})

	if err != nil {
		panic(err)
	}

	c.customerList.NumberOfCustomers = count

	if c.customerList.NumberOfCustomers < 10 {
		c.customerList.NumberOfPages = 1
	} else {
		c.customerList.NumberOfPages = int(count/10) + 1
	}

	c.customerList.Customers = &decryptedCustomers
	c.customerList.CurrentPage = pageNumber

	return &c.customerList, nil
}

func NewCustomerRepoImpl() CustomerRepoImpl {
	var customerRepoImpl CustomerRepoImpl

	return customerRepoImpl
}
