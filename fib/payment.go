package fib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"reflect"
)

// Request methods
const (
	POST = "POST"
	GET  = "GET"
)

// Request URLS
const (
	PaymentBasePath     = "protected/v1/payments"
	PaymentCreationPath = PaymentBasePath
	PaymentCheckPath    = PaymentBasePath + "/%s/status"
	PaymentRefundPath   = PaymentBasePath + "/%s/refund"
	PaymentCancelPath   = PaymentBasePath + "/%s/cancel"
)


// PaymentError Errors return with this type except authentication
type PaymentError struct {
	error     error
	ErrorBody *ErrorBody
}

// ErrorBody body for errors
type ErrorBody struct {
	TraceID string  `json:"traceId"`
	Errors  []Error `json:"errors"`
}

// Error details of the error
type Error struct {
	Code   string `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

// Payment CreatePayment request
type Payment struct {
	MonetaryValue     MonetaryValue `json:"monetaryValue"`
	StatusCallbackURL string        `json:"statusCallbackUrl"`
	Description       string        `json:"description"`
	ExpiresIn         string        `json:"expiresIn"`
}

type CreatePaymentResponse struct {
	PaymentID        uuid.UUID `json:"paymentId"`
	ReadableCode     string    `json:"readableCode"`
	QRCode           string    `json:"QrCode"`
	ValidUntil       string    `json:"validUntil"`
	PersonalAppLink  string    `json:"personalAppLink"`
	BusinessAppLink  string    `json:"businessAppLink"`
	CorporateAppLink string    `json:"corporateAppLink"`
}

// MonetaryValue includes amount and currency type currently ONLY IQD IS SUPPORTED
type MonetaryValue struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

// PaymentStatus returns one constant from below
type PaymentStatus string

const (
	PAID             PaymentStatus = "PAID"
	UNPAID           PaymentStatus = "UNPAID"
	DECLINED         PaymentStatus = "DECLINED"
	REFUND_REQUESTED PaymentStatus = "REFUND_REQUESTED"
	REFUNDED         PaymentStatus = "REFUNDED"
)

// PaymentDecliningReason returns one constant from below
type PaymentDecliningReason string

const (
	SERVER_FAILURE       PaymentDecliningReason = "SERVER_FAILURE"
	PAYMENT_EXPIRATION   PaymentDecliningReason = "PAYMENT_EXPIRATION"
	PAYMENT_CANCELLATION PaymentDecliningReason = "PAYMENT_CANCELLATION"
)

type PaidBy struct {
	Name string `json:"name"`
	IBAN string `json:"iban"`
}

type CheckPaymentResponse struct {
	PaymentID       uuid.UUID              `json:"paymentId"`
	Status          PaymentStatus          `json:"status"`
	PaidAt          string                 `json:"paidAt,omitempty"`
	MonetaryValue   MonetaryValue          `json:"amount"`
	DecliningReason PaymentDecliningReason `json:"decliningReason,omitempty"`
	DeclinedAt      string                 `json:"declinedAt,omitempty"`
	PaidBy          PaidBy                 `json:"paidBy,omitempty"`
}

type PaymentFunc func(*Payment)

// defaultPayment without optional parameters
func defaultPayment(amount int, currency string, statusCallBackURL string) Payment {
	return Payment{
		MonetaryValue: MonetaryValue{
			Amount:   amount,
			Currency: currency,
		},
		StatusCallbackURL: statusCallBackURL,
		Description:       "",
		ExpiresIn:         "",
	}
}

// WithDescription add description optional
func WithDescription(description string) PaymentFunc {
	return func(payment *Payment) {
		payment.Description = description
	}
}

// WithExpiresIn add expiresIn optional
func WithExpiresIn(expiresIn string) PaymentFunc {
	return func(payment *Payment) {
		payment.ExpiresIn = expiresIn
	}
}

// CreatePayment method creates a payment and returns CreatePaymentResponse, PaymentError
func (client *Client) CreatePayment(amount int, currency string, statusCallBackURL string, opts ...PaymentFunc) (CreatePaymentResponse, *PaymentError) {
	var createPaymentResponse CreatePaymentResponse
	payment := defaultPayment(amount, currency, statusCallBackURL)
	for _, opt := range opts {
		opt(&payment)
	}

	marshal, err := json.Marshal(payment)
	if err != nil {
		log.Fatal("cant encode body")
	}

	headers := client.buildHeaders()

	_, newErr := request(client.URL+PaymentCreationPath, headers, marshal, &createPaymentResponse, POST)
	return createPaymentResponse, newErr
}

// CheckPayment method checks payment and its status and returns CheckPaymentResponse, PaymentError
func (client *Client) CheckPayment(paymentID uuid.UUID) (CheckPaymentResponse, *PaymentError) {
	var checkPaymentResponse CheckPaymentResponse

	headers := client.buildHeaders()

	URL := fmt.Sprintf(client.URL+PaymentCheckPath, paymentID)
	_, err := request(URL, headers, nil, &checkPaymentResponse, GET)

	return checkPaymentResponse, err
}

// CancelPayment method to cancel a payment returns bool, PaymentError
func (client *Client) CancelPayment(paymentID uuid.UUID) (bool, *PaymentError) {
	headers := client.buildHeaders()

	URL := fmt.Sprintf(client.URL+PaymentCancelPath, paymentID)
	isCanceled, err := request(URL, headers, nil, nil, POST)

	if err != nil {
		return false, err
	}

	return reflect.ValueOf(isCanceled).Bool(), err
}

// RefundPayment method to refund a payment returns bool, PaymentError
func (client *Client) RefundPayment(paymentID uuid.UUID) (bool, *PaymentError) {
	headers := client.buildHeaders()

	URL := fmt.Sprintf(client.URL+PaymentRefundPath, paymentID)
	isRefundedAlready, err := request(URL, headers, nil, nil, POST)

	if err != nil {
		return false, err
	}

	return reflect.ValueOf(isRefundedAlready).Bool(), err
}

// request function used by other payment methods
func request(URL string, headers map[string]string, body []byte, responseBody interface{}, method string) (interface{}, *PaymentError) {
	req, err := http.NewRequest(method, URL, bytes.NewBuffer(body))
	if err != nil {
		log.Println("NewRequest:", err)
		return nil, &PaymentError{
			error:     err,
			ErrorBody: nil,
		}
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("NewRequest:", err)
		return nil, &PaymentError{
			error:     err,
			ErrorBody: nil,
		}
	}
	defer response.Body.Close()

	statusCode := response.StatusCode

	readableBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Error while reading body.")
		return nil, &PaymentError{
			error:     err,
			ErrorBody: nil,
		}
	}

	switch statusCode {
	case OK:
	case CREATED:
	case ACCEPTED:
		return true, nil
	case BAD_CONTENT, NOT_FOUND:
		var errBody ErrorBody
		if err := json.Unmarshal(readableBody, &errBody); err != nil {
			log.Println("Error unmarshalling ErrorBody.")
			return nil, &PaymentError{
				error:     err,
				ErrorBody: nil,
			}
		}
		return nil, &PaymentError{
			error:     nil,
			ErrorBody: &errBody,
		}
	case NO_CONTENT:
		return true, nil
	case UNAUTHORIZED:
		return nil, &PaymentError{
			error: nil,
			ErrorBody: &ErrorBody{
				TraceID: "",
				Errors: []Error{
					{
						Title:  "Unauthorized",
						Code:   "",
						Detail: "",
					},
				},
			},
		}
	default:
		log.Printf("Unhandled status code: %d", statusCode)
	}

	unmarshalJSON(readableBody, &responseBody)
	return nil, nil
}

func (client *Client) buildHeaders() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + client.Tokens.AccessToken,
		"Content-Type":  "application/json",
	}
}
