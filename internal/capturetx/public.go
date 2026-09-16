// Package capturetx owns the short, repository-wide transaction used by capture
// writers and capture drains.
package capturetx

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Source identifies one live capture document and its stable evidence name.
type Source struct {
	Name string
	Path string
}

// Bundle identifies the immutable generation currently under review.
type Bundle struct {
	ID      string
	State   string
	Sources []SourceFact
}

// SourceFact describes one source in a bundle.
type SourceFact struct {
	Name   string
	Bytes  int
	Digest string
}

// Append serializes one read-modify-write with drain cutover and recovery.
func Append(root string, source Source, compose func([]byte) ([]byte, error)) error {
	return withLock(root, func(common string) error {
		if err := recoverTransaction(common); err != nil {
			return err
		}
		current, err := readOptional(source.Path)
		if err != nil {
			return err
		}
		next, err := compose(current)
		if err != nil {
			return err
		}
		return replace(source.Path, next)
	})
}

// Begin seals sources as one immutable generation and immediately clears their
// live documents for later appends. A repeated call returns the open generation.
func Begin(root string, sources []Source) (Bundle, error) {
	var bundle Bundle
	err := withLock(root, func(common string) error {
		if err := recoverTransaction(common); err != nil {
			return err
		}
		if current, err := readManifest(common); err == nil {
			bundle = project(current)
			return nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		sorted := append([]Source(nil), sources...)
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].Name < sorted[j].Name })
		m := manifest{Schema: 1, State: "preparing"}
		var identity []byte
		for i, source := range sorted {
			data, err := readOptional(source.Path)
			if err != nil {
				return fmt.Errorf("read capture source %s: %w", source.Name, err)
			}
			m.Sources = append(m.Sources, manifestSource{
				Name: source.Name, Path: source.Path, Blob: fmt.Sprintf("source-%d.bin", i),
				Bytes: len(data), Digest: digest(data),
			})
			identity = append(identity, source.Name...)
			identity = append(identity, 0)
			identity = append(identity, data...)
			identity = append(identity, 0)
		}
		sum := sha256.Sum256(identity)
		m.ID = "d-" + hex.EncodeToString(sum[:6])
		dir := transactionDir(common)
		staged, err := os.MkdirTemp(common, directoryName+".tmp-")
		if err != nil {
			return fmt.Errorf("create capture transaction: %w", err)
		}
		defer os.RemoveAll(staged)
		for _, source := range m.Sources {
			data, err := readOptional(source.Path)
			if err != nil {
				return err
			}
			if digest(data) != source.Digest {
				return fmt.Errorf("capture source %s changed before cutover", source.Name)
			}
			if err := writeNew(filepath.Join(staged, source.Blob), data); err != nil {
				return err
			}
		}
		if err := writeManifestDir(staged, m); err != nil {
			return err
		}
		if err := os.Rename(staged, dir); err != nil {
			return fmt.Errorf("publish capture transaction: %w", err)
		}
		if err := recoverTransaction(common); err != nil {
			return err
		}
		sealed, err := readManifest(common)
		if err != nil {
			return err
		}
		bundle = project(sealed)
		return nil
	})
	return bundle, err
}

// Current returns the open generation, if one exists.
func Current(root string) (Bundle, bool, error) {
	var bundle Bundle
	var present bool
	err := withLock(root, func(common string) error {
		if err := recoverTransaction(common); err != nil {
			return err
		}
		m, err := readManifest(common)
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		if err != nil {
			return err
		}
		bundle, present = project(m), true
		return nil
	})
	return bundle, present, err
}

// Sources returns the requested documents from one immutable open generation.
func Sources(root string, names []string) (map[string][]byte, bool, error) {
	data := map[string][]byte{}
	var present bool
	err := withLock(root, func(common string) error {
		if err := recoverTransaction(common); err != nil {
			return err
		}
		m, err := readManifest(common)
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		if err != nil {
			return err
		}
		wanted := make(map[string]bool, len(names))
		for _, name := range names {
			wanted[name] = true
		}
		for _, source := range m.Sources {
			if !wanted[source.Name] {
				continue
			}
			body, err := os.ReadFile(filepath.Join(transactionDir(common), source.Blob))
			if err != nil {
				return err
			}
			if digest(body) != source.Digest {
				return fmt.Errorf("sealed capture source %s failed its digest", source.Name)
			}
			data[source.Name] = body
		}
		if len(data) != len(wanted) {
			return fmt.Errorf("sealed capture generation does not contain every requested source")
		}
		present = true
		return nil
	})
	return data, present, err
}

// Commit retires only the sealed generation. Live post-cut entries are unchanged.
func Commit(root, id string) error {
	return withLock(root, func(common string) error {
		if err := recoverTransaction(common); err != nil {
			return err
		}
		m, err := requireManifest(common, id)
		if err != nil {
			return err
		}
		m.State = "committing"
		if err := writeManifest(common, m); err != nil {
			return err
		}
		return recoverTransaction(common)
	})
}

// Abort restores the sealed generation before every post-cut live suffix.
func Abort(root, id string) error {
	return withLock(root, func(common string) error {
		if err := recoverTransaction(common); err != nil {
			return err
		}
		m, err := requireManifest(common, id)
		if err != nil {
			return err
		}
		dir := transactionDir(common)
		for i := range m.Sources {
			sealed, err := os.ReadFile(filepath.Join(dir, m.Sources[i].Blob))
			if err != nil {
				return err
			}
			live, err := readOptional(m.Sources[i].Path)
			if err != nil {
				return err
			}
			after := join(sealed, live)
			m.Sources[i].AfterBlob = fmt.Sprintf("after-%d.bin", i)
			m.Sources[i].AfterDigest = digest(after)
			if err := replace(filepath.Join(dir, m.Sources[i].AfterBlob), after); err != nil {
				return err
			}
		}
		m.State = "aborting"
		if err := writeManifest(common, m); err != nil {
			return err
		}
		return recoverTransaction(common)
	})
}
