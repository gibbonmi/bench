package compatibility

// capability ties a workflow requirement to its actual-interface probe.
type capability struct {
	name, class, action string
	operations          []string
}

func capabilities() []capability {
	return []capability{
		{"hook-behavior", "hooks", "observe the installed hook in the selected interface; a declaration or hook-process success is not actual-tool proof", []string{"work", "review", "recover"}},
		{"normal-shell", "commands", "run pwd through the actual chat shell with normal permissions", []string{"diagnose", "work", "review", "recover"}},
		{"repository-wrapper", "commands", "invoke the repository Bench wrapper through that same normal-permission chat tool", []string{"diagnose", "work", "review", "recover"}},
		{"file-access", "files", "write and read expected scratch bytes through the actual interface before dependent file work", []string{"work", "review", "recover"}},
		{"repository-rules", "rules", "read the repository agreement and its shared Bench rules in this independent chat", []string{"diagnose", "work", "review", "recover"}},
		{"repository-skill", "skills", "invoke the required repository skill from its installed path in this interface", []string{"work", "review", "recover"}},
		{"permission-policy", "permissions", "observe the required effective permission behavior without broadening policy", []string{"diagnose", "work", "review", "recover"}},
		{"worktree-isolation", "worktrees", "resolve this writer's distinct Bench assignment before its first write", []string{"work"}},
		{"review-outcome", "reviews", "complete the selected review operation through its actual route; an equivalent counts only after that operation succeeds", []string{"review"}},
		{"failed-interface-retest", "recovery", "retest the failed capability through the previously failed interface with normal permissions after repair", []string{"recover"}},
		{"desktop-presentation", "reviews", "use the requested desktop presentation tool only when the operation requires it and no observed equivalent completes that operation", []string{"present"}},
	}
}

func required(item capability, operation string) bool {
	if operation == "workflow" {
		return true
	}
	for _, candidate := range item.operations {
		if candidate == operation {
			return true
		}
	}
	return false
}

// CapabilityAction returns the probe owned by the capability inventory.
func CapabilityAction(name string) string {
	for _, item := range capabilities() {
		if item.name == name {
			return item.action
		}
	}
	return "unknown capability; identify the dependent operation before work"
}

// LiveObligations supplies the pending probes for one selected workflow operation.
func LiveObligations(operation string, includeHook bool) []LiveRow {
	var rows []LiveRow
	for _, item := range capabilities() {
		if required(item, operation) && (includeHook || item.name != "hook-behavior") {
			rows = append(rows, LiveRow{Capability: item.name, Action: item.action})
		}
	}
	return rows
}
