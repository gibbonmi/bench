package intent

import (
	"errors"
	"fmt"
)

func validEntry(entry Entry) error {
	if entry.Key == "" || entry.CreatedAt.IsZero() {
		return errors.New("entry requires key and creation time")
	}
	switch entry.Kind {
	case KindShift, KindWorktree, KindClaudeAgent:
		return nil
	default:
		return fmt.Errorf("entry %q has unknown writer kind %q", entry.Key, entry.Kind)
	}
}
