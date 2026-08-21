package worktree

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gibbonmi/bench/internal/jsonfile"
	"github.com/gibbonmi/bench/internal/landing"
)

const buildOutputDeclarationLimit = 16 * 1024

type buildOutputDeclaration struct {
	Schema int      `json:"schema"`
	Paths  []string `json:"paths"`
}

func loadBuildOutputs(root string) ([]string, []byte, error) {
	path := filepath.Join(root, ".bench", "build-outputs.json")
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, []byte("absent"), nil
	}
	if err != nil {
		return nil, nil, err
	}
	evidence := canonicalParts([]byte(info.Mode().String()), []byte(strconv.FormatInt(info.Size(), 10)))
	if !info.Mode().IsRegular() {
		return nil, evidence, errors.New("build-output declaration is not a regular file")
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, evidence, err
	}
	data, readErr := io.ReadAll(io.LimitReader(file, buildOutputDeclarationLimit+1))
	closeErr := file.Close()
	evidence = canonicalParts(evidence, data)
	if readErr != nil || closeErr != nil {
		return nil, evidence, errors.Join(readErr, closeErr)
	}
	if len(data) > buildOutputDeclarationLimit || !utf8.Valid(data) {
		return nil, evidence, errors.New("build-output declaration is malformed")
	}
	var declaration buildOutputDeclaration
	if err := jsonfile.DecodeDocument(data, &declaration); err != nil || declaration.Schema != 1 || declaration.Paths == nil {
		return nil, evidence, errors.New("build-output declaration is malformed")
	}
	for _, entry := range declaration.Paths {
		if !validBuildOutputPath(entry) {
			return nil, evidence, errors.New("build-output declaration is malformed")
		}
	}
	return declaration.Paths, evidence, nil
}

func validBuildOutputPath(entry string) bool {
	for _, r := range entry {
		if r < 0x20 || r == 0x7f {
			return false
		}
	}
	base := strings.TrimSuffix(entry, "/")
	native := filepath.FromSlash(base)
	return base != "" && base != "." && !filepath.IsAbs(native) && filepath.Clean(native) == native &&
		filepath.ToSlash(native) == base && native != ".." && !strings.HasPrefix(native, ".."+string(filepath.Separator))
}

func ignoredWithinLandingAllowance(inventory IgnoredInventory, declared []string) bool {
	return ignoredWithinDeclaredOutputs(inventory, declared, landing.RuntimeIgnoredPath)
}

func ignoredWithinDeclaredOutputs(inventory IgnoredInventory, declared []string, additional func(string) bool) bool {
	if inventory.Count == 0 || inventory.Uncertain || inventory.OverLimit || inventory.AtLeast {
		return false
	}
	for _, ignored := range inventory.Paths {
		if additional(ignored) {
			continue
		}
		contained := false
		for _, entry := range declared {
			if strings.HasSuffix(entry, "/") {
				contained = strings.HasPrefix(ignored, entry)
			} else {
				contained = ignored == entry
			}
			if contained {
				break
			}
		}
		if !contained {
			return false
		}
	}
	return true
}
