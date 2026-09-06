// Package bookmarks reads browser bookmark exports and answers questions
// about them. Every exported function is pure: given the same input it
// returns the same output, with no file or network access, so the parsing
// and analysis logic can be tested without touching disk.
package bookmarks

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Bookmark is a single entry pulled from a bookmark export.
type Bookmark struct {
	Title  string
	URL    string
	Folder string // slash-joined path, e.g. "Work/Reading"
	Added  time.Time
}

var (
	// Chrome, Firefox, and Safari all export the same "Netscape Bookmark
	// File Format", but it's tag soup rather than well-formed HTML (tags
	// are often left unclosed), so html/xml parsers choke on it. Line-based
	// regexes are what every browser's own importer does too.
	folderOpenRe = regexp.MustCompile(`(?i)<DT><H3[^>]*>(.*?)</H3>`)
	linkRe       = regexp.MustCompile(`(?i)<DT><A\s+([^>]*)>(.*?)</A>`)
	attrRe       = regexp.MustCompile(`(\w+)="([^"]*)"`)
	folderCloser = regexp.MustCompile(`(?i)^</DL>`)
)

// ParseNetscapeHTML parses the contents of a browser bookmark export
// (Chrome, Firefox, and Safari all use this format) into a flat list of
// bookmarks. Folder nesting is preserved as a slash-joined path on each
// Bookmark.
func ParseNetscapeHTML(content string) ([]Bookmark, error) {
	var result []Bookmark
	var folderStack []string

	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if m := folderOpenRe.FindStringSubmatch(line); m != nil {
			folderStack = append(folderStack, unescapeHTML(m[1]))
			continue
		}

		if folderCloser.MatchString(line) {
			if len(folderStack) > 0 {
				folderStack = folderStack[:len(folderStack)-1]
			}
			continue
		}

		if m := linkRe.FindStringSubmatch(line); m != nil {
			attrs := parseAttrs(m[1])
			href := attrs["href"]
			if href == "" {
				continue
			}
			result = append(result, Bookmark{
				Title:  unescapeHTML(m[2]),
				URL:    href,
				Folder: strings.Join(folderStack, "/"),
				Added:  parseAddDate(attrs["add_date"]),
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning bookmark export: %w", err)
	}

	return result, nil
}

func parseAttrs(raw string) map[string]string {
	attrs := make(map[string]string)
	for _, m := range attrRe.FindAllStringSubmatch(raw, -1) {
		attrs[strings.ToLower(m[1])] = m[2]
	}
	return attrs
}

func parseAddDate(v string) time.Time {
	if v == "" {
		return time.Time{}
	}
	sec, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return time.Time{}
	}
	return time.Unix(sec, 0).UTC()
}

func unescapeHTML(s string) string {
	replacer := strings.NewReplacer(
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
		"&quot;", `"`,
		"&#39;", "'",
	)
	return replacer.Replace(s)
}
