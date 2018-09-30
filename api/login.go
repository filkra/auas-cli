package api

import (
	"net/url"
)

const (
	FormParamLogin string = "LoginEmail"
	FormParamPassword string = "LoginPassword"
)

const (
	EmployeLoginURL string = "employee/login"
)

func (c *Client) Login(username string, password string) error {
	// Parse the request URL
	u, err := c.BaseURL.Parse(EmployeLoginURL)
	if err != nil {
		return err
	}

	// Send a POST request containing the credentials
	_, err = c.httpClient.PostForm(u.String(), url.Values {
		FormParamLogin : {username},
		FormParamPassword: {password},
	})

	if err != nil {
		return err
	}

	return nil
}