package fib

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// URLs
const (
	ProdURL = "https://fib.prod.fib.iq/"
	TestURL = "https://fib.stage.fib.iq/"
)

const (
	GrantType = "client_credentials"
)

const (
	AuthenticationPath = "auth/realms/fib-online-shop/protocol/openid-connect/token"
)

// Response status codes
const (
	OK                    = 200
	CREATED               = 201
	NO_CONTENT            = 204
	BAD_CONTENT           = 400
	UNAUTHORIZED          = 401
	NOT_FOUND             = 404
	INTERNAL_SERVER_ERROR = 500
	SERVICE_UNAVAILABLE   = 503
)

// Authentication properties
type Authentication struct {
	ClientId     string
	ClientSecret string
}

// Tokens properties
type Tokens struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresAfter time.Duration `json:"expires_in"`
	ExpiresAt    time.Time
}

// Client struct used in the entire website
type Client struct {
	URL            string
	Authentication Authentication
	Tokens         Tokens
	GrantType      string
}

type LoginError struct {
	Title       string `json:"error"`
	Description string `json:"error_description"`
}

func defaultConfigs(clientId string, clientSecret string, isTesting bool) Client {
	URL := ProdURL
	if isTesting {
		URL = TestURL
	}
	return Client{
		URL: URL,
		Authentication: Authentication{
			ClientId:     clientId,
			ClientSecret: clientSecret,
		},
		Tokens:    Tokens{},
		GrantType: GrantType,
	}
}

// New create a new Client
func New(clientId string, clientSecret string, isTesting bool) (*Client, *LoginError) {
	configs := defaultConfigs(clientId, clientSecret, isTesting)

	loginErr := authenticate(&configs)

	return &configs, loginErr
}

// authenticate function to authenticate user and get back token
func authenticate(configs *Client) *LoginError {
	authenticationURL := configs.URL + AuthenticationPath
	authenticationContentType := "application/x-www-form-urlencoded"
	formData := url.Values{
		"grant_type":    {configs.GrantType},
		"client_id":     {configs.Authentication.ClientId},
		"client_secret": {configs.Authentication.ClientSecret},
	}
	requestBody := strings.NewReader(formData.Encode())

	response, err := http.Post(authenticationURL, authenticationContentType, requestBody)
	if err != nil {
		log.Fatal("ERROR sending request request")
	}

	readableBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error while reading body.")
	}

	if response.StatusCode >= BAD_CONTENT && response.StatusCode < INTERNAL_SERVER_ERROR {
		var loginError LoginError
		unmarshalJSON(readableBody, &loginError)
		return &loginError
	}

	unmarshalJSON(readableBody, &configs.Tokens)

	configs.Tokens.ExpiresAt = time.Now().Add(configs.Tokens.ExpiresAfter * time.Second)

	return nil
}

func (tokens *Tokens) RefreshTokenIfNeeded(configs *Client) {
	if tokens.IsTokenExpired() {
		authenticate(configs)
	}
}

func (tokens *Tokens) IsTokenExpired() bool {
	return time.Now().After(tokens.ExpiresAt)
}

func unmarshalJSON(data []byte, v interface{}) {
	err := json.Unmarshal(data, v)
	if err != nil {
		log.Fatal("Error while parsing response body:", err)
	}
}
