package tests

import (
	"bytes"
	"context"
	"fmt"
	"freq/database"
	"freq/models"
	"freq/router"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

var resource *dockertest.Resource
var pool *dockertest.Pool

func TestProductHandler_FindAll(t *testing.T) {
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	pool = p

	opts := dockertest.RunOptions{
		Repository:   "mongo",
		Tag:          "latest",
		ExposedPorts: []string{"27017"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"27017": {
				{HostIP: "0.0.0.0", HostPort: "27017"},
			},
		},
	}

	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	os.Setenv("DB_URL", "mongodb://localhost:27017")

	defer func(pool *dockertest.Pool, r *dockertest.Resource) {
		err := pool.Purge(r)
		if err != nil {

		}
	}(pool, resource)

	conn := database.ConnectToDB()

	monId, _ := primitive.ObjectIDFromHex("61caa9598eaeae5425c9780f")

	prod := &models.Product{
		Id: monId,
	}

	_, err = conn.ProductCollection.InsertOne(context.TODO(), prod)
	if err != nil {
		return
	}

	tests := []struct {
		description string

		// Test input
		route string

		// Expected output
		expectedError bool
		expectedCode  int
		expectedBody  bool
		evaluator     string
		query         string
		queryValue    string
	}{
		{
			description:   "get all products route",
			route:         "/products",
			expectedError: false,
			expectedCode:  200,
			expectedBody:  true,
			evaluator:     "61caa9598eaeae5425c9780f",
			query:         "",
			queryValue:    "",
		},
		{
			description:   "get all products trending products error",
			route:         "/products",
			expectedError: false,
			expectedCode:  400,
			expectedBody:  true,
			evaluator:     "error",
			query:         "trending",
			queryValue:    "hello",
		},
		{
			description:   "get all products new products error",
			route:         "/products",
			expectedError: false,
			expectedCode:  400,
			expectedBody:  true,
			evaluator:     "error",
			query:         "new",
			queryValue:    "hello",
		},
		{
			description:   "no products",
			route:         "/products",
			expectedError: false,
			expectedCode:  500,
			expectedBody:  true,
			evaluator:     "no products",
			query:         "",
			queryValue:    "",
		},
	}

	// Setup the app as it is done in the main function
	app := router.Setup()

	// Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route
		// from the test case
		req, _ := http.NewRequest(
			"GET",
			test.route,
			nil,
		)

		q := req.URL.Query()
		q.Add(test.query, test.queryValue)
		req.URL.RawQuery = q.Encode()

		if test.description == "no products" {
			_, err2 := conn.ProductCollection.DeleteMany(context.TODO(), bson.M{})
			if err2 != nil {
				return
			}
		}

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		// verify that no error occured, that is not expected
		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses, the next
		// test case needs to be processed
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Reading the response body should work everytime, such that
		// the err variable should be nil
		assert.Nilf(t, err, test.description)

		// Verify, that the response body equals the expected body
		assert.Equalf(t, test.expectedBody, strings.Contains(string(body), test.evaluator), test.description)
	}
}

func TestProductHandler_Create(t *testing.T) {
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	pool = p

	opts := dockertest.RunOptions{
		Repository:   "mongo",
		Tag:          "latest",
		ExposedPorts: []string{"27017"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"27017": {
				{HostIP: "0.0.0.0", HostPort: "27017"},
			},
		},
	}

	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	os.Setenv("DB_URL", "mongodb://localhost:27017")

	defer func(pool *dockertest.Pool, r *dockertest.Resource) {
		err := pool.Purge(r)
		if err != nil {

		}
	}(pool, resource)

	//b, _ := json.Marshal(&models.User{Email: "a@a.com", Username: "a", Password: "pass"})

	tests := []struct {
		description string

		// Test input
		route string

		// Expected output
		expectedError bool
		expectedCode  int
		expectedBody  bool
		evaluator     string
		query         string
		queryValue    string
		body          []byte
	}{
		{
			description:   "unauthorized",
			route:         "/iriguchi/items",
			expectedError: false,
			expectedCode:  401,
			expectedBody:  true,
			evaluator:     "",
			query:         "",
			queryValue:    "",
			body:          nil,
		},
		//{
		//	description:   "logged in - create products",
		//	route:         "/iriguchi/items",
		//	expectedError: false,
		//	expectedCode:  201,
		//	expectedBody:  true,
		//	evaluator:     "",
		//	query:         "",
		//	queryValue:    "",
		//	body:          b,
		//},
	}

	conn := database.ConnectToDB()
	conn.AdminCollection.InsertOne(context.TODO(), &models.User{Email: "a@a.com", Username: "a", Password: "pass"})

	// Setup the app as it is done in the main function
	app := router.Setup()

	// Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route
		// from the test case
		req, _ := http.NewRequest(
			"POST",
			test.route,
			nil,
		)

		q := req.URL.Query()
		q.Add(test.query, test.queryValue)
		req.URL.RawQuery = q.Encode()

		if strings.Contains(test.description, "logged in") {
			data := url.Values{"email": {}}
			data.Set("body", "{email: a@a.com, password: pass}")
			re, _ := http.NewRequest(
				"POST",
				"/iriguchi/auth/login",
				bytes.NewBufferString(data.Encode()),
			)

			myR, _ := app.Test(re, -1)

			c := myR.Cookies()

			fmt.Println(c)

			req.AddCookie(c[0])
		}

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		// verify that no error occured, that is not expected
		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses, the next
		// test case needs to be processed
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Reading the response body should work everytime, such that
		// the err variable should be nil
		assert.Nilf(t, err, test.description)

		// Verify, that the response body equals the expected body
		assert.Equalf(t, test.expectedBody, strings.Contains(string(body), test.evaluator), test.description)
	}
}
