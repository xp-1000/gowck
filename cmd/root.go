/*
Copyright © 2020 XP-1000 <xp-1000@hotmail.fr>

This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package cmd

import (
	"log"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// nolint: gochecknoglobals
	cfgFile string
)

func setupRunCmd(rootCmd *cobra.Command) {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run",
		Long:  `Running webchecks`,
		RunE: func(cmd *cobra.Command, args []string) error {
			l := viper.GetString("log")
			url := viper.GetString("url")
			return run(url, l)
		},
	}

	rootCmd.AddCommand(runCmd)
	runCmd.Flags().String("url", "https://google.com", "url to check")

	if err := viper.BindPFlag("url", runCmd.Flags().Lookup("url")); err != nil {
		log.Fatal("Unable to bind flag:", err)
	}
}

func Execute() {
	// Setup global root command
	var rootCmd = &cobra.Command{
		Use:   "gowck",
		Short: "Golang written webchecker cli tool",
		Long: `Gowck is a CLI tool allowing to perform simple webchecks.
	It will gather useful metrics from configured endpoints.`,
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "yaml configuration file (default: ~/.gowk)")
	rootCmd.PersistentFlags().String("log", "knights.of.ni", "Syslog host")

	if err := viper.BindPFlag("log", rootCmd.PersistentFlags().Lookup("log")); err != nil {
		log.Fatal("Unable to bind flag:", err)
	}

	setupRunCmd(rootCmd)
	cobra.OnInitialize(initConfig)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gowck")
		viper.AddConfigPath(home)
	}

	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.

	if err := viper.ReadInConfig(); err != nil {
		if cfgFile != "" {
			log.Println("config specified but unable to read it, using defaults")
		}
	}
}
