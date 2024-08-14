package cmd

import (
	"bufio"
	"fmt"
	"gitstats/internals/gitutils"
	"os"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

type Contributor struct {
	Name         string
	Commits      int
	FirstCommit  time.Time
	LatestCommit time.Time
}

func GetContributorsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "contributors [path_to_repo]",
		Short: "Show contributor statistics",
		Long:  `Display a bar graph of contributor statistics and allow detailed view of individual contributors.`,
		Args:  cobra.ExactArgs(1),
		Run:   runContributors,
	}
}

func runContributors(cmd *cobra.Command, args []string) {
	repoPath := args[0]
	commits, err := gitutils.GetGitLog(repoPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	contributors := processCommits(commits)
	printBarGraph(contributors)

	fmt.Println("\nEnter the number of a contributor to see more details, or 'q' to quit:")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if input == "q" {
			break
		}
		var index int
		_, err := fmt.Sscanf(input, "%d", &index)
		if err != nil || index < 1 || index > len(contributors) {
			fmt.Println("Invalid input. Please enter a number between 1 and", len(contributors))
			continue
		}
		displayContributorDetails(contributors[index-1])
		fmt.Println("\nEnter another number or 'q' to quit:")
	}
}

func processCommits(commits []gitutils.CommitInfo) []Contributor {
	contributorMap := make(map[string]*Contributor)

	for _, commit := range commits {
		if contributor, exists := contributorMap[commit.Author]; exists {
			contributor.Commits++
			if commit.Timestamp.Before(contributor.FirstCommit) {
				contributor.FirstCommit = commit.Timestamp
			}
			if commit.Timestamp.After(contributor.LatestCommit) {
				contributor.LatestCommit = commit.Timestamp
			}
		} else {
			contributorMap[commit.Author] = &Contributor{
				Name:         commit.Author,
				Commits:      1,
				FirstCommit:  commit.Timestamp,
				LatestCommit: commit.Timestamp,
			}
		}
	}

	contributors := make([]Contributor, 0, len(contributorMap))
	for _, c := range contributorMap {
		contributors = append(contributors, *c)
	}

	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Commits > contributors[j].Commits
	})

	return contributors
}

func printColoredBar(length int, color string) {
	barChar := "█"
	fmt.Print(color)
	for i := 0; i < length; i++ {
		fmt.Print(barChar)
	}
	fmt.Print("\033[0m")
}

func printBarGraph(contributors []Contributor) {
	maxCommits := 0
	maxNameLength := 0
	for _, c := range contributors {
		if c.Commits > maxCommits {
			maxCommits = c.Commits
		}
		if len(c.Name) > maxNameLength {
			maxNameLength = len(c.Name)
		}
	}

	colors := []string{"\033[31m", "\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[91m", "\033[92m", "\033[93m", "\033[94m"}

	for i, c := range contributors {
		barLength := int((float64(c.Commits) / float64(maxCommits)) * 50)
		fmt.Printf("%2d. %-*s | ", i+1, maxNameLength, c.Name)
		printColoredBar(barLength, colors[i%len(colors)])
		fmt.Printf(" (%d commits)\n", c.Commits)
	}
}

func displayContributorDetails(contributor Contributor) {
	fmt.Printf("\nDetails for %s:\n", contributor.Name)
	fmt.Printf("Total Commits: %d\n", contributor.Commits)
	fmt.Printf("First Commit: %s\n", contributor.FirstCommit.Format("2006-01-02 15:04:05"))
	fmt.Printf("Latest Commit: %s\n", contributor.LatestCommit.Format("2006-01-02 15:04:05"))
}
