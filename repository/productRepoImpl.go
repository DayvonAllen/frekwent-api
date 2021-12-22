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

type ProductRepoImpl struct {
	product  models.Product
	products []models.Product
}

func (p ProductRepoImpl) Create(product *models.Product) error {
	conn := database.MongoConn

	product.Id = primitive.NewObjectID()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	_, err := conn.ProductCollection.InsertOne(context.TODO(), product)

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (p ProductRepoImpl) FindAll(page string, newProductQuery bool) (*[]models.Product, error) {
	conn := database.MongoConn

	findOptions := options.FindOptions{}
	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}

	findOptions.SetSkip((int64(pageNumber) - 1) * int64(perPage))
	findOptions.SetLimit(int64(perPage))

	if newProductQuery {
		findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	cur, err := conn.ProductCollection.Find(context.TODO(), bson.M{}, &findOptions)

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), p.products); err != nil {
		panic(err)
	}

	// Close the cursor once finished
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			panic(fmt.Errorf("error processing data"))
		}
	}(cur, context.TODO())

	return &p.products, nil
}

func (p ProductRepoImpl) FindByProductId(id primitive.ObjectID) (*models.Product, error) {
	conn := database.MongoConn

	err := conn.ProductCollection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&p.product)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, fmt.Errorf("error processing data")
	}

	return &p.product, nil
}

func (p ProductRepoImpl) UpdateById(product *models.Product) error {
	conn := database.MongoConn

	product.UpdatedAt = time.Now()

	_, err := conn.ProductCollection.UpdateByID(context.TODO(), product.Id, product)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return err
		}
		return fmt.Errorf("error processing data")
	}

	return nil
}

func (p ProductRepoImpl) DeleteById(id primitive.ObjectID) error {
	conn := database.MongoConn

	_, err := conn.ProductCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}})

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func NewProductRepoImpl() ProductRepoImpl {
	var productRepoImpl ProductRepoImpl

	return productRepoImpl
}
