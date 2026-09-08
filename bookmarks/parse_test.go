package bookmarks

import (
	"testing"
	"time"
)

func TestParseNetscapeHTML(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []Bookmark
	}{
		{
			name: "nested folders build slash-joined paths",
			in: `<!DOCTYPE NETSCAPE-Bookmark-file-1>
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
    <DT><H3>Work</H3>
    <DL><p>
        <DT><A HREF="https://example.com/foo" ADD_DATE="1690000000">Foo</A>
        <DT><H3>Reading</H3>
        <DL><p>
            <DT><A HREF="https://example.com/bar" ADD_DATE="1690000001">Bar</A>
        </DL><p>
    </DL><p>
    <DT><A HREF="https://other.com" ADD_DATE="1690000002">Other</A>
</DL><p>
`,
			want: []Bookmark{
				{Title: "Foo", URL: "https://example.com/foo", Folder: "Work", Added: time.Unix(1690000000, 0).UTC()},
				{Title: "Bar", URL: "https://example.com/bar", Folder: "Work/Reading", Added: time.Unix(1690000001, 0).UTC()},
				{Title: "Other", URL: "https://other.com", Folder: "", Added: time.Unix(1690000002, 0).UTC()},
			},
		},
		{
			name: "html entities in titles and folder names are unescaped",
			in: `<DL><p>
    <DT><H3>R&amp;D</H3>
    <DL><p>
        <DT><A HREF="https://example.com">Tom &amp; Jerry&#39;s &quot;show&quot;</A>
    </DL><p>
</DL><p>
`,
			want: []Bookmark{
				{Title: `Tom & Jerry's "show"`, URL: "https://example.com", Folder: "R&D"},
			},
		},
		{
			name: "links with an empty href are skipped",
			in: `<DL><p>
    <DT><A HREF="">No href here</A>
    <DT><A HREF="https://example.com">Has href</A>
</DL><p>
`,
			want: []Bookmark{
				{Title: "Has href", URL: "https://example.com", Folder: ""},
			},
		},
		{
			name: "missing or malformed add_date parses as zero time",
			in: `<DL><p>
    <DT><A HREF="https://example.com/a">No date</A>
    <DT><A HREF="https://example.com/b" ADD_DATE="not-a-number">Bad date</A>
</DL><p>
`,
			want: []Bookmark{
				{Title: "No date", URL: "https://example.com/a", Folder: ""},
				{Title: "Bad date", URL: "https://example.com/b", Folder: ""},
			},
		},
		{
			name: "closing a folder returns to its parent, not the root",
			in: `<DL><p>
    <DT><H3>A</H3>
    <DL><p>
        <DT><H3>B</H3>
        <DL><p>
            <DT><A HREF="https://example.com/deep">Deep</A>
        </DL><p>
        <DT><A HREF="https://example.com/back-in-a">Back in A</A>
    </DL><p>
</DL><p>
`,
			want: []Bookmark{
				{Title: "Deep", URL: "https://example.com/deep", Folder: "A/B"},
				{Title: "Back in A", URL: "https://example.com/back-in-a", Folder: "A"},
			},
		},
		{
			name: "empty input yields no bookmarks",
			in:   "",
			want: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := ParseNetscapeHTML(c.in)
			if err != nil {
				t.Fatalf("ParseNetscapeHTML returned error: %v", err)
			}
			if len(got) != len(c.want) {
				t.Fatalf("got %d bookmarks, want %d\ngot:  %+v\nwant: %+v", len(got), len(c.want), got, c.want)
			}
			for i := range got {
				if got[i] != c.want[i] {
					t.Errorf("bookmark %d = %+v, want %+v", i, got[i], c.want[i])
				}
			}
		})
	}
}
