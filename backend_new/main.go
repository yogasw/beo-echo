package main

import (
	"log"
	"os"

	"mockoon-control-panel/backend_new/cmd"
)

func main() {
	// Execute the root command
	if err := cmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
