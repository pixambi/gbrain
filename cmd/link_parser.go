package cmd

import (
	"regexp"
	"strings"
)

var linkPattern = regexp.MustCompile(`\[\[([^\[\]]+)\]\]`)

type Link struct {
	Title    string
	Position [2]int
}

func parseLinks(content string) []Link {
	var links []Link

	matches := linkPattern.FindAllStringSubmatchIndex(content, -1)
	for _, match := range matches {
		if len(match) >= 4 {
			start := match[2]
			end := match[3]
			title := content[start:end]
			links = append(links, Link{
				Title:    title,
				Position: [2]int{match[0], match[1]},
			})
		}
	}

	return links
}

func renderContent(content string, links []Link, currentLinkIndex int) string {
	if len(links) == 0 {
		return content
	}

	var result strings.Builder
	lastEnd := 0

	for i, link := range links {
		start, end := link.Position[0], link.Position[1]

		result.WriteString(content[lastEnd:start])

		linkTitle := link.Title

		if i == currentLinkIndex {
			result.WriteString(selectedLinkStyle.Render(linkTitle))
		} else {
			result.WriteString(linkStyle.Render(linkTitle))
		}

		lastEnd = end
	}

	result.WriteString(content[lastEnd:])
	return result.String()
}
