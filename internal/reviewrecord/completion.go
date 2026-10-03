package reviewrecord

func completionReconciler(record Record, plan Plan) string {
	if plan.Delegated() {
		return plan.Execution.OrchestratorSession
	}
	return record.ImplementationSession
}

// RecordCompletion derives and writes the completed entry for source. The current
// final verification results stay in the entry, and the final checkpoint validates
// the complete candidate before the transaction replaces the record file.
func RecordCompletion(root, spec, source string) (Completion, error) {
	digest, plan, err := atSource(root, spec, source)
	if err != nil {
		return Completion{}, err
	}
	tree, err := sourceTree(root, source)
	if err != nil {
		return Completion{}, err
	}
	var entry Completion
	err = write(root, spec, nil, func(record *Record) error {
		reconciliation := map[string]string{}
		for _, chunk := range plan.Chunks {
			for _, row := range chunk.Rows {
				reconciliation[row] = "covered"
			}
		}
		entry = Completion{
			State:          "completed",
			SourceDigest:   digest,
			Performer:      completionReconciler(*record, plan),
			Reconciliation: reconciliation,
			Verification:   append([]Verification{}, record.Completion.Verification...),
		}
		candidate := *record
		candidate.Completion = entry
		if err := checkSource(root, tree, source, candidate, "", true, true); err != nil {
			return err
		}
		record.Completion = entry
		return nil
	})
	return entry, err
}
