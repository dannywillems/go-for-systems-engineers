package main

import (
	"path/filepath"
	"regexp"
	"strings"
)

func splitLines(b []byte) []string {
	s := strings.ReplaceAll(string(b), "\r\n", "\n")
	s = strings.TrimSuffix(s, "\n")
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

func joinLines(lines []string, trailingNewline bool) []byte {
	s := strings.Join(lines, "\n")
	if trailingNewline {
		s += "\n"
	}
	return []byte(s)
}

func endsWithNewline(b []byte) bool {
	return len(b) > 0 && b[len(b)-1] == '\n'
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func fence(lang string, body []string) []string {
	out := make([]string, 0, len(body)+2)
	out = append(out, "```"+lang)
	out = append(out, body...)
	out = append(out, "```")
	return out
}

func trimTrailingBlank(lines []string) []string {
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func trimLeadingBlank(lines []string) []string {
	for len(lines) > 0 && strings.TrimSpace(lines[0]) == "" {
		lines = lines[1:]
	}
	return lines
}

// dedent removes the longest common leading run of spaces/tabs.
func dedent(lines []string) []string {
	prefix := ""
	first := true
	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			continue
		}
		lead := l[:len(l)-len(strings.TrimLeft(l, " \t"))]
		if first {
			prefix = lead
			first = false
			continue
		}
		prefix = commonPrefix(prefix, lead)
	}
	if prefix == "" {
		return lines
	}
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = strings.TrimPrefix(l, prefix)
	}
	return out
}

func commonPrefix(a, b string) string {
	n := min(len(a), len(b))
	i := 0
	for i < n && a[i] == b[i] {
		i++
	}
	return a[:i]
}

func langFromExt(path string) string {
	switch filepath.Ext(path) {
	case ".go":
		return "go"
	case ".rs":
		return "rust"
	case ".ml", ".mli":
		return "ocaml"
	case ".swift":
		return "swift"
	case ".kt", ".kts":
		return "kotlin"
	case ".sh":
		return "bash"
	case ".toml":
		return "toml"
	case ".json":
		return "json"
	case ".md":
		return "markdown"
	default:
		return "text"
	}
}

var (
	reAddr  = regexp.MustCompile(`0x[0-9a-fA-F]+`)
	reGoid  = regexp.MustCompile(`[Gg]oroutine \d+`)
	reGover = regexp.MustCompile(`go1\.\d+(\.\d+)?`)
	reMulti = regexp.MustCompile(`\n{3,}`)
)

// applyNormalizers scrubs machine- or run-specific noise so deterministic
// captures diff cleanly. Each normalizer is opt-in per output in capture.json.
func applyNormalizers(s string, names []string, moduleDir string) string {
	for _, n := range names {
		switch n {
		case "path":
			abs, _ := filepath.Abs(moduleDir)
			s = strings.ReplaceAll(s, abs, ".")
			s = strings.ReplaceAll(s, repoRoot(), ".")
		case "addr":
			s = reAddr.ReplaceAllString(s, "0xADDR")
		case "goid":
			s = reGoid.ReplaceAllString(s, "goroutine N")
		case "gover":
			s = reGover.ReplaceAllString(s, "goX")
		case "blank":
			s = reMulti.ReplaceAllString(s, "\n\n")
		}
	}
	return s
}

// newJSONCReader strips // line comments (outside strings) so capture.json can
// be annotated. Block comments are not supported.
func newJSONCReader(b []byte) *strings.Reader {
	var out strings.Builder
	inStr, esc := false, false
	for i := 0; i < len(b); i++ {
		c := b[i]
		if inStr {
			out.WriteByte(c)
			if esc {
				esc = false
			} else if c == '\\' {
				esc = true
			} else if c == '"' {
				inStr = false
			}
			continue
		}
		if c == '"' {
			inStr = true
			out.WriteByte(c)
			continue
		}
		if c == '/' && i+1 < len(b) && b[i+1] == '/' {
			for i < len(b) && b[i] != '\n' {
				i++
			}
			if i < len(b) {
				out.WriteByte('\n')
			}
			continue
		}
		out.WriteByte(c)
	}
	return strings.NewReader(out.String())
}
