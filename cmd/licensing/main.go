package main

import (
    "crypto/rand"
    "crypto/rsa"
    "fmt"
    "github.com/umbeluzi/licensing"
    "time"
)

func main() {
    fmt.Println("Umbeluzi Licensing System")

    // Generate RSA keys
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        fmt.Printf("Error generating private key: %v\n", err)
        return
    }
    publicKey := &privateKey.PublicKey

    // Example of generating a license
    licenseData := licensing.License{
        ID:          "license-123",
        Type:        "commercial",
        ExpiresAt:   "2025-01-01",
        IssuedAt:    time.Now().Format("2006-01-02"),
        Issuer:      "Umbeluzi Licensing System",
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

    // Example of validating a license
    validLicenseData, err := licensing.Validate(publicKey, licenseKey)
    if err != nil {
        fmt.Printf("Error validating license: %v\n", err)
        return
    }

    fmt.Printf("License is valid. License Data: %+v\n", validLicenseData)
}
