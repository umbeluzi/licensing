package cmd

import (
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "errors"
    "fmt"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/umbeluzi/licensing"
    "io/ioutil"
    "log"
)

var validateCmd = &cobra.Command{
    Use:   "validate",
    Short: "Validate a license",
    Long:  "Validate a license using a public key.",
    Run: func(cmd *cobra.Command, args []string) {
        publicKeyPath := viper.GetString("public-key")
        if publicKeyPath == "" {
            log.Fatal("public-key is required")
        }

        licenseKey := viper.GetString("license-key")
        if licenseKey == "" {
            log.Fatal("license-key is required")
        }

        publicKey, err := loadPublicKey(publicKeyPath)
        if err != nil {
            log.Fatal(err)
        }

        licenseData, err := licensing.Validate(publicKey, licenseKey)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf("License is valid. License Data: %+v\n", licenseData)
    },
}

func init() {
    rootCmd.AddCommand(validateCmd)

    validateCmd.Flags().String("public-key", "", "Path to the public key file")
    validateCmd.Flags().String("license-key", "", "License key to validate")

    viper.BindPFlag("public-key", validateCmd.Flags().Lookup("public-key"))
    viper.BindPFlag("license-key", validateCmd.Flags().Lookup("license-key"))
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    block, _ := pem.Decode(data)
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
