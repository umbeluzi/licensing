package licensing

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"
)

// License represents a license with various fields to define its properties.
type License struct {
	// Version represents the version of the license format.
	Version string `json:"version"`

	// Type represents the type of the license (e.g., trial, full, enterprise).
	Type string `json:"type"`

	// Issuer represents the entity that issued the license.
	Issuer string `json:"issuer"`

	// Subject represents the subject to whom the license is issued.
	Subject string `json:"subject"`

	// Audience represents the intended audience for the license.
	Audience []string `json:"audience"`

	// Features represents the features enabled by the license.
	Features []string `json:"features"`

	// Restrictions represents any restrictions associated with the license.
	Restrictions map[string]string `json:"restrictions"`

	// Metadata represents additional metadata associated with the license.
	Metadata map[string]string `json:"metadata"`

	// IssuedAt represents the date when the license was issued.
	IssuedAt string `json:"issued_at"`

	// ExpiresAt represents the date when the license expires.
	ExpiresAt string `json:"expires_at"`

	// Plans represents the subscription or service plans associated with the license.
	Plans []string `json:"plans"`
}

// Generate generates a license key for the given license data using the provided private key.
func Generate(privateKey *rsa.PrivateKey, data License) (string, error) {
	if data.Version == "" {
		data.Version = "1"
	}
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

// Validate validates a given license key using the provided public key.
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

	return data, nil
}

// LoadPublicKeyFromPEM loads an RSA public key from a PEM-encoded string.
func LoadPublicKeyFromPEM(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return rsaPub, nil
}

// CheckIssuer checks if the license issuer matches the expected issuer.
func (l License) CheckIssuer(expectedIssuer string) bool {
	return l.Issuer == expectedIssuer
}

// CheckSubject checks if the license subject matches the expected subject.
func (l License) CheckSubject(expectedSubject string) bool {
	return l.Subject == expectedSubject
}

// CheckAudience checks if the license audience includes the expected audience.
func (l License) CheckAudience(expectedAudience string) bool {
	for _, audience := range l.Audience {
		if audience == expectedAudience {
			return true
		}
	}
	return false
}

// CheckFeature checks if the license features include the expected feature.
func (l License) CheckFeature(expectedFeature string) bool {
	for _, feature := range l.Features {
		if feature == expectedFeature {
			return true
		}
	}
	return false
}

// CheckPlan checks if the license plans include the expected plan.
func (l License) CheckPlan(expectedPlan string) bool {
	for _, plan := range l.Plans {
		if plan == expectedPlan {
			return true
		}
	}
	return false
}

// CheckRestriction checks if the license restrictions include the expected key-value pair.
func (l License) CheckRestriction(key, expectedValue string) bool {
	value, exists := l.Restrictions[key]
	return exists && value == expectedValue
}

// IsValid checks if the license is currently valid based on its expiration date.
func (l License) IsValid() bool {
	expiresAt, err := time.Parse("2006-01-02", l.ExpiresAt)
	if err != nil {
		return false
	}
	return time.Now().Before(expiresAt)
}

// IsExpired checks if the license has expired.
func (l License) IsExpired() bool {
	expiresAt, err := time.Parse("2006-01-02", l.ExpiresAt)
	if err != nil {
		return false
	}
	return time.Now().After(expiresAt)
}

// CheckType checks if the license type matches the expected type.
func (l License) CheckType(expectedType string) bool {
	return l.Type == expectedType
}
