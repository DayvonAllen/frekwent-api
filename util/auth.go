package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"freq/config"
	"freq/helper"
	"freq/models"
	"github.com/globalsign/mgo/bson"
	"github.com/golang-jwt/jwt"
	"strconv"
	"strings"
	"time"
)

type Authentication struct {
	Id       bson.ObjectId
	Username string `bson:"username" json:"username"`
}

type LoginDetails struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

type Claims struct {
	jwt.StandardClaims
	Id       bson.ObjectId
	Username string
}

var k = config.Config("SECRET")

func (l Authentication) GenerateJWT(msg models.User) (string, error) {
	e, err := strconv.Atoi(config.Config("EXPIRATION"))

	if err != nil {
		return "", errors.New("authentication Failure")
	}

	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(e) * time.Minute).Unix(),
		},
		Id:       msg.Id,
		Username: msg.Username,
	}
	// always better to use a pointer with JSON
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	signedString, err := token.SignedString([]byte(k))

	if err != nil {
		return "", errors.New("authentication Failure")
	}
	return signedString, nil
}

func (l Authentication) SignToken(token []byte) ([]byte, error) {
	// second arg is a private key, key needs to be the same size as hasher
	// sha512 is 64 bits
	h := hmac.New(sha256.New, []byte(k))

	// hash is a writer
	_, err := h.Write(token)
	if err != nil {
		return nil, errors.New("authentication Failure")
	}

	return []byte(fmt.Sprintf("%x", h.Sum(nil))), nil
}

func (l Authentication) VerifySignature(token, sig []byte) (bool, error) {
	// sign message
	s, err := l.SignToken(token)
	// compare it
	return hmac.Equal(sig, s), err
}

func (l Authentication) IsLoggedIn(tokenValue string) (*Authentication, bool, error) {
	if tokenValue == "" {
		return nil, false, fmt.Errorf("no token")
	}

	data, err := helper.ExtractData(tokenValue)

	if err != nil {
		return nil, false, errors.New("authentication Failure")
	}

	validSig, err := l.VerifySignature([]byte(data[0]), []byte(data[1]))

	if err != nil {
		return nil, false, errors.New("authentication Failure")
	}

	if !validSig {
		return nil, false, errors.New("authentication Failure")
	}

	token, err := jwt.ParseWithClaims(data[0], &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() == jwt.SigningMethodHS256.Alg() {
			//verify token(we pass in our key to be verified)
			return []byte(k), nil
		}
		return nil, errors.New("authentication Failure")
	})

	if err != nil {
		return nil, false, errors.New("authentication Failure")
	}

	isEqual := token.Valid

	if isEqual {
		// user is logged in at this point
		// because we receive an interface type we need to assert which type we want to use that inherits it
		claims := token.Claims.(*Claims)

		l.Id = claims.Id
		l.Username = strings.ToLower(claims.Username)
		return &l, true, nil
	}

	return nil, false, fmt.Errorf("token is not valid")
}
