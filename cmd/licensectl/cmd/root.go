package cmd

import (
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "log"
    "os"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "licensectl",
    Short: "A command line tool for managing licenses",
    Long:  "licensectl is a CLI tool to generate and validate licenses.",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        log.Fatal(err)
    }
}

func init() {
    cobra.OnInitialize(initConfig)
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /home/edson/.licensectl.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
    if cfgFile != "" {
        // Use config file from the flag.
        viper.SetConfigFile(cfgFile)
    } else {
        // Find home directory.
        home, err := os.UserHomeDir()
        if err != nil {
            log.Fatal(err)
        }

        // Search config in home directory with name ".licensectl" (without extension).
        viper.AddConfigPath(home)
        viper.SetConfigName(".licensectl")
    }

    viper.AutomaticEnv() // read in environment variables that match

    // If a config file is found, read it in.
    if err := viper.ReadInConfig(); err == nil {
        log.Println("Using config file:", viper.ConfigFileUsed())
    }
}
