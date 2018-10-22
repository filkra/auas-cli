package sql

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	EmptyResponse = ""
)

const (
	FormParamQuery = "query"
)

const (
	QueryURL string = "plagcheck/query-to-csv"
)

func (c *Client) Query(query string) (string, error) {
	// Parse the request URL
	u, err := c.BaseURL.Parse(QueryURL)
	if err != nil {
		return EmptyResponse, err
	}

	// Create form parameter
	data := url.Values{}
	data.Set(FormParamQuery, query)

	// Create a query containing the SQL query
	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return EmptyResponse, err
	}

	// Set client credentials
	req.SetBasicAuth(c.credentials.Username, c.credentials.Password)

	// Set content type and length
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	// Send a POST request containing the SQL query
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return EmptyResponse, err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return EmptyResponse, err
	}

	return string(bodyBytes), nil
}
