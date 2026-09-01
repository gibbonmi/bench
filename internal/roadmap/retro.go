package roadmap

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/retros"
	"github.com/gibbonmi/bench/internal/toon"
	"github.com/gibbonmi/bench/internal/usage"
)

var retroGrammar = usage.Grammar{
	Cmd:     "bench retro",
	Help:    "usage: bench retro <slug> --body <markdown>",
	Flags:   []usage.Flag{{Name: "--body", HasValue: true, NoEmptyValue: true}},
	MaxArgs: 1,
}

// RetroCommand validates and writes one primary-local implementation retrospective.
func RetroCommand(args []string) (string, int) {
	parsed, line, code := usage.Parse(retroGrammar, args)
	if line != "" {
		return line + "\n", code
	}
	if len(parsed.Positionals) != 1 || !retros.ValidSlug(parsed.Positionals[0]) {
		return retroGrammar.Help + "\n", 2
	}
	body, present := parsed.Flags["--body"]
	if !present {
		return toon.MissingArg(retroGrammar.Cmd, "--body") + "\n", 2
	}
	if err := retros.Parse([]byte(body)); err != nil {
		return toon.Errorf("invalid retrospective", err.Error()) + "\n", 1
	}
	path := retros.Path(parsed.Positionals[0])
	root, refusal, code := inboxRoot(path)
	if refusal != "" {
		return refusal, code
	}
	file := filepath.Join(root, filepath.FromSlash(path))
	if err := refuseSymlinkComponents(root, path); err != nil {
		return cannotWrite(path, err), 1
	}
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		return cannotWrite(path, err), 1
	}
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return cannotWrite(path, err), 1
	}
	if _, err := f.Write([]byte(body)); err != nil {
		_ = f.Close()
		return cannotWrite(path, err), 1
	}
	if err := f.Close(); err != nil {
		return cannotWrite(path, err), 1
	}
	return "captured: " + parsed.Positionals[0] + "\n", 0
}

func refuseSymlinkComponents(root, relPath string) error {
	current := root
	for _, component := range strings.Split(filepath.Clean(filepath.FromSlash(relPath)), string(filepath.Separator)) {
		current = filepath.Join(current, component)
		info, err := os.Lstat(current)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("destination component is a symbolic link: %s", component)
		}
	}
	return nil
}
