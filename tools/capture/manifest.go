package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// manifest is a module-local capture.json declaring the command / source /
// file behind each generated block name.
type manifest struct {
	path     string // absolute path to capture.json
	dir      string // module directory (dir of capture.json)
	Outputs  map[string]outputSpec  `json:"outputs"`
	Snippets map[string]snippetSpec `json:"snippets"`
	Files    map[string]fileSpec    `json:"files"`
}

type outputSpec struct {
	Dir        string   `json:"dir"`        // working dir relative to module, default "."
	Cmd        []string `json:"cmd"`        // argv, run without a shell
	Env        []string `json:"env"`        // extra KEY=VALUE appended to os.Environ
	Combined   bool     `json:"combined"`   // merge stderr into stdout
	AllowError bool     `json:"allowError"` // do not fail if the command exits non-zero
	ShowCmd    bool     `json:"showCmd"`    // prefix output with a "$ cmd" line inside the fence
	Lang       string   `json:"lang"`       // fence language, default "text"
	Normalize  []string `json:"normalize"`  // named scrubbers: path, addr, goid, gover, blank
	Portable   *bool    `json:"portable"`   // default true; false => skipped by -check
}

type snippetSpec struct {
	File   string `json:"file"`   // source file relative to module
	Region string `json:"region"` // region:<Region>:start .. region:<Region>:end
	Lang   string `json:"lang"`   // fence language; default inferred from extension
}

type fileSpec struct {
	File     string `json:"file"`
	Lang     string `json:"lang"`
	Portable *bool  `json:"portable"` // default true
}

func (s outputSpec) portable() bool  { return s.Portable == nil || *s.Portable }
func (s fileSpec) portable() bool    { return s.Portable == nil || *s.Portable }

type manifestCache struct {
	byDir map[string]*manifest
}

func newManifestCache() *manifestCache { return &manifestCache{byDir: map[string]*manifest{}} }

// get returns the nearest ancestor capture.json for a markdown file's directory.
func (c *manifestCache) get(dir string) (*manifest, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	d := abs
	for {
		if mf, ok := c.byDir[d]; ok {
			return mf, nil
		}
		mp := filepath.Join(d, "capture.json")
		if _, err := os.Stat(mp); err == nil {
			mf, err := loadManifest(mp)
			if err != nil {
				return nil, err
			}
			c.byDir[d] = mf
			return mf, nil
		}
		parent := filepath.Dir(d)
		if parent == d {
			return nil, fmt.Errorf("no capture.json found for %s", dir)
		}
		d = parent
	}
}

func loadManifest(path string) (*manifest, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var mf manifest
	dec := json.NewDecoder(newJSONCReader(b))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&mf); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	mf.path = path
	mf.dir = filepath.Dir(path)
	return &mf, nil
}
