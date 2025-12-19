package main

import "github.com/spf13/cobra"

var checkCmd = cobra.Command{
	Use:               "check [file]",
	Short:             "Check a sieve filter for correctness",
	Long:              "",
	Example:           "proton-filters check test.siv",
	ValidArgs:         []cobra.Completion{},
	ValidArgsFunction: nil,
}

func init() {
	rootCmd.AddCommand(&checkCmd)
}
