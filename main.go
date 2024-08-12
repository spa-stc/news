package main

import (
	"os"

	"stpaulacademy.tech/newsletter/cmd"
)

func main() {
	if err := cmd.RootCMD.Execute(); err != nil {
		os.Exit(1)
	}
}
