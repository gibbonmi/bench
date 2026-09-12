package reviewrecord

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/gibbonmi/bench/internal/bounds"
	benchgit "github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/toon"
)

func RecordPath(spec string) (string, error) {
	slug, err := Slug(spec)
	if err != nil {
		return "", err
	}
	return "reviews/" + slug + ".md", nil
}

// Slug is the <slug> of a checkpoint spec path. The record path and the
// checkpoint's refusal route both read this one grammar, so the file the
// checkpoint opens and the spec its route names can never disagree.
func Slug(spec string) (string, error) {
	parts := strings.Split(spec, "/")
	if !safeRelative(spec) || len(parts) != 3 || parts[0] != "specs" || parts[2] != "spec.md" {
		return "", errors.New("invalid checkpoint spec path; use specs/<slug>/spec.md")
	}
	return parts[1], nil
}

func safeRelative(path string) bool {
	return fs.ValidPath(path) && path != "." && !strings.Contains(path, "\\") && strings.IndexFunc(path, unicode.IsControl) < 0 && toon.Representable(path)
}

func readFile(root, relative string) ([]byte, error) {
	if !safeRelative(relative) {
		return nil, errors.New("invalid escaping or control-byte path")
	}
	path := root
	parts := strings.Split(relative, "/")
	for i, component := range parts {
		path = filepath.Join(path, component)
		info, err := os.Lstat(path)
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrMissing
		}
		if err != nil {
			return nil, err
		}
		if info.Mode()&os.ModeSymlink != 0 || (i < len(parts)-1 && !info.IsDir()) {
			return nil, fmt.Errorf("invalid nonregular or symlink path %q", relative)
		}
	}
	c := bounds.ClassifyNoFollow(path)
	if c.State != bounds.StateParsed {
		return nil, fmt.Errorf("invalid record %q: %s %s", relative, c.State, c.Reason)
	}
	return c.Data, nil
}

func Read(root, spec string) (Record, error) {
	path, err := RecordPath(spec)
	if err != nil {
		return Record{}, err
	}
	data, err := readFile(root, path)
	if err != nil {
		return Record{}, err
	}
	return parseRecord(data, spec)
}

func ReadTree(root, tree, spec string) (Record, error) {
	path, err := RecordPath(spec)
	if err != nil {
		return Record{}, err
	}
	data, err := benchgit.ReadTreeFile(root, tree, path)
	if err != nil {
		return Record{}, err
	}
	return parseRecord(data, spec)
}

func parseRecord(data []byte, spec string) (Record, error) {
	payload, err := fenced(data, "bench-review-record")
	if err != nil {
		return Record{}, err
	}
	record, err := Parse(payload)
	if err != nil {
		return record, err
	}
	if record.Spec != spec {
		return record, errors.New("invalid record spec identity")
	}
	return record, nil
}
