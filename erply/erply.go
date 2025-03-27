package erply

import (
	"context"
	"fmt"
	"os"

	"github.com/imroc/req/v3"
	"go.uber.org/ratelimit"
	"golang.org/x/sync/semaphore"
)

type service struct {
	client *Client
}

type Client struct {
    *req.Client
    rateLimiter ratelimit.Limiter
    gate *semaphore.Weighted

	CRM *CRMService
}

type ClientConfig struct {
	Username string
	Password string
	Endpoint string
	ClientCode string
}

func NewClientConfigFromEnv() (*ClientConfig, error) {
	uN, exists := os.LookupEnv("ERPLY_USERNAME")
	if !exists {
		return nil, fmt.Errorf("ERPLY_USERNAME not set.")
	}

	pW, exists := os.LookupEnv("ERPLY_PASSWORD")
	if !exists {
		return nil, fmt.Errorf("ERPLY_PASSWORD not set.")
	}

	eP, exists := os.LookupEnv("ERPLY_ENDPOINT")
	if !exists {
		return nil, fmt.Errorf("ERPLY_ENDPOINT not set.")
	}

	cC, exists := os.LookupEnv("ERPLY_CLIENTCODE")
	if !exists {
		return nil, fmt.Errorf("ERPLY_CLIENTCODE not set.")
	}

	return &ClientConfig{
		Username: uN,
		Password: pW,
		Endpoint: eP,
		ClientCode: cC,
	}, nil

}

func NewClient(config *ClientConfig) (*Client, error) {
    c := Client{
        req.C(),
        nil,
        nil,
		nil,
    }

	c.CRM = &CRMService{ client: &c }	

    var vUser ErplyResponse[VerifyUserRecord]
    res, err := c.R().
        SetQueryParams(map[string]string {
            "request": "verifyUser",
            "clientCode": config.ClientCode,
            "username": config.Username,
            "password": config.Password,
			"sessionLength": "86400",
        }).
        SetSuccessResult(&vUser).
        Post(fmt.Sprintf("https://%s/api", config.Endpoint))
    if err != nil {
        return nil, err
    }

    if res.IsSuccessState() {
        sK := vUser.Records[0].SessionKey
        c.
            SetCommonHeader("sessionKey", sK).
            SetCommonHeader("clientCode", config.ClientCode).
            SetCommonHeaderNonCanonical("accept", "application/json").
			SetBaseURL(fmt.Sprintf("https://api-crm-%s/v1", config.Endpoint))
        return withMaxConcurrent(withThrottler(&c)), nil
    }


    return &c, err
}



func withThrottler(c *Client) *Client {
    c.rateLimiter = ratelimit.New(5)
    c.WrapRoundTripFunc(func (rt req.RoundTripper) req.RoundTripFunc {
        return func(req *req.Request) (resp *req.Response, err error) {
            c.rateLimiter.Take()
		    return rt.RoundTrip(req)
	    }
    })
    return c
}

func withMaxConcurrent(c *Client) *Client {
    c.gate = semaphore.NewWeighted(10)
    c.WrapRoundTripFunc(func (rt req.RoundTripper) req.RoundTripFunc {
        return func(req *req.Request) (resp *req.Response, err error) {
            ctx := context.Background()
            c.gate.Acquire(ctx, 1)
            defer c.gate.Release(1)
		    return rt.RoundTrip(req)
	    }
    })
    return c
}

type VerifyUserRecord struct {
    UserID                     string `json:"userID"`
    UserName                   string `json:"userName"`
    EmployeeID                 string `json:"employeeID"`
    EmployeeName               string `json:"employeeName"`
    GroupID                    string `json:"groupID"`
    GroupName                  string `json:"groupName"`
    IPAddress                  string `json:"ipAddress"`
    SessionKey                 string `json:"sessionKey"`
    SessionLength              int    `json:"sessionLength"`
    LoginURL                   string `json:"loginUrl"`
    IdentityToken string `json:"identityToken"`
    Token         string `json:"token"`
}

type Params map[string]string

type Status struct {
    RequestUnixTime int    `json:"requestUnixTime"`
    ResponseStatus  string `json:"responseStatus"`
    ErrorCode       int    `json:"errorCode"`
	RequestID 		int    `json:"requestID"`
}

type ErplyResponse[T any] struct {
	Status Status`json:"status"`
	Records []T `json:"records"`
}
