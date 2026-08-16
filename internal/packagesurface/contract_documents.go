package packagesurface

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	kitpayload "github.com/gibbonmi/bench"
)

const consumerPayloadPath = ".bench/consumer-payload.json"

// ContractDocumentInputs returns the consumer-managed documents lifecycle contracts read,
// plus the inventory that selects them. The payload remains the authoritative registry: a
// new consumer document joins this result by changing that file, not a second path list.
func ContractDocumentInputs(root string) ([]string, error) {
	rows, err := consumerPayloadRows(root)
	if err != nil {
		return nil, err
	}

	kitOnly := kitpayload.PayloadKitOnlyPrefixes(rows)
	paths := []string{consumerPayloadPath}
	for _, row := range kitpayload.PayloadConsumerRows(rows) {
		if row.Tree {
			documents, err := consumerTreeDocuments(root, row.Source, kitOnly)
			if err != nil {
				return nil, err
			}
			paths = append(paths, documents...)
			continue
		}
		if err := requireConsumerAsset(root, row.Source); err != nil {
			return nil, err
		}
		if strings.HasSuffix(row.Source, ".md") {
			paths = append(paths, row.Source)
		}
	}
	slices.Sort(paths)
	return slices.Compact(paths), nil
}

// consumerPayloadRows reads the graded root's allowlist through the canonical payload
// reader: classification before the open, and one parser for syntax and row validity,
// so this inventory cannot resolve a document set from rows the allowlist forbids. The
// allowlist is this check's subject, so absence is a defect here rather than optional.
func consumerPayloadRows(root string) ([]kitpayload.PayloadRow, error) {
	rows, absent, err := kitpayload.PayloadRowsAt(filepath.Join(root, filepath.FromSlash(consumerPayloadPath)))
	if absent {
		return nil, fmt.Errorf("%s is absent: the consumer payload inventory has no source", consumerPayloadPath)
	}
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func consumerTreeDocuments(root, source string, kitOnly []string) ([]string, error) {
	path := filepath.Join(root, filepath.FromSlash(source))
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat consumer inventory tree %q: %w", source, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("consumer inventory tree %q is not a directory", source)
	}
	var documents []string
	err = filepath.WalkDir(path, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if kitpayload.PayloadExcluded(rel, kitOnly) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.Type().IsRegular() && strings.HasSuffix(rel, ".md") {
			documents = append(documents, rel)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk consumer inventory tree %q: %w", source, err)
	}
	return documents, nil
}

func requireConsumerAsset(root, source string) error {
	info, err := os.Stat(filepath.Join(root, filepath.FromSlash(source)))
	if err != nil {
		return fmt.Errorf("stat consumer inventory asset %q: %w", source, err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("consumer inventory asset %q is not a regular file", source)
	}
	return nil
}
