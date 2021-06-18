package main

import (
	"fmt"
	"os"

	"github.com/clarke94/roulette-service/cmd/serve"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	var config string

	rootCmd := &cobra.Command{
		Use:   "roulette-service [cmd]",
		Short: "roulette-service",
	}

	rootCmd.PersistentFlags().StringVar(
		&config,
		"config",
		"",
		"config file passed in for different development environments.",
	)

	if err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config")); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(
		serve.New(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func initConfig() {
	if config := viper.GetString("config"); config != "" {
		fmt.Println("using config: ", config)
		viper.SetConfigFile(config)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
}
