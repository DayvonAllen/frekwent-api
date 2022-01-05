package repository

import (
	"errors"
	"fmt"
	"freq/database"
	"freq/models"
	bson2 "github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type MailMemberRepoImpl struct {
	mailMember  models.MailMember
	mailMembers []models.MailMember
}

func (m MailMemberRepoImpl) Create(mm *models.MailMember) error {
	conn := database.Sess

	fmt.Println(mm)

	err := conn.DB(database.DB).C(database.MAIL_MEMBERS).Find(bson.M{"memberEmail": mm.MemberEmail}).One(&m.mailMember)

	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err.Error() == "not found" {
			mm.CreatedAt = time.Now()
			mm.UpdatedAt = time.Now()

			err = conn.DB(database.DB).C(database.MAIL_MEMBERS).Insert(mm)

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
	conn := database.Sess

	err := conn.DB(database.DB).C(database.MAIL_MEMBERS).Find(nil).All(&m.mailMembers)

	if err != nil {
		return nil, err
	}

	return &m.mailMembers, nil
}

func (m MailMemberRepoImpl) DeleteById(id bson2.ObjectId) error {
	conn := database.Sess

	err := conn.DB(database.DB).C(database.MAIL_MEMBERS).RemoveId(id)

	if err != nil {
		return err
	}

	return nil
}

func NewMailMemberRepoImpl() MailMemberRepo {
	var mailMemberRepoImpl MailMemberRepo

	return mailMemberRepoImpl
}
