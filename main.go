package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run script.go <directory>")
		os.Exit(1)
	}

	root := os.Args[1]
	var data [][]string

	// Walk through the directory
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %q: %v\n", path, err)
			return nil
		}

		// Check if the current directory contains a .git folder
		if info.IsDir() && strings.HasSuffix(path, ".git") {
			repoDir := filepath.Dir(path)
			dirName := filepath.Base(repoDir)
			status := getGitStatus(repoDir)
			data = append(data, []string{dirName, repoDir, status})
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the directory: %v\n", err)
		os.Exit(1)
	}

	// Print results in a table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Directory Name", "Path", "Git Status"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

// getGitStatus executes `git status --short` in the given repository directory
func getGitStatus(repoDir string) string {
	cmd := exec.Command("git", "-C", repoDir, "status", "--short")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("Error: %s", stderr.String())
	}

	return strings.TrimSpace(out.String())
}
