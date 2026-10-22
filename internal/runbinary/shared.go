package runbinary

import (
	"context"
	"sync"
)

// Shared holds one selection for every caller of a test package run. Raw and Present carry
// the Env value that the creator read, so a later change to the process environment does
// not move the selection. The first Path call selects through Factory.Select, and every
// later call answers that selection or its refusal.
type Shared struct {
	Factory Factory
	Raw     string
	Present bool

	once      sync.Once
	selection *Selection
	err       error
}

// Path answers the shared executable. Only the first call grades sourceRoot.
func (s *Shared) Path(ctx context.Context, sourceRoot string) (string, error) {
	s.once.Do(func() { s.selection, s.err = s.Factory.Select(ctx, sourceRoot, s.Raw, s.Present) })
	if s.err != nil {
		return "", s.err
	}
	return s.selection.Path, nil
}

// Close releases an owned selection after every caller has returned.
func (s *Shared) Close() error {
	return s.selection.Close()
}
