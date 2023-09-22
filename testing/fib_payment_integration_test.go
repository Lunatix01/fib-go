package testing

import (
	"github.com/Lunatix01/fib-go/fib"
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

func TestCheckPayment(t *testing.T) {
	// given
	client, _ := fib.New(clientID, clientSecret, true)
	createdPaymentResponse, _ := client.CreatePayment(amount, currency, callbackURL)

	// when
	paymentCheckResponse, _ := client.CheckPayment(createdPaymentResponse.PaymentID)
	paymentID := paymentCheckResponse.PaymentID
	paymentStatus := paymentCheckResponse.Status
	monetaryValue := paymentCheckResponse.MonetaryValue
	currency := monetaryValue.Currency
	amount := monetaryValue.Amount
	declinedAt := paymentCheckResponse.DeclinedAt
	paidAt := paymentCheckResponse.PaidAt
	paymentDecliningReason := paymentCheckResponse.DecliningReason
	paidBy := paymentCheckResponse.PaidBy
	name := paidBy.Name
	iban := paidBy.IBAN

	// then
	assert.Equal(t, reflect.TypeOf(paymentID), reflect.TypeOf(uuid.New()))
	assert.Equal(t, paymentStatus, fib.UNPAID)
	assert.NotNil(t, monetaryValue)
	assert.Equal(t, currency, "IQD")
	assert.Equal(t, amount, 1000)
	assert.Empty(t, declinedAt)
	assert.Empty(t, paidAt)
	assert.Empty(t, paymentDecliningReason)
	assert.Empty(t, name)
	assert.Empty(t, iban)

}

func TestCheckPaymentNotFound(t *testing.T) {
	// given
	client, _ := fib.New(clientID, clientSecret, true)
	fakeUUID := uuid.New()

	// when
	_, err := client.CheckPayment(fakeUUID)

	// then
	assert.NotEmpty(t, err.ErrorBody.TraceID)
	assert.Equal(t, err.ErrorBody.Errors[0].Code, "PAYMENT_NOT_FOUND")
	assert.Equal(t, err.ErrorBody.Errors[0].Title, "")
	assert.Equal(t, err.ErrorBody.Errors[0].Detail, "")

}

func TestCancelPayment(t *testing.T) {
	// given
	client, _ := fib.New(clientID, clientSecret, true)
	_, _ = client.CreatePayment(amount, currency, callbackURL)
	fakeUUID := uuid.New()

	// when
	_, err := client.CancelPayment(fakeUUID)

	// then
	assert.NotEmpty(t, err.ErrorBody.TraceID)
	assert.Equal(t, err.ErrorBody.Errors[0].Code, "PAYMENT_NOT_FOUND")
	assert.Equal(t, err.ErrorBody.Errors[0].Title, "")
	assert.Equal(t, err.ErrorBody.Errors[0].Detail, "")
}

func TestRefundNotSucceed(t *testing.T) {
	// find a way to test refund ?
}
