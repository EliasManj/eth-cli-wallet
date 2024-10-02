/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/EliasManj/go-wallet/cmd/account"
	"github.com/EliasManj/go-wallet/cmd/network"
	"github.com/EliasManj/go-wallet/cmd/send"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var debugMode bool

var rootCmd = &cobra.Command{
	Use:   "toolbox",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func setDefaults() {
	viper.SetDefault("port", "8080")
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "Enable debug mode")

	cobra.OnInitialize(initConfig)
	setDefaults()

	if debugMode {
		fmt.Println("Initializing root command")
		fmt.Println("name:", viper.Get("name"))
	}

	// Add my subcommand palette
	rootCmd.AddCommand(network.NetworkCmd)
	rootCmd.AddCommand(account.AccountCmd)
	rootCmd.AddCommand(send.SendCmd)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if debugMode {
		fmt.Println("Initializing configuration")
	}

	if cfgFile != "" {
		if debugMode {
			fmt.Println("Using config file from flag:", cfgFile)
		}
		viper.SetConfigFile(cfgFile)
	} else {
		// Add the current directory as a search path
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get current directory: %v\n", err)
			os.Exit(1)
		}
		if debugMode {
			fmt.Println("Current directory:", currentDir)
		}
		viper.AddConfigPath(currentDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err == nil {
		if debugMode {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	} else {
		if debugMode {
			fmt.Fprintf(os.Stderr, "No config file found: %v\n", err)
		}
	}

	viper.AutomaticEnv()
	checkAndCreateDBPath()
}

func checkAndCreateDBPath() {
	if debugMode {
		fmt.Println("CREATING DB PATH")
	}
	dbPath := viper.GetString("database_file_path")
	if dbPath == "" {
		fmt.Fprintln(os.Stderr, "Database file path not set.")
		return
	}
	dir := filepath.Dir(dbPath)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if debugMode {
			fmt.Println("Creating directories for database file path:", dir)
		}
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directories: %v\n", err)
			os.Exit(1)
		}
	}
}
