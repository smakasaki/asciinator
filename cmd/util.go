package cmd

import (
	"fmt"
	"path"
)

// return false if error occurs
func checkArgsAndFlags(args []string) bool {
	if len(args) < 1 {
		fmt.Println("Error: Need at least 1 input path")
		return false
	}

	for _, arg := range args {
		extenstion := path.Ext(arg)
		if extenstion != ".jpg" && extenstion != ".jpeg" && extenstion != ".png" {
			return false
		}
	}

	return true
}
