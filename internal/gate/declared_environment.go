package gate

import "os"

// suppliedEnvironment splits the declared environment names into the name and value
// pairs this process supplies and the names it does not, both in the order given. The
// gate subject and the doctor row read one split, so they never disagree on a name.
func suppliedEnvironment(names []string) (supplied [][2]string, missing []string) {
	for _, name := range names {
		if value, ok := os.LookupEnv(name); ok {
			supplied = append(supplied, [2]string{name, value})
		} else {
			missing = append(missing, name)
		}
	}
	return supplied, missing
}

// DeclaredEnvironment returns the environment names that root's gate input manifest
// declares, with the manifest's validation. It reports false for an absent or invalid
// manifest, which declares nothing, so a caller can tell that state from an empty list.
func DeclaredEnvironment(root string) ([]string, bool) {
	m, _, _ := loadManifest(root)
	if m == nil {
		return nil, false
	}
	return m.Environment, true
}

// UnsuppliedEnvironment returns each environment name root's gate input manifest
// declares that this process does not supply. An absent or invalid manifest declares
// nothing, so it returns none; the gate reports that manifest state itself.
func UnsuppliedEnvironment(root string) []string {
	names, _ := DeclaredEnvironment(root)
	_, missing := suppliedEnvironment(names)
	return missing
}
