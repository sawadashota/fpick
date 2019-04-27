package main

import (
	"fmt"
	"os"

	"github.com/sawadashota/fpick/cmd"
)

func main() {
	if err := cmd.Cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
