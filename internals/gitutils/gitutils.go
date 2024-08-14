package gitutils

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type CommitInfo struct {
	Hash      string
	Author    string
	Timestamp time.Time
	Message   string
}

func GetGitLog(repoPath string) ([]CommitInfo, error) {
	cmd := exec.Command("git", "-C", repoPath, "log", "--pretty=format:%H|%an|%at|%s")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error accessing git repository: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	commits := make([]CommitInfo, 0, len(lines))

	for _, line := range lines {
		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue
		}

		timestamp, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			// Log the error and skip this commit
			fmt.Printf("Error parsing timestamp for commit %s: %v\n", parts[0], err)
			continue
		}
		commitTime := time.Unix(timestamp, 0)

		commits = append(commits, CommitInfo{
			Hash:      parts[0],
			Author:    parts[1],
			Timestamp: commitTime,
			Message:   parts[3],
		})
	}

	return commits, nil
}

func GetFileChanges(repoPath string, ignoreSubstrings []string) (map[string]int, error) {
	cmd := exec.Command("git", "-C", repoPath, "log", "--name-only", "--pretty=format:")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error executing git command: %v", err)
	}

	fileChanges := make(map[string]int)
	files := strings.Split(string(output), "\n")

	for _, file := range files {
		file = strings.TrimSpace(file)
		if file == "" {
			continue
		}

		if shouldIgnore(file, ignoreSubstrings) {
			continue
		}

		fileChanges[file]++
	}

	return fileChanges, nil
}

func shouldIgnore(file string, ignoreSubstrings []string) bool {
	for _, substr := range ignoreSubstrings {
		if strings.Contains(file, substr) {
			return true
		}
	}
	return false
}
