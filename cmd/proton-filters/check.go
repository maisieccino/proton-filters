package main

import (
	"github.com/maisieccino/proton-filters/internal/cmd"
	"github.com/spf13/cobra"
)

var checkCmd = cobra.Command{
	Use:     "check [file]",
	Short:   "Check a sieve filter for correctness",
	Long:    "",
	Example: "proton-filters check test.siv",
	RunE:    cmd.Check,
	Args:    cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(&checkCmd)
}
