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
	"time"
)

type MailMemberRepoImpl struct {
	mailMember  models.MailMember
	mailMembers []models.MailMember
}

func (m MailMemberRepoImpl) Create(mm *models.MailMember) error {
	conn := database.ConnectToDB()

	err := conn.MailMemberCollection.FindOne(context.TODO(), bson.D{{"memberEmail", mm.MemberEmail}}).Decode(&m.mailMember)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			mm.Id = primitive.NewObjectID()
			mm.CreatedAt = time.Now()
			mm.UpdatedAt = time.Now()

			_, err := conn.MailMemberCollection.InsertOne(context.TODO(), mm)

			if err != nil {
				return fmt.Errorf("error processing data")
			}

			return nil
		}
		return fmt.Errorf("error processing data")
	}

	return errors.New("mail member with that email already exists")
}

func (m MailMemberRepoImpl) FindAll() (*[]models.MailMember, error) {
	conn := database.ConnectToDB()

	cur, err := conn.MailMemberCollection.Find(context.TODO(), bson.M{})

	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &m.mailMembers); err != nil {
		panic(err)
	}

	if m.mailMembers == nil {
		return nil, errors.New("no mail members in the database")
	}

	// Close the cursor once finished
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			panic(fmt.Errorf("error processing data"))
		}
	}(cur, context.TODO())

	return &m.mailMembers, nil
}

func (m MailMemberRepoImpl) DeleteById(id primitive.ObjectID) error {
	conn := database.ConnectToDB()

	res, _ := conn.MailMemberCollection.DeleteOne(context.TODO(), bson.D{{"_id", id}})

	if res.DeletedCount == 0 {
		return errors.New("failed to delete Mail Member")
	}

	return nil
}

func NewMailMemberRepoImpl() MailMemberRepo {
	var mailMemberRepoImpl MailMemberRepo

	return mailMemberRepoImpl
}
