package adopt

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Manifest struct {
	KitVersion string
	hashes     map[string]string
}

func (m Manifest) Hash(rel string) string {
	if m.hashes == nil {
		return ""
	}
	return m.hashes[rel]
}

func (m Manifest) Rows() []manifestRow {
	rows := make([]manifestRow, 0, len(m.hashes))
	for rel, hash := range m.hashes {
		rows = append(rows, manifestRow{rel, hash})
	}
	return rows
}

func ReadManifest(path string) (Manifest, error) {
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return Manifest{hashes: map[string]string{}}, nil
	}
	if err != nil {
		return Manifest{}, err
	}
	defer f.Close()
	m := Manifest{hashes: map[string]string{}}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue
		}
		if parts[0] == "#kit" {
			m.KitVersion = parts[1]
			continue
		}
		if strings.HasPrefix(parts[0], "#") {
			continue
		}
		m.hashes[parts[0]] = parts[1]
	}
	return m, sc.Err()
}

type manifestRow struct {
	rel, hash string
}

func writeManifest(path, version string, rows []manifestRow) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	_, err = f.Write(manifestBytes(version, rows))
	cerr := f.Close()
	if err != nil {
		_ = os.Remove(tmp)
		return err
	}
	if cerr != nil {
		_ = os.Remove(tmp)
		return cerr
	}
	return os.Rename(tmp, path)
}

func manifestBytes(version string, rows []manifestRow) []byte {
	var b strings.Builder
	b.WriteString("#kit\t" + version + "\n")
	for _, r := range rows {
		b.WriteString(r.rel + "\t" + r.hash + "\n")
	}
	return []byte(b.String())
}

func fingerprintPath(path string) (string, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(path)
		if err != nil {
			return "", err
		}
		return hashBytes([]byte("symlink:" + target + "\n")), nil
	}
	if !info.Mode().IsRegular() {
		return "", errors.New("not a regular file")
	}
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func hashBytes(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
