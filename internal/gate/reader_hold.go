package gate

// The reader hold: a red whole landing run writes the hold of its green marker, and the
// landing plan runs the whole gate while the marker is that commit. A whole green landing
// moves the marker, so the hold stops applying with no write. The same red run names its
// suspects: the reader landings since the marker.

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/gibbonmi/bench/internal/durablefile"
	benchgit "github.com/gibbonmi/bench/internal/git"
)

// holdWholeGate writes the reader hold of a red whole run whose landing value names a green
// marker, and answers the suspects of the red. A run without such a landing value writes
// nothing and has no suspects. A background run is speculative, so its red writes no hold either.
func holdWholeGate(ctx context.Context, replacer durablefile.Replacer, storageRoot string) ([]string, error) {
	landing, landed := landingOf(ctx)
	if !landed || landing.marker == "" || backgroundGradeOf(ctx) != nil {
		return nil, nil
	}
	gitdir, err := benchgit.CommonDir(storageRoot)
	if err != nil {
		return nil, err
	}
	if err := writeReaderHold(replacer, gitdir, landing.marker); err != nil {
		return nil, err
	}
	// The suspects only point the operator at commits to inspect, so a chain that cannot be
	// read names none and leaves the red as it is.
	_, suspects, err := readerChain(storageRoot, landing)
	if err != nil {
		return nil, nil
	}
	return suspects, nil
}

// writeReaderHold installs the hold of marker in the common Git directory gitdir through
// replacer.
func writeReaderHold(replacer durablefile.Replacer, gitdir, marker string) error {
	return durableReplaceRecordAt(replacer, gitdir, benchgit.ReaderHoldFile, []byte(marker))
}

// readerHeld answers whether the reader hold of root names marker. An absent hold names no
// marker.
func readerHeld(root, marker string) (bool, error) {
	gitdir, err := benchgit.CommonDir(root)
	if err != nil {
		return false, err
	}
	data, err := os.ReadFile(filepath.Join(gitdir, benchgit.ReaderHoldFile))
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return string(bytes.TrimSuffix(data, []byte("\n"))) == marker, nil
}
