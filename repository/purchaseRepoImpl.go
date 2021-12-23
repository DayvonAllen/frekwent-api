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
	"log"
	"strconv"
	"sync"
	"time"
)

type PurchaseRepoImpl struct {
	purchase     models.Purchase
	purchases    []models.Purchase
	purchaseList models.PurchaseList
}

func (p PurchaseRepoImpl) Purchase(purchase *models.Purchase) error {
	conn := database.ConnectToDB()

	purchase = helper.EncryptPI(purchase)

	customer := new(models.Customer)

	customer.Id = primitive.NewObjectID()
	customer.UpdatedAt = time.Now()
	customer.CreatedAt = time.Now()
	customer.FirstName = purchase.FirstName
	customer.LastName = purchase.LastName
	customer.Email = purchase.Email
	customer.PurchasedItems = purchase.PurchasedItems
	customer.StreetAddress = purchase.StreetAddress
	customer.OptionalAddress = purchase.OptionalAddress
	customer.City = purchase.City
	customer.State = purchase.State
	customer.ZipCode = purchase.ZipCode
	customer.InfoEmailOptIn = purchase.InfoEmailOptIn

	err := CustomerRepoImpl{}.Create(customer)

	_, err = conn.PurchaseCollection.InsertOne(context.TODO(), purchase)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	go func(connection *database.Connection) {
		email := new(models.Email)
		email.CustomerEmail = customer.Email
		email.From = "frekwent@frekwent.com"
		email.Subject = "test subject"
		email.Content = "test content"
		email.Type = "purchase"
		email.Status = "pending"

		err := EmailRepoImpl{}.Create(email)

		if err != nil {
			log.Println(errors.New(fmt.Sprintf("error sending an email to %s", email.CustomerEmail)))
		}
	}(conn)

	return nil
}

func (p PurchaseRepoImpl) FindAll(page string, newPurchaseQuery bool) (*models.PurchaseList, error) {
	conn := database.ConnectToDB()

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
		return nil, errors.New("error finding purchases")
	}

	if err = cur.All(context.TODO(), &p.purchases); err != nil {
		log.Println(errors.New("error finding purchases"))
	}

	if p.purchases == nil {
		return nil, errors.New("no purchases in the database")
	}

	// Close the cursor once finished
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			panic(fmt.Errorf("error processing data"))
		}
	}(cur, context.TODO())

	count, err := conn.PurchaseCollection.CountDocuments(context.TODO(), bson.M{})

	if err != nil {
		return nil, errors.New("error finding purchases")
	}

	p.purchaseList.NumberOfPurchases = count

	if p.purchaseList.NumberOfPurchases < 10 {
		p.purchaseList.NumberOfPages = 1
	} else {
		p.purchaseList.NumberOfPages = int(count/10) + 1
	}

	p.purchaseList.Purchases = &p.purchases
	p.purchaseList.CurrentPage = pageNumber

	return &p.purchaseList, nil
}

func (p PurchaseRepoImpl) FindByPurchaseById(id primitive.ObjectID) (*models.Purchase, error) {
	conn := database.ConnectToDB()

	err := conn.PurchaseCollection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&p.purchase)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(fmt.Sprintf("no purchase with the ID %v", id))
		}
		return nil, fmt.Errorf("error processing data")
	}

	p.purchase = *helper.DecryptPI(&p.purchase)

	return &p.purchase, nil
}

func (p PurchaseRepoImpl) FindByPurchaseConfirmationId(id string) (*models.Purchase, error) {
	conn := database.ConnectToDB()

	err := conn.PurchaseCollection.FindOne(context.TODO(), bson.D{{"purchaseConfirmationId", id}}).Decode(&p.purchase)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(fmt.Sprintf("error finding purchases with purchase confirmation ID %v", id))
		}
		return nil, fmt.Errorf("error processing data")
	}

	p.purchase = *helper.DecryptPI(&p.purchase)

	return &p.purchase, nil
}

func (p PurchaseRepoImpl) UpdateShippedStatus(dto *models.PurchaseShippedDTO) error {
	conn := database.ConnectToDB()

	opts := options.FindOneAndUpdate()
	filter := bson.D{{"_id", dto.Id}}
	update := bson.D{{"$set", bson.D{{"shipped", dto.Shipped},
		{"trackingId", dto.TrackingId}}}}

	err := conn.PurchaseCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&p.purchase)

	if err != nil {
		return errors.New("error updating purchase's shipping status")
	}

	return nil
}

func (p PurchaseRepoImpl) UpdateDeliveredStatus(dto *models.PurchaseDeliveredDTO) error {
	conn := database.ConnectToDB()

	opts := options.FindOneAndUpdate()
	filter := bson.D{{"_id", dto.Id}}
	update := bson.D{{"$set", bson.D{{"delivered", dto.Delivered}}}}

	err := conn.PurchaseCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&p.purchase)

	if err != nil {
		return errors.New("error updating purchase's delivered status")
	}

	return nil
}

func (p PurchaseRepoImpl) UpdatePurchaseAddress(dto *models.PurchaseAddressDTO) error {
	conn := database.ConnectToDB()

	key := config.Config("KEY")

	encrypt := helper.Encryption{Key: []byte(key)}

	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(dto.StreetAddress)

		if err != nil {
			log.Println(err)
		}

		dto.StreetAddress = pi
	}()

	go func() {
		defer wg.Done()

		if len(dto.OptionalAddress) > 0 {
			pi, err := encrypt.Encrypt(dto.OptionalAddress)

			if err != nil {
				log.Println(err)
			}

			dto.OptionalAddress = pi
		}
	}()

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(dto.City)

		if err != nil {
			log.Println(err)
		}

		dto.City = pi
	}()

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(dto.State)

		if err != nil {
			log.Println(err)
		}

		dto.State = pi
	}()

	go func() {
		defer wg.Done()
		pi, err := encrypt.Encrypt(dto.ZipCode)

		if err != nil {
			log.Println(err)
		}

		dto.ZipCode = pi
	}()

	wg.Wait()

	opts := options.FindOneAndUpdate()
	filter := bson.D{{"_id", dto.Id}}
	update := bson.D{{"$set", bson.D{{"streetAddress", dto.StreetAddress},
		{"optionalAddress", dto.OptionalAddress},
		{"city", dto.City},
		{"state", dto.State},
		{"zipCode", dto.ZipCode}}}}

	err := conn.PurchaseCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&p.purchase)

	if err != nil {
		return errors.New("error updating purchase's address")
	}

	return nil
}

func (p PurchaseRepoImpl) UpdateTrackingNumber(dto *models.PurchaseTrackingDTO) error {
	conn := database.ConnectToDB()

	opts := options.FindOneAndUpdate()
	filter := bson.D{{"_id", dto.Id}}
	update := bson.D{{"$set", bson.D{{"trackingId", dto.TrackingId}}}}

	err := conn.PurchaseCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&p.purchase)

	if err != nil {
		return errors.New("error updating purchase's tracking ID")
	}

	return nil
}

func NewPurchaseRepoImpl() PurchaseRepoImpl {
	var purchaseRepoImpl PurchaseRepoImpl

	return purchaseRepoImpl
}
