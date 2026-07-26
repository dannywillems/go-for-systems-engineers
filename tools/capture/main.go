// Command capture is the falsifiability engine of this repository.
//
// Every factual claim in a README must be backed by a program in the repo whose
// output is injected into the README between markers. capture (re)generates
// those blocks so no number, timing, size, or memory behavior is ever hand-typed.
//
// A markdown file may contain three kinds of generated block:
//
//		<!-- BEGIN:output NAME -->   ... <!-- END:output NAME -->
//		<!-- BEGIN:snippet NAME -->  ... <!-- END:snippet NAME -->
//		<!-- BEGIN:file NAME -->     ... <!-- END:file NAME -->
//
//	  - output  runs a command and injects its stdout (or combined output).
//	  - snippet extracts a region of a real source file (never hand-copied code).
//	  - file    injects the full contents of a committed file (e.g. bench results).
//
// The command / source / file behind each NAME is declared in the nearest
// ancestor capture.json. Run with -check to verify freshness without writing
// (the docs-fresh CI gate); run without it to regenerate in place.
//
// Blocks marked "portable": false are regenerated locally but SKIPPED by -check,
// because they capture architecture-specific output (assembly, escape analysis
// on a specific target) that legitimately differs between the author's machine
// and a CI runner. See README for the honest limits of this.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	beginRe = regexp.MustCompile(`^<!--\s*BEGIN:(output|snippet|file)\s+([A-Za-z0-9._-]+)\s*-->\s*$`)
	endRe   = regexp.MustCompile(`^<!--\s*END:(output|snippet|file)\s+([A-Za-z0-9._-]+)\s*-->\s*$`)
)

func main() {
	root := flag.String("root", ".", "directory tree to scan for markdown files")
	check := flag.Bool("check", false, "verify freshness without writing; exit 1 if any portable block would change")
	verbose := flag.Bool("v", false, "log every file processed")
	flag.Parse()

	mds, err := findMarkdown(*root)
	if err != nil {
		fatal(err)
	}

	changed := 0
	mc := newManifestCache()
	for _, md := range mds {
		orig, err := os.ReadFile(md)
		if err != nil {
			fatal(err)
		}
		out, portableDiff, err := process(md, orig, mc, *check)
		if err != nil {
			fatal(fmt.Errorf("%s: %w", md, err))
		}
		if !bytes.Equal(orig, out) {
			if *check {
				if portableDiff {
					changed++
					fmt.Fprintf(os.Stderr, "STALE: %s\n", rel(md))
				}
				continue
			}
			if err := os.WriteFile(md, out, 0o644); err != nil {
				fatal(err)
			}
			fmt.Fprintf(os.Stderr, "wrote: %s\n", rel(md))
		} else if *verbose {
			fmt.Fprintf(os.Stderr, "fresh: %s\n", rel(md))
		}
	}
	if *check && changed > 0 {
		fmt.Fprintf(os.Stderr, "\n%d file(s) out of date. Run `make docs` and commit.\n", changed)
		os.Exit(1)
	}
}

// process rewrites the generated blocks of one markdown file. It returns the new
// bytes and whether any *portable* block changed (the -check gate ignores
// non-portable blocks). Under -check the returned bytes still reflect ALL blocks
// so a pure non-portable change is visible to the caller but not gated.
func process(mdPath string, src []byte, mc *manifestCache, check bool) ([]byte, bool, error) {
	lines := splitLines(src)
	var out []string
	portableDiff := false
	i := 0
	for i < len(lines) {
		m := beginRe.FindStringSubmatch(lines[i])
		if m == nil {
			out = append(out, lines[i])
			i++
			continue
		}
		kind, name := m[1], m[2]
		// find matching END
		j := i + 1
		for j < len(lines) {
			e := endRe.FindStringSubmatch(lines[j])
			if e != nil {
				if e[1] != kind || e[2] != name {
					return nil, false, fmt.Errorf("marker mismatch: BEGIN:%s %s closed by END:%s %s at line %d", kind, name, e[1], e[2], j+1)
				}
				break
			}
			j++
		}
		if j >= len(lines) {
			return nil, false, fmt.Errorf("unterminated BEGIN:%s %s at line %d", kind, name, i+1)
		}
		mf, err := mc.get(filepath.Dir(mdPath))
		if err != nil {
			return nil, false, err
		}
		block, portable, err := generate(mf, kind, name)
		if err != nil {
			return nil, false, fmt.Errorf("%s:%s: %w", kind, name, err)
		}
		old := lines[i+1 : j]
		if portable && !equalStrings(old, block) {
			portableDiff = true
		}
		out = append(out, lines[i])
		out = append(out, block...)
		out = append(out, lines[j])
		i = j + 1
	}
	return joinLines(out, endsWithNewline(src)), portableDiff, nil
}

func generate(mf *manifest, kind, name string) (block []string, portable bool, err error) {
	switch kind {
	case "output":
		spec, ok := mf.Outputs[name]
		if !ok {
			return nil, false, fmt.Errorf("no output %q in %s", name, mf.path)
		}
		return renderOutput(mf, spec)
	case "snippet":
		spec, ok := mf.Snippets[name]
		if !ok {
			return nil, false, fmt.Errorf("no snippet %q in %s", name, mf.path)
		}
		return renderSnippet(mf, spec)
	case "file":
		spec, ok := mf.Files[name]
		if !ok {
			return nil, false, fmt.Errorf("no file %q in %s", name, mf.path)
		}
		return renderFile(mf, spec)
	}
	return nil, false, fmt.Errorf("unknown kind %q", kind)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "capture:", err)
	os.Exit(1)
}

func findMarkdown(root string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			switch d.Name() {
			case ".git", "target", "_build", "node_modules":
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(p, ".md") {
			out = append(out, p)
		}
		return nil
	})
	sort.Strings(out)
	return out, err
}

var repoRootCache string

func repoRoot() string {
	if repoRootCache != "" {
		return repoRootCache
	}
	d, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(d, ".git")); err == nil {
			repoRootCache = d
			return d
		}
		parent := filepath.Dir(d)
		if parent == d {
			break
		}
		d = parent
	}
	repoRootCache, _ = os.Getwd()
	return repoRootCache
}

func rel(p string) string {
	if r, err := filepath.Rel(repoRoot(), p); err == nil {
		return r
	}
	return p
}
