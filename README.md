# First Iraqi Bank Go SDK
[![Go Reference](https://pkg.go.dev/badge/github.com/lunatix01/fib-go)](https://pkg.go.dev/github.com/lunatix01/fib-go)
[![Run Go Tests](https://github.com/Lunatix01/fib-go/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/Lunatix01/fib-go/actions/workflows/test.yml)

Welcome to FIB Go, the Go SDK for integrating with the First Iraqi Bank's Online Payments Service. This SDK allows you to easily accept payments, check payment statuses, handle refunds, and much more, all within your Go applications.

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
For comprehensive information on how to handle errors, please refer to the Error Handling section on our [Website](https://fibgo.wiki/payment/error-handling).

## Contributing
We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for more details.

### Contributing to Documentation

The documentation is available in this [Repository](https://github.com/Lunatix01/fib-go-doc)

If you find an error, or omission, or have an idea for improving the documentation, we would love to hear from you! Feel free to open an issue or create a pull request.

1. **Open an Issue**: If you find a problem or have a suggestion for the documentation, start by opening an issue to discuss it.

2. **Create a Pull Request**: Once the issue is acknowledged, you can proceed to fork the repository, clone it locally, and make your changes. Please make sure your pull request is linked to the issue.

3. **Branch Naming Convention**: Use the following naming conventions for your branches:
    - For documentation improvements: `improvement/doc-BRANCH-NAME`

4. **Testing**: Make sure to test your changes locally before submitting the pull request.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
