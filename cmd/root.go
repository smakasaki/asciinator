/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/smakasaki/asciinator/internal/image_processor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	customMap string
	colored   bool

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "asciinator [image path]", // TODO add url support
		Short: "Transform an image into ASCII art",
		Long:  "This tool will transform an image into ASCII art and will print it to the terminal.\n For now, only local files are supported.",
		Run: func(cmd *cobra.Command, args []string) {
			if !checkArgsAndFlags(args) {
				os.Exit(1)
			}

			flags := image_processor.Flags{
				CustomMap: customMap,
				Colored:   colored,
			}

			for _, imagePath := range args {
				if err := printAscii(imagePath, flags); err != nil {
					return
				}
			}
		},
	}
)

func printAscii(imagePath string, flags image_processor.Flags) error {
	if asciiArt, err := image_processor.Convert(imagePath, flags); err == nil {
		fmt.Println(asciiArt)
	} else {
		fmt.Println(err)
		return err
	}
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.asciinator.yaml)")
	rootCmd.PersistentFlags().StringVarP(&customMap, "custom-map", "m", "", "Custom character map")
	rootCmd.PersistentFlags().BoolVarP(&colored, "colored", "c", false, "Colored output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".asciinator" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".asciinator")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
