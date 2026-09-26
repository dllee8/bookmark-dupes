package bookmarks

import "testing"

func TestNormalizeURL(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"strips www and trailing slash", "https://www.Example.com/foo/", "https://example.com/foo"},
		{"already normalized", "https://example.com/foo", "https://example.com/foo"},
		{"keeps query string", "http://example.com/foo?a=1", "http://example.com/foo?a=1"},
		{"root path collapses", "https://example.com/", "https://example.com"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := NormalizeURL(c.in)
			if got != c.want {
				t.Errorf("NormalizeURL(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestFindDuplicates(t *testing.T) {
	input := []Bookmark{
		{Title: "Example", URL: "https://example.com/", Folder: "Work"},
		{Title: "Example mirror", URL: "https://www.example.com", Folder: "Personal"},
		{Title: "Other", URL: "https://other.com", Folder: "Work"},
	}

	got := FindDuplicates(input)

	if len(got) != 1 {
		t.Fatalf("expected 1 duplicate group, got %d", len(got))
	}
	if len(got[0].Bookmarks) != 2 {
		t.Fatalf("expected 2 bookmarks in duplicate group, got %d", len(got[0].Bookmarks))
	}
}

func TestFilterByFolder(t *testing.T) {
	input := []Bookmark{
		{Title: "Root", URL: "https://a.com", Folder: ""},
		{Title: "Work", URL: "https://b.com", Folder: "Work"},
		{Title: "Work sub", URL: "https://c.com", Folder: "Work/Reading"},
		{Title: "Different case", URL: "https://d.com", Folder: "WORK/Reading"},
		{Title: "Other", URL: "https://e.com", Folder: "Personal"},
		{Title: "Prefix but not subfolder", URL: "https://f.com", Folder: "Workshop"},
	}

	got := FilterByFolder(input, "Work")

	want := []string{"https://b.com", "https://c.com", "https://d.com"}
	if len(got) != len(want) {
		t.Fatalf("expected %d bookmarks, got %d: %+v", len(want), len(got), got)
	}
	for i, b := range got {
		if b.URL != want[i] {
			t.Errorf("bookmark %d = %q, want %q", i, b.URL, want[i])
		}
	}
}

func TestFilterByFolderEmptyReturnsAll(t *testing.T) {
	input := []Bookmark{
		{Title: "A", URL: "https://a.com", Folder: "Work"},
		{Title: "B", URL: "https://b.com", Folder: ""},
	}

	got := FilterByFolder(input, "")

	if len(got) != len(input) {
		t.Fatalf("expected %d bookmarks, got %d", len(input), len(got))
	}
}

func TestFindDuplicatesNoDuplicates(t *testing.T) {
	input := []Bookmark{
		{Title: "A", URL: "https://a.com"},
		{Title: "B", URL: "https://b.com"},
	}

	if got := FindDuplicates(input); len(got) != 0 {
		t.Fatalf("expected no duplicate groups, got %d", len(got))
	}
}
