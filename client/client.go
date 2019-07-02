package client

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
	RespChan       chan *Response
	//Channel for error handling
	ErrChan chan error
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
		RespChan:          make(chan *Response),
		ErrChan:           make(chan error),
		requestsPerSecond: time.Second * time.Duration(timer), //Makes a request every given second
	}
}

// Make a http GET Request to the Server
func (c *Client) Get() {
	req, err := http.NewRequest(http.MethodGet, c.coindeskConfig.host+c.coindeskConfig.path, nil)
	if err != nil {
		c.ErrChan <- errors.New("NewReqeust failed: " + c.coindeskConfig.host + c.coindeskConfig.path)
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
		c.ErrChan <- fmt.Errorf("client failed to do request %v", err)
		return
	}
	defer httpResp.Body.Close()

	decoder := json.NewDecoder(httpResp.Body)

	var s Response
	err = decoder.Decode(&s)
	c.RespChan <- &s
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
