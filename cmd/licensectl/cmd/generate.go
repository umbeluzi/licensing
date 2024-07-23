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
    "time"
)

var generateCmd = &cobra.Command{
    Use:   "generate",
    Short: "Generate a license",
    Long:  "Generate a license using a private key.",
    Run: func(cmd *cobra.Command, args []string) {
        privateKeyPath := viper.GetString("private_key")
        if privateKeyPath == "" {
            log.Fatal("private_key is required")
        }

        privateKey, err := loadPrivateKey(privateKeyPath)
        if err != nil {
            log.Fatal(err)
        }

        licenseData := licensing.License{
            ID:          viper.GetString("id"),
            Type:        viper.GetString("type"),
            ExpiresAt:   viper.GetString("expires_at"),
            IssuedAt:    time.Now().Format("2006-01-02"),
            Issuer:      viper.GetString("issuer"),
            Audience:    viper.GetStringSlice("audience"),
            Features:    viper.GetStringSlice("features"),
            Restrictions: viper.GetStringMapString("restrictions"),
            Metadata:    viper.GetStringMapString("metadata"),
        }

        licenseKey, err := licensing.Generate(privateKey, licenseData)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Println("Generated License:", licenseKey)
    },
}

func init() {
    rootCmd.AddCommand(generateCmd)

    generateCmd.Flags().String("private_key", "", "Path to the private key file")
    generateCmd.Flags().String("id", "", "License ID")
    generateCmd.Flags().String("type", "", "License type")
    generateCmd.Flags().String("expires_at", "", "License expiration date")
    generateCmd.Flags().String("issuer", "", "License issuer")
    generateCmd.Flags().StringSlice("audience", []string{}, "License audience")
    generateCmd.Flags().StringSlice("features", []string{}, "License features")
    generateCmd.Flags().StringToString("restrictions", map[string]string{}, "License restrictions")
    generateCmd.Flags().StringToString("metadata", map[string]string{}, "License metadata")

    viper.BindPFlag("private_key", generateCmd.Flags().Lookup("private_key"))
    viper.BindPFlag("id", generateCmd.Flags().Lookup("id"))
    viper.BindPFlag("type", generateCmd.Flags().Lookup("type"))
    viper.BindPFlag("expires_at", generateCmd.Flags().Lookup("expires_at"))
    viper.BindPFlag("issuer", generateCmd.Flags().Lookup("issuer"))
    viper.BindPFlag("audience", generateCmd.Flags().Lookup("audience"))
    viper.BindPFlag("features", generateCmd.Flags().Lookup("features"))
    viper.BindPFlag("restrictions", generateCmd.Flags().Lookup("restrictions"))
    viper.BindPFlag("metadata", generateCmd.Flags().Lookup("metadata"))
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    block, _ := pem.Decode(data)
    if block == nil || block.Type != "RSA PRIVATE KEY" {
        return nil, errors.New("failed to decode PEM block containing private key")
    }

    return x509.ParsePKCS1PrivateKey(block.Bytes)
}
