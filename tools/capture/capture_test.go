package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestExtractRegion(t *testing.T) {
	src := []byte(strings.Join([]string{
		"package x",
		"// region:demo:start",
		"    func f() int {",
		"        return 1",
		"    }",
		"// region:demo:end",
		"var y = 2",
	}, "\n"))
	got, err := extractRegion(src, "demo")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"func f() int {", "    return 1", "}"}
	if !equalStrings(got, want) {
		t.Fatalf("dedent wrong:\n got=%q\nwant=%q", got, want)
	}
}

func TestExtractRegionMissing(t *testing.T) {
	if _, err := extractRegion([]byte("nothing here"), "nope"); err == nil {
		t.Fatal("expected error for missing region")
	}
}

func TestApplyNormalizers(t *testing.T) {
	in := "ptr 0xc000123 goroutine 42 built with go1.26.5"
	got := applyNormalizers(in, []string{"addr", "goid", "gover"}, ".")
	want := "ptr 0xADDR goroutine N built with goX"
	if got != want {
		t.Fatalf("normalize:\n got=%q\nwant=%q", got, want)
	}
}

func TestJSONCStripsComments(t *testing.T) {
	in := []byte(`{
  // COMMENTLINE
  "k": "v // KEEPME",
  "n": 1 // TRAILINGCOMMENT
}`)
	r := newJSONCReader(in)
	b, _ := io.ReadAll(r)
	s := string(b)
	if strings.Contains(s, "COMMENTLINE") || strings.Contains(s, "TRAILINGCOMMENT") {
		t.Fatalf("comment not stripped: %q", s)
	}
	if !strings.Contains(s, "v // KEEPME") {
		t.Fatalf("in-string slashes wrongly stripped: %q", s)
	}
}

func TestProcessRewritesOutputBlock(t *testing.T) {
	// A file whose output block is stale gets rewritten to the command's stdout.
	dir := t.TempDir()
	writeFile(t, dir, "capture.json", `{
  "outputs": { "hi": { "cmd": ["printf", "hello\nworld"], "lang": "text" } }
}`)
	md := "# t\n\n<!-- BEGIN:output hi -->\nstale\n<!-- END:output hi -->\n"
	mc := newManifestCache()
	out, portableDiff, err := process(dir+"/README.md", []byte(md), mc, false)
	if err != nil {
		t.Fatal(err)
	}
	if !portableDiff {
		t.Fatal("expected a portable diff")
	}
	s := string(out)
	if !strings.Contains(s, "```text\nhello\nworld\n```") {
		t.Fatalf("block not rewritten:\n%s", s)
	}
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(dir+"/"+name, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
