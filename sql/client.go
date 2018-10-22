package sql

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

const (
	baseUrl        string = "https://auas.cs.uni-duesseldorf.de"
)

type Credentials struct {
	Username string
	Password string
}

type Client struct {
	BaseURL *url.URL
	httpClient *http.Client
	credentials Credentials
}

func NewClient(client *http.Client) (*Client, error) {
	// Create a default client if none is specified
	if client == nil {
		client = http.DefaultClient
	}

	// Attach a cookie jar to the client
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	client.Jar = jar

	// Parse the base url and create the http client
	base, _ := url.Parse(baseUrl)
	ret := &Client{httpClient: client, BaseURL: base}

	// Get the username
	user, present := os.LookupEnv("SQL_USER")
	if present == false {
		log.Fatal("Please specify a username within the environment variable SQL_USER")
	}

	// Get the password
	password, present := os.LookupEnv("SQL_PASS")
	if present == false {
		log.Fatal("Please specify a password within the environment variable SQL_PASS")
	}

	ret.credentials = Credentials{Username: user, Password: password}

	return ret, nil
}
