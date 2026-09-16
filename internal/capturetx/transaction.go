package capturetx

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"

	benchgit "github.com/gibbonmi/bench/internal/git"
)

const (
	directoryName = "bench-capture-drain"
	manifestName  = "manifest.json"
	lockName      = "bench-capture-drain.lock"
	lockWait      = 2 * time.Second
)

type manifest struct {
	Schema  int              `json:"schema"`
	ID      string           `json:"id"`
	State   string           `json:"state"`
	Sources []manifestSource `json:"sources"`
}

type manifestSource struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Blob        string `json:"blob"`
	Bytes       int    `json:"bytes"`
	Digest      string `json:"digest"`
	AfterBlob   string `json:"after_blob,omitempty"`
	AfterDigest string `json:"after_digest,omitempty"`
}

func recoverTransaction(common string) error {
	m, err := readManifest(common)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	switch m.State {
	case "preparing":
		for _, source := range m.Sources {
			current, err := readOptional(source.Path)
			if err != nil {
				return err
			}
			if len(current) == 0 {
				continue
			}
			if digest(current) != source.Digest {
				return fmt.Errorf("capture source %s changed during cutover", source.Name)
			}
			if err := replace(source.Path, nil); err != nil {
				return err
			}
		}
		m.State = "sealed"
		return writeManifest(common, m)
	case "sealed":
		return nil
	case "committing":
		return os.RemoveAll(transactionDir(common))
	case "aborting":
		for _, source := range m.Sources {
			after, err := os.ReadFile(filepath.Join(transactionDir(common), source.AfterBlob))
			if err != nil {
				return err
			}
			current, err := readOptional(source.Path)
			if err != nil {
				return err
			}
			if digest(current) == source.AfterDigest {
				continue
			}
			sealed, err := os.ReadFile(filepath.Join(transactionDir(common), source.Blob))
			if err != nil {
				return err
			}
			if !bytes.Equal(after, join(sealed, current)) {
				return fmt.Errorf("capture source %s changed during abort", source.Name)
			}
			if err := replace(source.Path, after); err != nil {
				return err
			}
		}
		return os.RemoveAll(transactionDir(common))
	default:
		return fmt.Errorf("capture transaction %s has unknown state %q", m.ID, m.State)
	}
}

func requireManifest(common, id string) (manifest, error) {
	m, err := readManifest(common)
	if err != nil {
		return manifest{}, err
	}
	if m.ID != id {
		return manifest{}, fmt.Errorf("capture transaction %s is not open", id)
	}
	if m.State != "sealed" {
		return manifest{}, fmt.Errorf("capture transaction %s is %s", id, m.State)
	}
	return m, nil
}

func project(m manifest) Bundle {
	b := Bundle{ID: m.ID, State: m.State}
	for _, source := range m.Sources {
		b.Sources = append(b.Sources, SourceFact{Name: source.Name, Bytes: source.Bytes, Digest: source.Digest})
	}
	return b
}

func withLock(root string, run func(common string) error) error {
	common, err := benchgit.CommonDir(root)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filepath.Join(common, lockName), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()
	deadline := time.Now().Add(lockWait)
	for {
		err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			break
		}
		if !errors.Is(err, syscall.EWOULDBLOCK) || time.Now().After(deadline) {
			return fmt.Errorf("capture transaction lock is busy: %w", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
	defer syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
	return run(common)
}

func transactionDir(common string) string { return filepath.Join(common, directoryName) }

func readManifest(common string) (manifest, error) {
	data, err := os.ReadFile(filepath.Join(transactionDir(common), manifestName))
	if err != nil {
		return manifest{}, err
	}
	var m manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return manifest{}, fmt.Errorf("parse capture transaction manifest: %w", err)
	}
	if m.Schema != 1 || m.ID == "" {
		return manifest{}, fmt.Errorf("unsupported capture transaction manifest")
	}
	return m, nil
}

func writeManifest(common string, m manifest) error {
	return writeManifestDir(transactionDir(common), m)
}

func writeManifestDir(dir string, m manifest) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return replace(filepath.Join(dir, manifestName), append(data, '\n'))
}

func readOptional(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	return data, err
}

func writeNew(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return err
	}
	return file.Close()
}

func replace(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(name, 0o644); err != nil {
		return err
	}
	return os.Rename(name, path)
}

func digest(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func join(first, second []byte) []byte {
	if len(first) == 0 {
		return append([]byte(nil), second...)
	}
	if len(second) == 0 {
		return append([]byte(nil), first...)
	}
	out := append([]byte(nil), first...)
	if out[len(out)-1] != '\n' {
		out = append(out, '\n')
	}
	return append(out, second...)
}
