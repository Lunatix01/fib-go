package testing

import (
	"github.com/lunatix01/fib-go/fib"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var clientID = os.Getenv("CLIENT_ID")
var clientSecret = os.Getenv("CLIENT_SECRET")

const (
	someFakeName string = "Luna"
)

func TestClient(t *testing.T) {
	// when
	client, _ := fib.New(clientID, clientSecret, true)

	// then
	assert.NotEmpty(t, client.URL)
	assert.NotEmpty(t, client.Tokens)
	assert.NotEmpty(t, client.GrantType)
	assert.NotEmpty(t, client.Authentication.ClientId)
	assert.NotEmpty(t, client.Authentication.ClientSecret)
}

func TestClientWithFakeClientIDFail(t *testing.T) {
	// when
	fakeClientID := "luna"
	_, err := fib.New(fakeClientID, clientSecret, true)

	// then
	assert.Equal(t, err.Title, "unauthorized_client")
	assert.Equal(t, err.Description, "INVALID_CREDENTIALS: Invalid client credentials")

}

func TestClientWithFakeClientSecretFail(t *testing.T) {
	// when
	_, err := fib.New(clientID, someFakeName, true)

	// then
	assert.Equal(t, err.Title, "unauthorized_client")
	assert.Equal(t, err.Description, "Invalid client secret")

}
