package canary

import (
	"errors"
	"path/filepath"
)

// Dir is the canary fixture directory under a repository root.
func Dir(root string) string {
	return filepath.Join(root, "tests", "canary")
}

// RootFixtures returns the fixture inventory under a repository root, keyed by base name.
// A root with no fixture gives an empty inventory, not an error.
func RootFixtures(root string) (map[string]Fixture, error) {
	return fixturesOf(rootRecords(root))
}

// rootRecords discovers the fixtures under a repository root. An absent or empty canary
// directory holds no fixture; that is an answer, not an error.
func rootRecords(root string) ([]fixtureRecord, error) {
	records, err := discoverFixtures(Dir(root))
	if errors.Is(err, ErrNoFixtures) {
		return nil, nil
	}
	return records, err
}

// fixturesOf keys each discovered fixture record by base name with its owning check.
func fixturesOf(records []fixtureRecord, err error) (map[string]Fixture, error) {
	if err != nil {
		return nil, err
	}
	result := make(map[string]Fixture, len(records))
	for _, record := range records {
		_, check, err := fixtureCheck(record.dir)
		if err != nil {
			return nil, err
		}
		result[filepath.Base(record.dir)] = Fixture{
			Dir: record.dir, Family: record.family, Check: fixtureScope(record.family, check),
		}
	}
	return result, nil
}
