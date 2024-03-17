package moengage

import (
	"errors"
	"fmt"
	"moengage/helpers"
	"moengage/internal"
	"moengage/pkg/moengage/inform"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	baseURL    string
	appId      string
	dataKey    string
	informKey  string
	pushKey    string
	httpClient http.Client
	//Data       moengage.Data
	Inform inform.Inform
}

func NewClientFromEnv(options ...func(*Client)) (Client, error) {
	if os.Getenv("MOENGAGE_BASE_URL") == "" {
		return Client{}, errors.New("MOENGAGE_BASE_URL environment variable is not set")
	}

	if os.Getenv("MOENGAGE_APP_ID") == "" {
		return Client{}, errors.New("MOENGAGE_APP_ID environment variable is not set")
	}

	options = append(options, func(c *Client) {
		c.dataKey = os.Getenv("MOENGAGE_DATA_KEY")
		c.informKey = os.Getenv("MOENGAGE_INFORM_KEY")
		c.pushKey = os.Getenv("MOENGAGE_PUSH_KEY")
	})

	return NewClient(os.Getenv("MOENGAGE_BASE_URL"), os.Getenv("MOENGAGE_APP_ID"), options...)
}

func NewClient(baseURL string, appId string, options ...func(*Client)) (Client, error) {
	apiURL, err := validateURL(baseURL)
	if err != nil {
		return Client{}, err
	}

	informURL, err := validateURL(fmt.Sprintf("%s-%s", "inform", baseURL))
	if err != nil {
		return Client{}, err
	}

	c := Client{baseURL: apiURL, appId: appId, httpClient: http.Client{}}

	for _, opt := range options {
		opt(&c)
	}

	helpers.Log(helpers.Debug, fmt.Sprintf("Base URL: %s", apiURL))
	helpers.Log(helpers.Debug, fmt.Sprintf("Inform URL: %s", informURL))

	c.Inform = &inform.Channel{
		ReqHandler: internal.HTTPHandler{AppId: c.appId, BaseURL: informURL, APIKey: c.informKey, HTTPClient: c.httpClient},
	}

	return c, nil
}

func validateURL(baseURL string) (string, error) {
	_, err := url.ParseRequestURI(baseURL)
	if err != nil {
		baseURL = "https://" + baseURL
		_, err = url.ParseRequestURI(baseURL)
	}

	return baseURL, err
}
