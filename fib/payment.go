package fib

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
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
