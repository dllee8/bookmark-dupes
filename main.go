// Command bookmarkdupes reads a browser bookmark export and reports
// bookmarks that point at the same page, wherever they're filed.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"bookmarkdupes/bookmarks"
)

func main() {
	jsonOutput := flag.Bool("json", false, "print duplicate groups as JSON instead of plain text")
	folder := flag.String("folder", "", "only look for duplicates under this folder (and its subfolders)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: bookmarkdupes [--json] [--folder <path>] <exported-bookmarks.html>")
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	data, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reading file: %v\n", err)
		os.Exit(1)
	}

	marks, err := bookmarks.ParseNetscapeHTML(string(data))
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing bookmarks: %v\n", err)
		os.Exit(1)
	}

	marks = bookmarks.FilterByFolder(marks, *folder)

	dupes := bookmarks.FindDuplicates(marks)

	if *jsonOutput {
		// Encode as [] rather than null when there are no duplicates, so
		// consumers don't need a nil check before ranging over the result.
		if dupes == nil {
			dupes = []bookmarks.DuplicateGroup{}
		}
		if err := json.NewEncoder(os.Stdout).Encode(dupes); err != nil {
			fmt.Fprintf(os.Stderr, "encoding json: %v\n", err)
			os.Exit(1)
		}
		return
	}

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
