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
	return len(undeclaredLandingIgnoredPaths(inventory, declared)) == 0 &&
		inventory.Count > 0 && !inventory.Uncertain && !inventory.OverLimit && !inventory.AtLeast
}

func undeclaredLandingIgnoredPaths(inventory IgnoredInventory, declared []string) []string {
	var foreign []string
	for _, ignored := range inventory.Paths {
		if landing.RuntimeIgnoredPath(ignored) || landing.LocalCapturePath(ignored) || ignoredWithinBuildOutput(ignored, declared) {
			continue
		}
		foreign = append(foreign, ignored)
	}
	return foreign
}

func ignoredWithinBuildOutput(ignored string, declared []string) bool {
	for _, entry := range declared {
		if (strings.HasSuffix(entry, "/") && strings.HasPrefix(ignored, entry)) || ignored == entry {
			return true
		}
	}
	return false
}

func ignoredWithinDeclaredOutputs(inventory IgnoredInventory, declared []string, additional func(string) bool) bool {
	if inventory.Count == 0 || inventory.Uncertain || inventory.OverLimit || inventory.AtLeast {
		return false
	}
	for _, ignored := range inventory.Paths {
		if additional(ignored) {
			continue
		}
		if !ignoredWithinBuildOutput(ignored, declared) {
			return false
		}
	}
	return true
}
