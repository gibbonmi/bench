package roadmap

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gibbonmi/bench/internal/capturetx"
	"github.com/gibbonmi/bench/internal/git"
	"github.com/gibbonmi/bench/internal/learnings"
	"github.com/gibbonmi/bench/internal/toon"
)

const captureUsage = "usage: bench capture drain [show|commit|abort] [<drain-id>]"

// CaptureCommand owns the lifecycle of one immutable capture-drain generation.
func CaptureCommand(args []string) (string, int) {
	if len(args) == 1 && (args[0] == "--help" || args[0] == "-h") {
		return captureUsage + "\n", 0
	}
	if len(args) == 0 || args[0] != "drain" {
		return captureUsage + "\n", 2
	}
	root, err := git.Root()
	if err != nil {
		return toon.NotInRepo() + "\n", 1
	}
	if len(args) == 1 {
		sources, err := captureSources()
		if err != nil {
			return toon.Errorf("capture drain failed", err.Error()) + "\n", 1
		}
		bundle, err := capturetx.Begin(root, sources)
		if err != nil {
			return toon.Errorf("capture drain failed", err.Error()) + "\n", 1
		}
		return renderCaptureBundle(bundle)
	}
	if len(args) != 3 {
		return captureUsage + "\n", 2
	}
	action, id := args[1], strings.TrimSpace(args[2])
	if id == "" {
		return captureUsage + "\n", 2
	}
	switch action {
	case "show":
		bundle, present, err := capturetx.Current(root)
		if err != nil {
			return toon.Errorf("capture drain failed", err.Error()) + "\n", 1
		}
		if !present || bundle.ID != id {
			return toon.Errorf("capture drain not found", id) + "\n", 1
		}
		return renderCaptureBundle(bundle)
	case "commit":
		if err := capturetx.Commit(root, id); err != nil {
			return toon.Errorf("capture drain commit failed", err.Error()) + "\n", 1
		}
		return "committed: " + id + "\n", 0
	case "abort":
		if err := capturetx.Abort(root, id); err != nil {
			return toon.Errorf("capture drain abort failed", err.Error()) + "\n", 1
		}
		return "aborted: " + id + "\n", 0
	default:
		return captureUsage + "\n", 2
	}
}

func captureSources() ([]capturetx.Source, error) {
	labels := []string{IdeasFile, learnings.JournalPath}
	sources := make([]capturetx.Source, 0, len(labels))
	for _, label := range labels {
		root, refusal, code := inboxRoot(label)
		if refusal != "" {
			return nil, fmt.Errorf("source %s: %s (exit %d)", label, strings.TrimSpace(refusal), code)
		}
		sources = append(sources, capturetx.Source{Name: label, Path: filepath.Join(root, filepath.FromSlash(label))})
	}
	return sources, nil
}

func renderCaptureBundle(bundle capturetx.Bundle) (string, int) {
	var out strings.Builder
	block, err := toon.TableTyped("drain", []string{"id", "state", "sources"}, [][]any{{bundle.ID, bundle.State, len(bundle.Sources)}})
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	out.WriteString(block)
	rows := make([][]any, len(bundle.Sources))
	for i, source := range bundle.Sources {
		rows[i] = []any{source.Name, source.Bytes, "sha256:" + source.Digest}
	}
	block, err = toon.TableTyped("sources", []string{"source", "bytes", "digest"}, rows)
	if err != nil {
		return toon.RenderError(err) + "\n", 1
	}
	out.WriteString(block)
	return out.String(), 0
}
