package fib

import "github.com/google/uuid"

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
