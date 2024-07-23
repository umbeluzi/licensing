# Umbeluzi Licensing System

This project provides a licensing system similar to GitHub, GitLab, and Sourcegraph licenses.

## Features

- License generation
- License validation
- Secure key management

## License Structure

The `License` struct includes the following fields:
- `ID`: The unique identifier for the license.
- `Type`: The type of the license (e.g., commercial, trial).
- `ExpiresAt`: The expiration date of the license.
- `IssuedAt`: The issuance date of the license.
- `Issuer`: The issuer of the license.
- `Audience`: The audience for the license (can be a single string or an array of strings).
- `Features`: A list of features included in the license.
- `Restrictions`: A map of restrictions associated with the license.
- `Metadata`: A map of additional metadata associated with the license.

## Getting Started

To run the project, navigate to the `cmd/licensing` directory and run the following command:

```sh
go run main.go
```

## Use Cases

### How to Use This Library in an Application

1. **Import the library into your application.**
2. **Generate RSA keys** (in a real-world application, you would store these keys securely).
3. **Generate a license** using the private key.
4. **Validate the license** using the public key.

#### Example Application

**Directory Structure**

```
myapp/
├── cmd/
│   └── myapp/
│       └── main.go
├── go.mod
└── go.sum
```

**`go.mod` File**

Create a `go.mod` file for your application:
```sh
module github.com/yourusername/myapp

go 1.18

require github.com/umbeluzi/licensing latest
```

**Main Application Entry Point**

**`cmd/myapp/main.go`**
```go
package main

import (
    "crypto/rand"
    "crypto/rsa"
    "fmt"
    "time"

    "github.com/umbeluzi/licensing"
)

func main() {
    fmt.Println("My Application with Licensing")

    // Generate RSA keys
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        fmt.Printf("Error generating private key: %v\n", err)
        return
    }
    publicKey := &privateKey.PublicKey

    // Generate a license
    licenseData := licensing.License{
        ID:          "license-123",
        Type:        "commercial",
        ExpiresAt:   "2025-01-01",
        IssuedAt:    time.Now().Format("2006-01-02"),
        Issuer:      "My Application",
        Audience:    []string{"user-123"},
        Features:    []string{"feature1", "feature2"},
        Restrictions: map[string]string{"region": "US"},
        Metadata:    map[string]string{"version": "1.0"},
    }

    licenseKey, err := licensing.Generate(privateKey, licenseData)
    if err != nil {
        fmt.Printf("Error generating license: %v\n", err)
        return
    }

    fmt.Printf("Generated License: %s\n", licenseKey)

    // Validate the license
    validLicenseData, err := licensing.Validate(publicKey, licenseKey)
    if err != nil {
        fmt.Printf("Error validating license: %v\n", err)
        return
    }

    fmt.Printf("License is valid. License Data: %+v\n", validLicenseData)
}
```

**Running the Application**

1. **Initialize the Go module**:
   ```sh
   go mod tidy
   ```

2. **Run the application**:
   ```sh
   go run cmd/myapp/main.go
   ```

## License

This project is licensed under the terms of the [MPL-2.0](../LICENSE).
