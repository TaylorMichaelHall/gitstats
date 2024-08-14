package main

import (
	"fmt"
	"gitstats/cmd"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gitstats",
		Short: "A tool for analyzing Git repository statistics",
		Long:  `gitstats is a CLI tool that provides various statistics about a Git repository.`,
	}

	rootCmd.AddCommand(cmd.GetContributorsCmd(), cmd.GetFilesCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
