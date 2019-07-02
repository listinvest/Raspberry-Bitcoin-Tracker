package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	//Default httpClient
	httpClient     *http.Client
	coindeskConfig config
	respChan       chan *Response
	//Channel for error handling
	errChan chan error
	//limits the request
	requestsPerSecond time.Duration
}
type config struct {
	host string
	path string
}

func NewClient(timer int) *Client {
	return &Client{
		//Config for Coindesk API
		coindeskConfig: config{
			host: "https://api.coindesk.com/",
			path: "v1/bpi/currentprice.json",
		},
		respChan:          make(chan *Response),
		errChan:           make(chan error),
		requestsPerSecond: time.Second * time.Duration(timer), //Makes a request every given second
	}
}

// Make a http GET Request to the Server
func (c *Client) get() {
	req, err := http.NewRequest(http.MethodGet, c.coindeskConfig.host+c.coindeskConfig.path, nil)
	if err != nil {
		c.errChan <- errors.New("NewReqeust failed: " + c.coindeskConfig.host + c.coindeskConfig.path)
		return
	}

	//start a new goroutine for each request when timer is over
	throttle := time.Tick(c.requestsPerSecond)
	for {
		go c.do(req)
		<-throttle
	}
}

func (c *Client) do(req *http.Request) {
	client := c.httpClient
	if client == nil {
		client = http.DefaultClient
	}

	httpResp, err := client.Do(req)
	if err != nil {
		c.errChan <- fmt.Errorf("client failed to do request %v", err)
		return
	}
	defer httpResp.Body.Close()

	decoder := json.NewDecoder(httpResp.Body)

	var s Response
	err = decoder.Decode(&s)
	c.respChan <- &s
}

//Reponse from Coindesk
type Response struct {
	Time struct {
		Updated    string    `json:"updated"`
		UpdatedISO time.Time `json:"updatedISO"`
	} `json:"time"`
	ChartName string `json:"chartName"`
	Bpi       struct {
		USD struct {
			Code      string  `json:"code"`
			Rate      string  `json:"rate"`
			RateFloat float64 `json:"rate_float"`
		} `json:"USD"`
		EUR struct {
			Code      string  `json:"code"`
			Rate      string  `json:"rate"`
			RateFloat float64 `json:"rate_float"`
		} `json:"EUR"`
	} `json:"bpi"`
}
