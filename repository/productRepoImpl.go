package repository

import (
	"context"
	"errors"
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
	product     models.Product
	products    []models.Product
	productList models.ProductList
}

func (p ProductRepoImpl) Create(product *models.Product) error {
	conn := database.ConnectToDB()

	err := conn.ProductCollection.FindOne(context.TODO(), bson.D{{"name", product.Name}}).Decode(&p.product)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			product.Id = primitive.NewObjectID()
			product.CreatedAt = time.Now()
			product.UpdatedAt = time.Now()

			_, err := conn.ProductCollection.InsertOne(context.TODO(), product)

			if err != nil {
				return fmt.Errorf("error processing data")
			}

			return nil
		}
		return fmt.Errorf("error processing data")
	}

	return errors.New("product with that name already exists")
}

func (p ProductRepoImpl) FindAll(page string, newProductQuery bool) (*models.ProductList, error) {
	conn := database.ConnectToDB()

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

	if err = cur.All(context.TODO(), &p.products); err != nil {
		panic(err)
	}

	if p.products == nil {
		return nil, errors.New("no products in the database")
	}

	// Close the cursor once finished
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			panic(fmt.Errorf("error processing data"))
		}
	}(cur, context.TODO())

	count, err := conn.ProductCollection.CountDocuments(context.TODO(), bson.M{})

	if err != nil {
		panic(err)
	}

	p.productList.NumberOfProducts = count

	if p.productList.NumberOfProducts < 10 {
		p.productList.NumberOfPages = 1
	} else {
		p.productList.NumberOfPages = int(count/10) + 1
	}

	p.productList.Products = &p.products
	p.productList.CurrentPage = pageNumber

	return &p.productList, nil
}

func (p ProductRepoImpl) FindAllByCategory(category string, page string, newProductQuery bool) (*models.ProductList, error) {
	conn := database.ConnectToDB()

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

	cur, err := conn.ProductCollection.Find(context.TODO(), bson.D{{"category", category}}, &findOptions)

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &p.products); err != nil {
		panic(err)
	}

	if p.products == nil {
		return nil, errors.New("no products in the database")
	}

	// Close the cursor once finished
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			panic(fmt.Errorf("error processing data"))
		}
	}(cur, context.TODO())

	count, err := conn.ProductCollection.CountDocuments(context.TODO(), bson.D{{"category", category}})

	if err != nil {
		panic(err)
	}

	p.productList.NumberOfProducts = count

	if p.productList.NumberOfProducts < 10 {
		p.productList.NumberOfPages = 1
	} else {
		p.productList.NumberOfPages = int(count/10) + 1
	}

	p.productList.Products = &p.products
	p.productList.CurrentPage = pageNumber

	return &p.productList, nil
}

func (p ProductRepoImpl) FindByProductId(id primitive.ObjectID) (*models.Product, error) {
	conn := database.ConnectToDB()

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

func (p ProductRepoImpl) FindByProductName(name string) (*models.Product, error) {
	conn := database.ConnectToDB()

	err := conn.ProductCollection.FindOne(context.TODO(), bson.D{{"name", name}}).Decode(&p.product)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, fmt.Errorf("error processing data")
	}

	return &p.product, nil
}

func (p ProductRepoImpl) UpdateName(name string, id primitive.ObjectID) (*models.Product, error) {
	conn := database.ConnectToDB()

	prod := new(models.Product)

	err := conn.ProductCollection.FindOne(context.TODO(), bson.D{{"name", name}}).Decode(prod)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			opts := options.FindOneAndUpdate()
			filter := bson.D{{"_id", id}}
			update := bson.D{{"$set", bson.D{{"name", name},
				{"updatedAt", time.Now()}}}}

			err = conn.ProductCollection.FindOneAndUpdate(context.TODO(),
				filter, update, opts).Decode(&p.product)

			if err != nil {
				return nil, err
			}

			p.product.Name = name

			return &p.product, nil
		}
		return nil, fmt.Errorf("error processing data")
	}

	return nil, errors.New("product with that name already exists")
}

func (p ProductRepoImpl) UpdateQuantity(quantity uint16, id primitive.ObjectID) (*models.Product, error) {
	conn := database.ConnectToDB()

	opts := options.FindOneAndUpdate()
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"quantity", quantity},
		{"updatedAt", time.Now()}}}}

	err := conn.ProductCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&p.product)

	if err != nil {
		return nil, err
	}

	p.product.Quantity = quantity

	return &p.product, nil
}

func (p ProductRepoImpl) UpdatePrice(price string, id primitive.ObjectID) (*models.Product, error) {
	conn := database.ConnectToDB()

	opts := options.FindOneAndUpdate()
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"price", price},
		{"updatedAt", time.Now()}}}}

	err := conn.ProductCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&p.product)

	if err != nil {
		return nil, err
	}

	p.product.Price = price

	return &p.product, nil
}

func (p ProductRepoImpl) UpdateDescription(desc string, id primitive.ObjectID) (*models.Product, error) {
	conn := database.ConnectToDB()

	opts := options.FindOneAndUpdate()
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"description", desc},
		{"updatedAt", time.Now()}}}}

	err := conn.ProductCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&p.product)

	if err != nil {
		return nil, err
	}

	p.product.Description = desc

	return &p.product, nil
}

func (p ProductRepoImpl) UpdateIngredients(ingredients *[]string, id primitive.ObjectID) (*models.Product, error) {
	conn := database.ConnectToDB()

	opts := options.FindOneAndUpdate()
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"ingredients", ingredients},
		{"updatedAt", time.Now()}}}}

	err := conn.ProductCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&p.product)

	if err != nil {
		return nil, err
	}

	p.product.Ingredients = *ingredients

	return &p.product, nil
}

func (p ProductRepoImpl) UpdateCategory(category string, id primitive.ObjectID) (*models.Product, error) {
	conn := database.ConnectToDB()

	opts := options.FindOneAndUpdate()
	filter := bson.D{{"_id", id}}
	update := bson.D{{"$set", bson.D{{"category", category},
		{"updatedAt", time.Now()}}}}

	err := conn.ProductCollection.FindOneAndUpdate(context.TODO(),
		filter, update, opts).Decode(&p.product)

	if err != nil {
		return nil, err
	}

	p.product.Category = category

	return &p.product, nil
}

func (p ProductRepoImpl) DeleteById(id primitive.ObjectID) error {
	conn := database.ConnectToDB()

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
