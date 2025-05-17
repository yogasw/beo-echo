package main

import (
	"log"
	"os"

	"beo-echo/backend/cmd"
)

func main() {
	// Execute the root command
	if err := cmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
