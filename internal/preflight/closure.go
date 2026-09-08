package preflight

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
