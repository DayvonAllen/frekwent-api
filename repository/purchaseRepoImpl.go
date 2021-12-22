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

type PurchaseRepoImpl struct {
	purchase  models.Purchase
	purchases []models.Purchase
}

func (p PurchaseRepoImpl) Purchase(purchase *models.Purchase) error {
	conn := database.MongoConn

	key := config.Config("KEY")

	encrypt := helper.Encryption{Key: []byte(key)}

	var wg sync.WaitGroup
	wg.Add(9)

	go func() {
		defer wg.Done()
		purchase.Id = primitive.NewObjectID()
		purchase.PurchaseConfirmationId = helper.RandomString(32)
		purchase.Shipped = false
		purchase.Delivered = false
		purchase.TrackingId = ""
		purchase.CreatedAt = time.Now()
		purchase.UpdatedAt = time.Now()
	}()

	go func() {
		defer wg.Done()
		cc, err := encrypt.Encrypt(purchase.CreditCardNumber)

		if err != nil {
			panic(err)
		}

		purchase.CreditCardNumber = cc
	}()

	go func() {
		defer wg.Done()
		cc, err := encrypt.Encrypt(purchase.CreditCardSecurityCode)

		if err != nil {
			panic(err)
		}

		purchase.CreditCardSecurityCode = cc
	}()

	go func() {
		defer wg.Done()
		cc, err := encrypt.Encrypt(purchase.CreditCardExpirationDate)

		if err != nil {
			panic(err)
		}

		purchase.CreditCardExpirationDate = cc
	}()

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(purchase.StreetAddress)

		if err != nil {
			panic(err)
		}

		purchase.StreetAddress = pi
	}()

	go func() {
		defer wg.Done()

		if len(purchase.OptionalAddress) > 0 {
			pi, err := encrypt.Encrypt(purchase.OptionalAddress)

			if err != nil {
				panic(err)
			}

			purchase.OptionalAddress = pi
		}
	}()

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(purchase.City)

		if err != nil {
			panic(err)
		}

		purchase.City = pi
	}()

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(purchase.State)

		if err != nil {
			panic(err)
		}

		purchase.State = pi
	}()

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(purchase.ZipCode)

		if err != nil {
			panic(err)
		}

		purchase.ZipCode = pi
	}()

	wg.Wait()

	_, err := conn.PurchaseCollection.InsertOne(context.TODO(), purchase)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (p PurchaseRepoImpl) FindAll(page string, newPurchaseQuery bool) (*[]models.Purchase, error) {
	conn := database.MongoConn

	findOptions := options.FindOptions{}
	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}

	findOptions.SetSkip((int64(pageNumber) - 1) * int64(perPage))
	findOptions.SetLimit(int64(perPage))

	if newPurchaseQuery {
		findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	cur, err := conn.PurchaseCollection.Find(context.TODO(), bson.M{}, &findOptions)

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), p.purchases); err != nil {
		panic(err)
	}

	// Close the cursor once finished
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			panic(fmt.Errorf("error processing data"))
		}
	}(cur, context.TODO())

	return &p.purchases, nil
}

func (p PurchaseRepoImpl) FindByPurchaseById(id primitive.ObjectID) (*models.Purchase, error) {
	conn := database.MongoConn

	key := config.Config("KEY")

	decrypt := helper.Encryption{Key: []byte(key)}

	err := conn.PurchaseCollection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&p.purchase)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, fmt.Errorf("error processing data")
	}

	var wg sync.WaitGroup
	wg.Add(8)

	go func() {
		defer wg.Done()
		cc, err := decrypt.Decrypt(p.purchase.CreditCardNumber)

		if err != nil {
			panic(err)
		}

		p.purchase.CreditCardNumber = cc
	}()

	go func() {
		defer wg.Done()
		cc, err := decrypt.Decrypt(p.purchase.CreditCardSecurityCode)

		if err != nil {
			panic(err)
		}

		p.purchase.CreditCardSecurityCode = cc
	}()

	go func() {
		defer wg.Done()
		cc, err := decrypt.Decrypt(p.purchase.CreditCardExpirationDate)

		if err != nil {
			panic(err)
		}

		p.purchase.CreditCardExpirationDate = cc
	}()

	go func() {
		defer wg.Done()
		pi, err := decrypt.Decrypt(p.purchase.StreetAddress)

		if err != nil {
			panic(err)
		}

		p.purchase.StreetAddress = pi
	}()

	go func() {
		defer wg.Done()

		if len(p.purchase.OptionalAddress) > 0 {
			pi, err := decrypt.Decrypt(p.purchase.OptionalAddress)

			if err != nil {
				panic(err)
			}

			p.purchase.OptionalAddress = pi
		}
	}()

	go func() {
		defer wg.Done()
		pi, err := decrypt.Decrypt(p.purchase.City)

		if err != nil {
			panic(err)
		}

		p.purchase.City = pi
	}()

	go func() {
		defer wg.Done()
		pi, err := decrypt.Decrypt(p.purchase.State)

		if err != nil {
			panic(err)
		}

		p.purchase.State = pi
	}()

	go func() {
		defer wg.Done()
		pi, err := decrypt.Decrypt(p.purchase.ZipCode)

		if err != nil {
			panic(err)
		}

		p.purchase.ZipCode = pi
	}()

	wg.Wait()

	return &p.purchase, nil
}

func (p PurchaseRepoImpl) UpdateShippedStatus(dto *models.PurchaseShippedDTO) error {
	conn := database.MongoConn

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", dto.Id}}
	update := bson.D{{"$set", bson.D{{"shipped", dto.Shipped},
		{"trackingId", dto.TrackingId}}}}

	conn.PurchaseCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return nil
}

func (p PurchaseRepoImpl) UpdateDeliveredStatus(dto *models.PurchaseDeliveredDTO) error {
	conn := database.MongoConn

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", dto.Id}}
	update := bson.D{{"$set", bson.D{{"delivered", dto.Delivered}}}}

	conn.PurchaseCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return nil
}

func (p PurchaseRepoImpl) UpdatePurchaseAddress(dto *models.PurchaseAddressDTO) error {
	conn := database.MongoConn

	key := config.Config("KEY")

	encrypt := helper.Encryption{Key: []byte(key)}

	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(dto.StreetAddress)

		if err != nil {
			panic(err)
		}

		dto.StreetAddress = pi
	}()

	go func() {
		defer wg.Done()

		if len(dto.OptionalAddress) > 0 {
			pi, err := encrypt.Encrypt(dto.OptionalAddress)

			if err != nil {
				panic(err)
			}

			dto.OptionalAddress = pi
		}
	}()

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(dto.City)

		if err != nil {
			panic(err)
		}

		dto.City = pi
	}()

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(dto.State)

		if err != nil {
			panic(err)
		}

		dto.State = pi
	}()

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(dto.ZipCode)

		if err != nil {
			panic(err)
		}

		dto.ZipCode = pi
	}()

	wg.Wait()

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", dto.Id}}
	update := bson.D{{"$set", bson.D{{"streetAddress", dto.StreetAddress},
		{"optionalAddress", dto.OptionalAddress},
		{"city", dto.City},
		{"state", dto.State},
		{"zipCode", dto.ZipCode}}}}

	conn.PurchaseCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return nil
}

func (p PurchaseRepoImpl) UpdateTrackingNumber(dto *models.PurchaseTrackingDTO) error {
	conn := database.MongoConn

	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", dto.Id}}
	update := bson.D{{"$set", bson.D{{"trackingId", dto.TrackingId}}}}

	conn.PurchaseCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts)

	return nil
}

func NewPurchaseRepoImpl() PurchaseRepoImpl {
	var purchaseRepoImpl PurchaseRepoImpl

	return purchaseRepoImpl
}
