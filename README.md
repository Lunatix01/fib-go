# First Iraqi Bank Go SDK 
[![Run Go Tests](https://github.com/Lunatix01/fib-go/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/Lunatix01/fib-go/actions/workflows/test.yml)

Welcome to FIB Go, the official Go SDK for integrating with the First Iraqi Bank's Online Payments Service. This SDK allows you to easily accept payments, check payment statuses, handle refunds, and much more, all within your Go applications.

## Features
1. Easy Payment Creation
2. Payment Status Checks
3. Payment Refunds
4. Secure Authentication
5. Comprehensive Error Handling

## Installation

```bash
go get -u github.com/Lunatix01/fib-go
```

## Documentation
For detailed documentation, code samples, and best practices, please visit our [documentation](https://www.fibgo.wiki/).

## Quick Start

```go
client, err := fib.New(clientID, clientSecret, isTesting)
if err != nil {
    log.Fatalf("Error creating FIB client: %s - %s", err.Title, err.Description)
}

response, paymentErr := client.CreatePayment(500, "IQD", "http://callback.url")
if paymentErr != nil {
    log.Fatal("Error creating payment:", paymentErr.ErrorBody)
}
```

## Error Handling
For comprehensive information on how to handle errors, please refer to the Error Handling section on our [website](https://fibgo.wiki/error-handling).

## Contributing
We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for more details.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE.md) file for details
