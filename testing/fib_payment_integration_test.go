package testing

import (
	"fib-go/fib"
	"github.com/go-chrono/chrono"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

const (
	amount      int    = 1000
	currency    string = "IQD"
	callbackURL string = "http://localhost:1337/"
)

func TestCreatePayment(t *testing.T) {
	// given
	client, _ := fib.New(clientID, clientSecret, true)

	// when
	createdPaymentResponse, _ := client.CreatePayment(amount, currency, callbackURL)

	// then
	paymentID := createdPaymentResponse.PaymentID
	QRCode := createdPaymentResponse.QRCode
	readableCode := createdPaymentResponse.ReadableCode
	validUntil := createdPaymentResponse.ValidUntil
	assert.Equal(t, reflect.TypeOf(paymentID), reflect.TypeOf(uuid.New()))
	assert.NotEmpty(t, readableCode)
	assert.NotEmpty(t, QRCode)
	assert.NotEmpty(t, validUntil)

}

func TestCreatePaymentWithDescription(t *testing.T) {
	// given
	client, _ := fib.New(clientID, clientSecret, true)
	description := "This is a random desc"

	// when
	createdPaymentResponse, _ := client.CreatePayment(amount, currency, callbackURL, fib.WithDescription(description))

	// then
	paymentID := createdPaymentResponse.PaymentID
	QRCode := createdPaymentResponse.QRCode
	readableCode := createdPaymentResponse.ReadableCode
	validUntil := createdPaymentResponse.ValidUntil
	assert.Equal(t, reflect.TypeOf(paymentID), reflect.TypeOf(uuid.New()))
	assert.NotEmpty(t, readableCode)
	assert.NotEmpty(t, QRCode)
	assert.NotEmpty(t, validUntil)

}

func TestCreatePaymentWithExpiresIn(t *testing.T) {
	// given
	client, _ := fib.New(clientID, clientSecret, true)
	expiresIn := chrono.DurationOf(3*chrono.Hour + 7*chrono.Minute + 500*chrono.Millisecond)

	// when
	createdPaymentResponse, _ := client.CreatePayment(amount, currency, callbackURL, fib.WithExpiresIn(expiresIn.String()))

	// then
	paymentID := createdPaymentResponse.PaymentID
	QRCode := createdPaymentResponse.QRCode
	readableCode := createdPaymentResponse.ReadableCode
	validUntil := createdPaymentResponse.ValidUntil
	assert.Equal(t, reflect.TypeOf(paymentID), reflect.TypeOf(uuid.New()))
	assert.NotEmpty(t, readableCode)
	assert.NotEmpty(t, QRCode)
	assert.NotEmpty(t, validUntil)

}

func TestCreatePaymentFails(t *testing.T) {
	// given
	client, _ := fib.New(clientID, clientSecret, true)
	fakeCallBack := "evilwebsite"

	// when
	_, err := client.CreatePayment(amount, currency, fakeCallBack)

	// then
	assert.NotEmpty(t, err.ErrorBody.TraceID)
	assert.Equal(t, err.ErrorBody.Errors[0].Code, "INVALID_REQUEST")
	assert.Equal(t, err.ErrorBody.Errors[0].Title, "general_invalid_request_title")
	assert.Equal(t, err.ErrorBody.Errors[0].Detail, "general_invalid_request_details")

}
