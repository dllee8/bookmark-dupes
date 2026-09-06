package bookmarks

import (
	"net/url"
	"sort"
	"strings"
)

// DuplicateGroup is a set of bookmarks that point at the same page under
// NormalizeURL's rules.
type DuplicateGroup struct {
	NormalizedURL string
	Bookmarks     []Bookmark
}

// NormalizeURL reduces a URL to a form suitable for comparing bookmarks
// that point at the same page but were saved slightly differently: it
// lowercases the scheme and host, drops a leading "www.", and trims a
// trailing slash from the path. It deliberately keeps the query string,
// since "?id=1" and "?id=2" are usually different pages.
func NormalizeURL(raw string) string {
	trimmed := strings.TrimSpace(raw)

	u, err := url.Parse(trimmed)
	if err != nil || u.Host == "" {
		// Not a parseable absolute URL; fall back to a light-touch
		// normalization so obviously-identical strings still match.
		return strings.ToLower(strings.TrimRight(trimmed, "/"))
	}

	scheme := strings.ToLower(u.Scheme)
	host := strings.ToLower(strings.TrimPrefix(u.Host, "www."))
	path := strings.TrimRight(u.Path, "/")

	normalized := scheme + "://" + host + path
	if u.RawQuery != "" {
		normalized += "?" + u.RawQuery
	}
	return normalized
}

// FindDuplicates groups bookmarks that resolve to the same NormalizeURL
// value and returns only the groups with more than one member, ordered
// from most-duplicated to least.
func FindDuplicates(bookmarks []Bookmark) []DuplicateGroup {
	groups := make(map[string][]Bookmark)
	var order []string

	for _, b := range bookmarks {
		key := NormalizeURL(b.URL)
		if _, seen := groups[key]; !seen {
			order = append(order, key)
		}
		groups[key] = append(groups[key], b)
	}

	var result []DuplicateGroup
	for _, key := range order {
		if len(groups[key]) > 1 {
			result = append(result, DuplicateGroup{
				NormalizedURL: key,
				Bookmarks:     groups[key],
			})
		}
	}

	sort.SliceStable(result, func(i, j int) bool {
		return len(result[i].Bookmarks) > len(result[j].Bookmarks)
	})

	return result
}
