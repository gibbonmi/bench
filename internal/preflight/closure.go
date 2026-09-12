package preflight

import "strings"

type closureKind string

const (
	fixtureClosure  closureKind = "fixture"
	registryClosure closureKind = "registry"
)

type closureRequirement struct {
	ticket string
	entry  string
	path   string
	kind   closureKind
}

// missingClosures derives omitted closure facts for legacy verdicts and proposals.
func missingClosures(f Facts, kind closureKind) []closureRequirement {
	var missing []closureRequirement
	seen := map[string]bool{}
	for _, ticket := range f.Tickets {
		owned := ownedPaths(ticket)
		for _, entry := range ticket.Writes {
			required := f.WritesBoundFiles[entry]
			if kind == fixtureClosure {
				required = f.WritesFixturePins[entry]
			}
			for _, path := range required {
				if pathCovered(path, owned) {
					continue
				}
				key := ticket.Name + "\x00" + entry + "\x00" + path + "\x00" + string(kind)
				if !seen[key] {
					seen[key] = true
					missing = append(missing, closureRequirement{ticket.Name, entry, path, kind})
				}
			}
		}
	}
	return missing
}

// fixtureClosureCheck grades the red-capable fixture into the ticket. A ticket
// that edits a fixture-pinned line without naming the owning fixture directory
// leaves the proof outside the charge, and the bite breaks unnoticed.
func fixtureClosureCheck(f Facts) CheckResult {
	unnamed := closureMessages(missingClosures(f, fixtureClosure))
	if len(unnamed) > 0 {
		return red("fixture-closure", "Writes: entry names a fixture-pinned path without naming the fixture: "+strings.Join(unnamed, ", "))
	}
	return green("fixture-closure")
}

// registryClosureCheck grades the declared binding into the ticket. A ticket that
// writes a bound package and omits a bound registry finds that registry mid-build
// and pays a repair round.
func registryClosureCheck(f Facts) CheckResult {
	omitted := closureMessages(missingClosures(f, registryClosure))
	if len(omitted) > 0 {
		return red("registry-closure", "Writes: entry names a bound package without naming every bound file: "+strings.Join(omitted, ", "))
	}
	return green("registry-closure")
}

func closureMessages(requirements []closureRequirement) []string {
	result := make([]string, len(requirements))
	for i, requirement := range requirements {
		verb := "requires"
		if requirement.kind == fixtureClosure {
			verb = "is pinned by"
		}
		result[i] = requirement.ticket + ": " + requirement.entry + " " + verb + " " + requirement.path
	}
	return result
}
