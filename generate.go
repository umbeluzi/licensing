package licensing

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/base64"
    "encoding/pem"
    "encoding/json"
    "errors"
    "fmt"
)

type License struct {
    ID           string            `json:"id"`
    Type         string            `json:"type"`
    ExpiresAt    string            `json:"expires_at"`
    IssuedAt     string            `json:"issued_at"`
    Issuer       string            `json:"issuer"`
    Audience     []string          `json:"audience"`
    Features     []string          `json:"features"`
    Restrictions map[string]string `json:"restrictions"`
    Metadata     map[string]string `json:"metadata"`
}

func Generate(privateKey *rsa.PrivateKey, data License) (string, error) {
    dataJSON, err := json.Marshal(data)
    if err != nil {
        return "", err
    }

    hashed := sha256.Sum256(dataJSON)
    signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, 0, hashed[:])
    if err != nil {
        return "", err
    }

    licenseKey := fmt.Sprintf("%s.%s", base64.StdEncoding.EncodeToString(dataJSON), base64.StdEncoding.EncodeToString(signature))
    return licenseKey, nil
}

func GenerateFromString(privateKeyPEM string, data License) (string, error) {
    block, _ := pem.Decode([]byte(privateKeyPEM))
    if block == nil || block.Type != "RSA PRIVATE KEY" {
        return "", errors.New("failed to decode PEM block containing private key")
    }

    privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return "", err
    }

    return Generate(privateKey, data)
}
