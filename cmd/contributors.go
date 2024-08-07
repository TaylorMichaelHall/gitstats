package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

type Contributor struct {
	Name         string
	Commits      int
	FirstCommit  time.Time
	LatestCommit time.Time
	LinesAdded   int
	LinesRemoved int
}

type CommitInfo struct {
	Name      string
	Timestamp time.Time
	Added     int
	Removed   int
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
	logEntries, err := getGitLog(repoPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	contributors := processGitLog(logEntries)
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

func getGitLog(repoPath string) ([]string, error) {
	cmd := exec.Command("git", "-C", repoPath, "log", "--format=%an|%at", "--numstat")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error accessing git repository: %v", err)
	}
	return strings.Split(string(output), "\n"), nil
}

func processLogEntry(entry string, commitChan chan<- CommitInfo) {
	if strings.Contains(entry, "|") {
		parts := strings.Split(entry, "|")
		name, timestampStr := parts[0], parts[1]
		timestamp, err := strconv.ParseInt(strings.TrimSpace(timestampStr), 10, 64)
		if err != nil {
			fmt.Printf("Error parsing timestamp for %s: %v\n", name, err)
			return
		}
		commitTime := time.Unix(timestamp, 0)
		commitChan <- CommitInfo{Name: name, Timestamp: commitTime}
	} else if len(entry) > 0 {
		parts := strings.Fields(entry)
		if len(parts) == 3 {
			added, _ := strconv.Atoi(parts[0])
			removed, _ := strconv.Atoi(parts[1])
			commitChan <- CommitInfo{Added: added, Removed: removed}
		}
	}
}

func processGitLog(logEntries []string) []Contributor {
	numWorkers := runtime.NumCPU()
	commitChan := make(chan CommitInfo, len(logEntries))
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, entry := range logEntries {
				processLogEntry(entry, commitChan)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(commitChan)
	}()

	contributors := make(map[string]*Contributor)
	var currentContributor *Contributor

	for commit := range commitChan {
		if commit.Name != "" {
			if contributor, exists := contributors[commit.Name]; exists {
				currentContributor = contributor
				currentContributor.Commits++
				if commit.Timestamp.Before(contributor.FirstCommit) {
					contributor.FirstCommit = commit.Timestamp
				}
				if commit.Timestamp.After(contributor.LatestCommit) {
					contributor.LatestCommit = commit.Timestamp
				}
			} else {
				currentContributor = &Contributor{
					Name:         commit.Name,
					Commits:      1,
					FirstCommit:  commit.Timestamp,
					LatestCommit: commit.Timestamp,
				}
				contributors[commit.Name] = currentContributor
			}
		} else if currentContributor != nil {
			currentContributor.LinesAdded += commit.Added
			currentContributor.LinesRemoved += commit.Removed
		}
	}

	result := make([]Contributor, 0, len(contributors))
	for _, c := range contributors {
		result = append(result, *c)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Commits > result[j].Commits
	})

	return result
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
	fmt.Printf("Lines Added: %d\n", contributor.LinesAdded)
	fmt.Printf("Lines Removed: %d\n", contributor.LinesRemoved)
	fmt.Printf("Net Lines: %d\n", contributor.LinesAdded-contributor.LinesRemoved)
}