package dataparser

import (
	"fmt"
	"os/exec"
	"strings"
)

var (
	authors           map[string]*lineCommitsFiles
	allFiles          []string
	commits           map[string]map[string]bool
	authorOrCommitter = "author "
)

func getAllFiles() ([]string, error) {
	cmd := exec.Command("git", "-C", info.repository, "ls-tree", "--name-only", "-r", info.revision)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running git ls-tree: %v", err)
	}
	files := strings.Split(string(output), "\n")
	files = files[0 : len(files)-1]
	var matchedFiles []string
	for _, file := range files {
		if info.isFileFine(file) {
			matchedFiles = append(matchedFiles, file)
		}
	}
	return matchedFiles, nil
}

func processGitBlame() error {
	for _, file := range allFiles {
		cmd := exec.Command("git", "-C", info.repository, "blame", "--line-porcelain", info.revision, file)
		change := 1
		if authorOrCommitter == "committer " {
			change = 5
		}
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("error running git blame: %v", err)
		}

		lines := strings.Split(string(output), "\n")
		if string(output) == "" {
			err := processEmptyFile(file)
			if err != nil {
				return fmt.Errorf("error processing empty file: %v", err)
			}
			continue
		}
		hasAuthorChanged := make(map[string]bool)

		for i, line := range lines {
			if strings.HasPrefix(line, authorOrCommitter) {
				author := strings.TrimPrefix(line, authorOrCommitter)
				if _, ok := authors[author]; !ok {
					authors[author] = &lineCommitsFiles{0, 0, 0}
				}
				if _, ok := hasAuthorChanged[author]; !ok {
					hasAuthorChanged[author] = true
					authors[author].fileCount++
				}
				if _, ok := commits[author]; !ok {
					commits[author] = make(map[string]bool, 0)
				}
				authors[author].lines++
				commit := strings.Fields(lines[i-change])[0]
				commits[author][commit] = true
			}
		}
	}
	for author := range commits {
		authors[author].commit = len(commits[author])
	}

	return nil
}

func processEmptyFile(file string) error {
	cmd := exec.Command("git", "-C", info.repository, "log", "-n1", "--format=%H%n%an", info.revision, "--", file)
	if info.usecommitter {
		cmd = exec.Command("git", "-C", info.repository, "log", "-n1", "--format=%H%n%cn", info.revision, "--", file)
	}
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error running git shortlog: %v", err)
	}
	lines := strings.Split(string(output), "\n")
	author := lines[1]
	commit := lines[0]
	if _, ok := authors[author]; !ok {
		authors[author] = &lineCommitsFiles{0, 0, 0}
	}
	if _, ok := commits[author]; !ok {
		commits[author] = make(map[string]bool)
	}
	commits[author][commit] = true
	authors[author].fileCount++
	return nil
}

func processFileCount() error {
	return nil
}

func Init(
	repository string,
	revision string,
	orderby string,
	usecommitter bool,
	format string,
	extensions string,
	languages string,
	exclude string,
	restrictto string,
) error {
	var err error
	err = initParseInfo(
		repository,
		revision,
		orderby,
		usecommitter,
		format,
		extensions,
		languages,
		exclude,
		restrictto,
	)
	fmt.Printf("%v %v %v %v %v %v %v %v %v\n", repository, revision, orderby, usecommitter, format, extensions, languages, exclude, restrictto)
	if err != nil {
		return err
	}
	if usecommitter {
		authorOrCommitter = "committer "
	}
	authors = make(map[string]*lineCommitsFiles)
	commits = make(map[string]map[string]bool)
	allFiles, err = getAllFiles()
	if err != nil {
		return err
	}
	err = processGitBlame()
	if err != nil {
		return err
	}
	err = processFileCount()
	if err != nil {
		return err
	}
	return nil
}
