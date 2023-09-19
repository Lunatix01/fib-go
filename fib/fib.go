package fib

import "time"

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
