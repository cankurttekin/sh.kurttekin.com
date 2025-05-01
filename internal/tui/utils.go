package tui

import (
	"regexp"
)

// FindLinks extracts all URLs from a slice of content strings
func FindLinks(content []string) []string {
	links := []string{}
	re := regexp.MustCompile(`https?://\S+`)

	for _, line := range content {
		matches := re.FindAllString(line, -1)
		links = append(links, matches...)
	}

	return links
}
