package licensing

import (
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/base64"
    "encoding/pem"
    "encoding/json"
    "errors"
    "strings"
    "time"
)

func Validate(publicKey *rsa.PublicKey, licenseKey string) (License, error) {
    parts := strings.Split(licenseKey, ".")
    if len(parts) != 2 {
        return License{}, errors.New("invalid license format")
    }

    dataJSON, err := base64.StdEncoding.DecodeString(parts[0])
    if err != nil {
        return License{}, err
    }

    var data License
    if err := json.Unmarshal(dataJSON, &data); err != nil {
        return License{}, err
    }

    signature, err := base64.StdEncoding.DecodeString(parts[1])
    if err != nil {
        return License{}, err
    }

    hashed := sha256.Sum256(dataJSON)
    if err := rsa.VerifyPKCS1v15(publicKey, 0, hashed[:], signature); err != nil {
        return License{}, err
    }

    expiresAt, err := time.Parse("2006-01-02", data.ExpiresAt)
    if err != nil {
        return License{}, err
    }

    if time.Now().After(expiresAt) {
        return License{}, errors.New("license has expired")
    }

    return data, nil
}

func ValidateFromString(publicKeyPEM string, licenseKey string) (License, error) {
    block, _ := pem.Decode([]byte(publicKeyPEM))
    if block == nil || block.Type != "PUBLIC KEY" {
        return License{}, errors.New("failed to decode PEM block containing public key")
    }

    pub, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return License{}, err
    }

    rsaPub, ok := pub.(*rsa.PublicKey)
    if !ok {
        return License{}, errors.New("not an RSA public key")
    }

    return Validate(rsaPub, licenseKey)
}
