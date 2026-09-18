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

// UnsuppliedEnvironment returns each environment name root's gate input manifest
// declares that this process does not supply. An absent or invalid manifest declares
// nothing, so it returns none; the gate reports that manifest state itself.
func UnsuppliedEnvironment(root string) []string {
	m, _, _ := loadManifest(root)
	if m == nil {
		return nil
	}
	_, missing := suppliedEnvironment(m.Environment)
	return missing
}
