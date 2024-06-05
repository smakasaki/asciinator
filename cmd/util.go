package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/smakasaki/asciinator/internal/art_processor"
	"github.com/smakasaki/asciinator/internal/handlers"
	"github.com/spf13/cobra"
)

var (
	customMap string
	colored   bool
)

func mainCommand(cmd *cobra.Command, args []string) {
	if !checkArgsAndFlags(args) {
		os.Exit(1)
	}

	flags := art_processor.Flags{
		CustomMap: customMap,
		Colored:   colored,
	}

	for _, imagePath := range args {
		if err := processImage(imagePath, flags); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func processImage(imagePath string, flags art_processor.Flags) error {
	chain := handlers.CreateChain()
	asciiArt, err := chain.Handle(imagePath, flags)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(asciiArt)
	return nil
}

// checkArgsAndFlags returns false if an error occurs
func checkArgsAndFlags(args []string) bool {
	if len(args) < 1 {
		fmt.Println("Error: Need at least 1 input path or URL")
		return false
	}

	for _, arg := range args {
		_, err := url.Parse(arg)
		if err != nil {
			fmt.Printf("Error: Invalid URL '%s'\n", arg)
			continue
		}

		extension := path.Ext(arg)
		if extension != ".jpg" && extension != ".jpeg" && extension != ".png" && extension != ".webp" && extension != ".gif" {
			fmt.Printf("Error: Unsupported file extension '%s' for '%s'\n", extension, arg)
			return false
		}
	}

	return true
}
