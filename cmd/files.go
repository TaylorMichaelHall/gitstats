package cmd

import (
	"fmt"
	"gitstats/internals/gitutils"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

type FileChange struct {
	Name        string
	ChangeCount int
}

func GetFilesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "files [path_to_repo]",
		Short: "Show file change frequency",
		Long:  `Display a chart of the most frequently changed files in the repository.`,
		Args:  cobra.ExactArgs(1),
		Run:   runFiles,
	}

	cmd.Flags().StringSliceP("ignore", "i", []string{}, "Substrings to ignore in file paths (can be used multiple times)")
	return cmd
}

func runFiles(cmd *cobra.Command, args []string) {
	repoPath := args[0]
	ignoreSubstrings, _ := cmd.Flags().GetStringSlice("ignore")

	fileChanges, err := gitutils.GetFileChanges(repoPath, ignoreSubstrings)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	sortedFileChanges := sortFileChanges(fileChanges)
	printFileChangeChart(sortedFileChanges)
}

func sortFileChanges(fileChanges map[string]int) []FileChange {
	var changes []FileChange
	for name, count := range fileChanges {
		changes = append(changes, FileChange{Name: name, ChangeCount: count})
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].ChangeCount > changes[j].ChangeCount
	})

	return changes
}

func printFileChangeChart(fileChanges []FileChange) {
	fmt.Println("File Change Frequency Chart (Top 25):")
	fmt.Println("-------------------------------------")

	maxCount := 0
	maxNameLength := 0
	if len(fileChanges) > 0 {
		maxCount = fileChanges[0].ChangeCount
	}
	for _, fc := range fileChanges[:min(25, len(fileChanges))] {
		if len(fc.Name) > maxNameLength {
			maxNameLength = len(fc.Name)
		}
	}

	colors := []string{
		"\033[31m", "\033[32m", "\033[33m", "\033[34m", "\033[35m",
		"\033[36m", "\033[91m", "\033[92m", "\033[93m", "\033[94m",
	}

	for i, fc := range fileChanges[:min(25, len(fileChanges))] {
		barLength := 0
		if maxCount > 0 {
			barLength = (fc.ChangeCount * 50) / maxCount
		}
		colorIndex := i % len(colors)

		fmt.Printf("%2d. %-*s |", i+1, maxNameLength, fc.Name)
		fmt.Print(colors[colorIndex])
		for j := 0; j < barLength; j++ {
			fmt.Print("█")
		}
		fmt.Print("\033[0m") // Reset color
		fmt.Printf(" (%d)\n", fc.ChangeCount)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
