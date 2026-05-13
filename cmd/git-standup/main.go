package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func usage() {
	fmt.Print(`git-standup — What did you do yesterday?

Usage: git-standup [options]

Options:
  -d int         Number of days back (default 1)
  -a string      Author name filter
  -b string      Branch (default current)
  -r             Include remote branches
  -s             Short format (one line per commit)
  -j             JSON output (agent-friendly)
  -q             Quiet (no headers, just data)

Examples:
  git-standup              # Yesterday's commits
  git-standup -d 3         # Last 3 days
  git-standup -a "Bark"    # Only commits by author
  git-standup -j           # JSON output for agents
  git-standup -jq          # JSON output, no stderr
`)
}

type Commit struct {
	Hash    string `json:"hash"`
	Date    string `json:"date"`
	Author  string `json:"author"`
	Message string `json:"message"`
}

type StandupReport struct {
	Days       int       `json:"days"`
	Total      int       `json:"total"`
	Authors    []string  `json:"authors"`
	Commits    []Commit  `json:"commits"`
	Repo       string    `json:"repo"`
	Generated  string    `json:"generated"`
}

func main() {
	var (
		days   = flag.Int("d", 1, "Number of days back")
		author = flag.String("a", "", "Author name filter")
		branch = flag.String("b", "", "Branch (default current)")
		remote = flag.Bool("r", false, "Include remote branches")
		short  = flag.Bool("s", false, "Short format")
		jsonOut= flag.Bool("j", false, "JSON output")
		quiet  = flag.Bool("q", false, "Quiet mode")
	)
	flag.Usage = usage
	flag.Parse()

	// Env var overrides
	if envDays := os.Getenv("STANDUP_DAYS"); envDays != "" && *days == 1 {
		fmt.Sscanf(envDays, "%d", days)
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -*days)

	// Get repo name
	repo := getRepoName()

	// Build git log command
	args := []string{
		"log",
		"--since", startDate.Format("2006-01-02"),
		"--until", endDate.Format("2006-01-02"),
		"--date=short",
		"--format=%H|%ad|%an|%s",
	}

	if *author != "" {
		args = append(args, "--author", *author)
	}

	if *branch != "" {
		args = append(args, *branch)
	} else if !*remote {
		args = append(args, "HEAD")
	}

	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		if !*quiet {
			fmt.Fprintf(os.Stderr, "Error running git log: %v\n", err)
		}
		os.Exit(1)
	}

	// Parse commits
	var commits []Commit
	authorSet := make(map[string]bool)
	
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 4)
		if len(parts) == 4 {
			c := Commit{
				Hash:    parts[0][:7],
				Date:    parts[1],
				Author:  parts[2],
				Message: parts[3],
			}
			commits = append(commits, c)
			authorSet[parts[2]] = true
		}
	}

	// Build author list
	var authors []string
	for a := range authorSet {
		authors = append(authors, a)
	}

	if *jsonOut {
		report := StandupReport{
			Days:      *days,
			Total:     len(commits),
			Authors:   authors,
			Commits:   commits,
			Repo:      repo,
			Generated: time.Now().Format(time.RFC3339),
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(report)
		return
	}

	if len(commits) == 0 {
		if !*quiet {
			fmt.Printf("No commits in the last %d day(s)\n", *days)
		}
		os.Exit(0)
	}

	if !*quiet {
		fmt.Printf("Commits in the last %d day(s):\n\n", *days)
	}

	for _, c := range commits {
		if *short {
			fmt.Printf("%s %s\n", c.Hash, c.Message)
		} else {
			fmt.Printf("%s | %s | %s | %s\n", c.Hash, c.Date, c.Author, c.Message)
		}
	}

	if !*quiet {
		fmt.Printf("\nTotal: %d commit(s)\n", len(commits))
		if len(authors) > 1 {
			fmt.Printf("Authors: %s\n", strings.Join(authors, ", "))
		}
	}
}

func getRepoName() string {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		out, _ = exec.Command("git", "rev-parse", "--show-toplevel").Output()
		return strings.TrimSpace(string(out))
	}
	url := strings.TrimSpace(string(out))
	// Extract repo name from URL
	if i := strings.LastIndex(url, "/"); i != -1 {
		name := url[i+1:]
		name = strings.TrimSuffix(name, ".git")
		return name
	}
	return url
}
