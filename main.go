// Command bookmarkdupes reads a browser bookmark export and reports
// bookmarks that point at the same page, wherever they're filed.
package main

import (
	"fmt"
	"os"

	"bookmarkdupes/bookmarks"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: bookmarkdupes <exported-bookmarks.html>")
		os.Exit(1)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "reading file: %v\n", err)
		os.Exit(1)
	}

	marks, err := bookmarks.ParseNetscapeHTML(string(data))
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing bookmarks: %v\n", err)
		os.Exit(1)
	}

	dupes := bookmarks.FindDuplicates(marks)
	if len(dupes) == 0 {
		fmt.Println("no duplicate bookmarks found")
		return
	}

	for _, group := range dupes {
		fmt.Printf("%s (%d copies)\n", group.NormalizedURL, len(group.Bookmarks))
		for _, b := range group.Bookmarks {
			folder := b.Folder
			if folder == "" {
				folder = "(no folder)"
			}
			fmt.Printf("  - %q in %s\n", b.Title, folder)
		}
	}
}
