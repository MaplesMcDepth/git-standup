package main

import (
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

Examples:
  git-standup              # Yesterday's commits
  git-standup -d 3         # Last 3 days
  git-standup -a "Bark"    # Only commits by Bark
  git-standup -s           # Short format
`)
}

func main() {
	var (
		days   = flag.Int("d", 1, "Number of days back")
		author = flag.String("a", "", "Author name filter")
		branch = flag.String("b", "", "Branch (default current)")
		remote = flag.Bool("r", false, "Include remote branches")
		short  = flag.Bool("s", false, "Short format")
	)
	flag.Usage = usage
	flag.Parse()

	// Calculate date range
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -*days)

	// Build git log command
	args := []string{
		"log",
		"--since", startDate.Format("2006-01-02"),
		"--until", endDate.Format("2006-01-02"),
		"--date=short",
	}

	if *author != "" {
		args = append(args, "--author", *author)
	}

	if *branch != "" {
		args = append(args, *branch)
	} else if !*remote {
		args = append(args, "HEAD")
	}

	if *short {
		args = append(args, "--format=%h %s")
	} else {
		args = append(args, "--format=%h | %ad | %an | %s")
	}

	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running git log: %v\n", err)
		os.Exit(1)
	}

	output := strings.TrimSpace(string(out))
	if output == "" {
		fmt.Printf("No commits in the last %d day(s)\n", *days)
		return
	}

	fmt.Printf("Commits in the last %d day(s):\n\n", *days)
	fmt.Println(output)

	// Show stats
	lines := strings.Split(output, "\n")
	fmt.Printf("\nTotal: %d commit(s)\n", len(lines))
}
