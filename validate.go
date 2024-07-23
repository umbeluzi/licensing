package licensing

import (
    "crypto/rsa"
    "crypto/sha256"
    "encoding/base64"
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
