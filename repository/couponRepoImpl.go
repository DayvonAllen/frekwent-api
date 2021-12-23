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

type CouponRepoImpl struct {
	Coupon     models.Coupon
	Coupons    []models.Coupon
	CouponList models.CouponList
}

func (c CouponRepoImpl) Create(coupon *models.Coupon) error {
	conn := database.ConnectToDB()

	_, err := c.FindByCode(coupon.Code)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			coupon.Id = primitive.NewObjectID()
			coupon.CreatedAt = time.Now()
			coupon.UpdatedAt = time.Now()

			_, err = conn.CouponCollection.InsertOne(context.TODO(), coupon)

			if err != nil {
				return fmt.Errorf("error processing data")
			}

			return nil
		}
		return fmt.Errorf("error processing data")
	}

	return errors.New("coupon already exists")
}

func (c CouponRepoImpl) FindAll(page string, newCouponQuery bool) (*models.CouponList, error) {
	conn := database.ConnectToDB()

	findOptions := options.FindOptions{}
	perPage := 10
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return nil, fmt.Errorf("page must be a number")
	}

	findOptions.SetSkip((int64(pageNumber) - 1) * int64(perPage))
	findOptions.SetLimit(int64(perPage))

	if newCouponQuery {
		findOptions.SetSort(bson.D{{"createdAt", -1}})
	}

	cur, err := conn.CouponCollection.Find(context.TODO(), bson.M{}, &findOptions)

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &c.Coupons); err != nil {
		panic(err)
	}

	if c.Coupons == nil {
		return nil, errors.New("no coupons found")
	}

	// Close the cursor once finished
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			panic(fmt.Errorf("error processing data"))
		}
	}(cur, context.TODO())

	count, err := conn.CouponCollection.CountDocuments(context.TODO(), bson.M{})

	if err != nil {
		panic(err)
	}

	c.CouponList.NumberOfCoupons = count

	if c.CouponList.NumberOfCoupons < 10 {
		c.CouponList.NumberOfPages = 1
	} else {
		c.CouponList.NumberOfPages = int(count/10) + 1
	}

	c.CouponList.Coupons = &c.Coupons
	c.CouponList.CurrentPage = pageNumber

	return &c.CouponList, nil
}

func (c CouponRepoImpl) FindByCode(code string) (*models.Coupon, error) {
	conn := database.ConnectToDB()

	err := conn.CouponCollection.FindOne(context.TODO(), bson.D{{"code", code}}).Decode(&c.Coupon)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, fmt.Errorf("error processing data")
	}

	return &c.Coupon, nil
}

func (c CouponRepoImpl) DeleteByCode(code string) error {
	conn := database.ConnectToDB()

	_, err := conn.CouponCollection.DeleteOne(context.TODO(), bson.D{{"code", code}})

	if err != nil {
		return fmt.Errorf("error processing data")
	}

	return nil
}

func NewCouponRepoImpl() CouponRepoImpl {
	var couponRepoImpl CouponRepoImpl

	return couponRepoImpl
}
