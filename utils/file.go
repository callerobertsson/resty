package utils

import (
	"errors"
	"os"
	"regexp"
	"strings"
)

var reHttpFile = regexp.MustCompile("\\.http$")

// FileExists checks if file path f exists
func FileExists(f string) bool {
	_, err := os.Stat(f)
	return !errors.Is(err, os.ErrNotExist)
}

// GetHttpFile returns file paths for all .http file in directory path d with a recursive option
func GetHttpFilePaths(d string, recure bool) []string {
	fs, _ := getMatchingFiles(d, reHttpFile, recure)
	return fs
}

// Get matching entries in dir, and if recursive, all subdirs
// Returns list of matches and the total count of files in tree
func getMatchingFiles(dir string, r *regexp.Regexp, recursive bool) ([]string, int) {

	fis, err := os.ReadDir(dir)
	if err != nil {
		return []string{}, 0
	}

	entries := []string{}
	dirs := []string{}
	total := 0

	for _, fi := range fis {
		if strings.Index(fi.Name(), ".") != 0 && fi.IsDir() {
			dirs = append(dirs, dir+"/"+fi.Name())
			continue
		}
		total++
		if r.MatchString(fi.Name()) {
			entries = append(entries, dir+"/"+fi.Name())
		}
	}

	if recursive {
		for _, d := range dirs {
			newEntries, newTotal := getMatchingFiles(d, r, recursive)
			total += newTotal
			entries = append(entries, newEntries...)
		}
	}

	return entries, total
}
