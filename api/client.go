package api

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

const (
	baseUrl        string = "https://auas.cs.uni-duesseldorf.de/"
)

type Client struct {
	BaseURL *url.URL
	httpClient *http.Client
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
	user, present := os.LookupEnv("AUAS_USER")
	if present == false {
		log.Fatal("Please specify a username within the environment variable AUAS_USER")
	}

	// Get the password
	password, present := os.LookupEnv("AUAS_PASS")
	if present == false {
		log.Fatal("Please specify a password within the environment variable AUAS_PASS")
	}

	// Login using the client
	err = ret.Login(user, password)
	if err != nil {
		log.Fatal(err)
	}

	return ret, nil
}
