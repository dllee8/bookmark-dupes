# bookmarkdupes

Answers one question: which of my browser bookmarks are duplicates?

After a few years of "save it for later" you end up with the same article
saved three times under three different folders, once with `www.` and once
without, once with a trailing slash. Browsers don't tell you this. This is
a small command-line tool that reads a bookmark export and lists the
duplicates, grouped, with the folder each copy lives in.

## Usage

Export your bookmarks to HTML first:

- Chrome/Edge: `chrome://bookmarks` → menu → "Export bookmarks"
- Firefox: Bookmarks → Manage Bookmarks → Import and Backup → "Export Bookmarks to HTML"
- Safari: File → Export Bookmarks

Then run:

```
go run . ~/Downloads/bookmarks.html
```

Example output:

```
https://example.com/pricing (3 copies)
  - "Example — Pricing" in Work/Reference
  - "pricing page (check later)" in Inbox
  - "Example pricing" in Personal
no duplicate bookmarks found
```

(The second line only prints if there are none — the tool prints one or
the other, not both.)

## How it decides two bookmarks are duplicates

Two bookmarks are considered the same page if, after normalizing the URL,
they're identical:

- scheme and host are lowercased
- a leading `www.` is dropped
- a trailing slash on the path is dropped
- the query string is kept as-is (`?id=1` and `?id=2` are different pages)

This lives in `bookmarks.NormalizeURL` and is deliberately conservative —
it won't catch a link that changed domains (e.g. a URL shortener) or one
with reordered query parameters.

## Design

The bookmark HTML parser (`bookmarks.ParseNetscapeHTML`) and the duplicate
finder (`bookmarks.FindDuplicates`) are both pure functions: given the same
input string or slice, they always return the same output, and neither
touches a file, the network, or the clock (except to convert a bookmark's
own `ADD_DATE` attribute). `main.go` is the only place that touches disk.
That split is what makes the core logic testable without fixture files on
disk — see `bookmarks/duplicates_test.go`.

## Status

Early skeleton. Parses the standard Netscape bookmark export format used
by Chrome, Firefox, and Safari, and finds exact-after-normalization
duplicates. See the issues/roadmap for what's next.

## License

MIT, see LICENSE.
