package main

import (
	"fmt"

	"github.com/maisieccino/proton-filters/internal/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Example: "proton-filters [command]",
	Run:     cmd.Root,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
