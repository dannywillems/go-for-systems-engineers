package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func renderOutput(mf *manifest, s outputSpec) ([]string, bool, error) {
	if len(s.Cmd) == 0 {
		return nil, false, fmt.Errorf("empty cmd")
	}
	wd := mf.dir
	if s.Dir != "" {
		wd = filepath.Join(mf.dir, s.Dir)
	}
	cmd := exec.Command(s.Cmd[0], s.Cmd[1:]...)
	cmd.Dir = wd
	cmd.Env = append(os.Environ(), s.Env...)
	var stdout, stderr bytes.Buffer
	if s.Combined {
		cmd.Stdout = &stdout
		cmd.Stderr = &stdout
	} else {
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
	}
	runErr := cmd.Run()
	if runErr != nil && !s.AllowError {
		return nil, false, fmt.Errorf("running %v in %s: %w\nstderr:\n%s", s.Cmd, wd, runErr, stderr.String())
	}
	text := stdout.String()
	text = applyNormalizers(text, s.Normalize, mf.dir)

	lang := s.Lang
	if lang == "" {
		lang = "text"
	}
	var body []string
	if s.ShowCmd {
		body = append(body, "$ "+shellJoin(s.Cmd))
	}
	body = append(body, splitLines([]byte(text))...)
	body = trimTrailingBlank(body)
	return fence(lang, body), s.portable(), nil
}

func renderSnippet(mf *manifest, s snippetSpec) ([]string, bool, error) {
	fp := filepath.Join(mf.dir, s.File)
	b, err := os.ReadFile(fp)
	if err != nil {
		return nil, false, err
	}
	region, err := extractRegion(b, s.Region)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", s.File, err)
	}
	lang := s.Lang
	if lang == "" {
		lang = langFromExt(s.File)
	}
	return fence(lang, region), true, nil
}

func renderFile(mf *manifest, s fileSpec) ([]string, bool, error) {
	fp := filepath.Join(mf.dir, s.File)
	b, err := os.ReadFile(fp)
	if err != nil {
		return nil, false, err
	}
	lang := s.Lang
	if lang == "" {
		lang = langFromExt(s.File)
	}
	body := trimTrailingBlank(splitLines(b))
	return fence(lang, body), s.portable(), nil
}

// extractRegion returns the lines strictly between `region:<name>:start` and
// `region:<name>:end` (matched anywhere on a line, so any comment syntax works),
// with common leading indentation removed.
func extractRegion(src []byte, name string) ([]string, error) {
	start := "region:" + name + ":start"
	end := "region:" + name + ":end"
	lines := splitLines(src)
	si, ei := -1, -1
	for i, l := range lines {
		if si < 0 && strings.Contains(l, start) {
			si = i
			continue
		}
		if si >= 0 && strings.Contains(l, end) {
			ei = i
			break
		}
	}
	if si < 0 {
		return nil, fmt.Errorf("region %q start not found", name)
	}
	if ei < 0 {
		return nil, fmt.Errorf("region %q end not found", name)
	}
	return dedent(trimTrailingBlank(trimLeadingBlank(lines[si+1 : ei]))), nil
}

func shellJoin(argv []string) string {
	parts := make([]string, len(argv))
	for i, a := range argv {
		if strings.ContainsAny(a, " \t'\"") {
			parts[i] = "'" + strings.ReplaceAll(a, "'", `'\''`) + "'"
		} else {
			parts[i] = a
		}
	}
	return strings.Join(parts, " ")
}
