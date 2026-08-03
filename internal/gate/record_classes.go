package gate

// storeRecordClasses is every record class in the evidence store that declares an exact
// field set, keyed by the name a failure reports. The classes are alternatives rather than a
// spectrum, so the family is enumerated here once and each member is measured against every
// sibling — a class added to the store joins that comparison without a second edit.
//
// The ready verdict classes (full, reduced, partial) are folded in from readyFieldClasses in
// verdict.go rather than restated, so this map cannot drift from the set verdict.go actually
// validates against. The enumeration spans verdicts, component slots, the conformance-check
// ledger, and attestations equally, so it lives outside every individual class owner.
var storeRecordClasses = func() map[string][]string {
	classes := make(map[string][]string, len(readyFieldClasses)+3)
	for name, fields := range readyFieldClasses {
		classes[name] = fields
	}
	classes["component slot"] = componentSlotFields
	classes["build attestation"] = buildAttestationFields
	classes["conformance check slot store"] = conformanceCheckSlotStoreFields
	return classes
}()

// verdictClasses names the members of storeRecordClasses that answer for the whole tree,
// derived from readyFieldClasses' keys rather than restated — every class registered there is
// a verdict class by definition, and none of storeRecordClasses' other members come from that
// map. The verdict sets deliberately nest — reduced and partial are each the full set plus
// their own inherited fields — so they are the family's one designed overlap, and every other
// class has to stay disjoint from all three.
var verdictClasses = func() []string {
	names := make([]string, 0, len(readyFieldClasses))
	for name := range readyFieldClasses {
		names = append(names, name)
	}
	return names
}()

// recordClassSharesVerdictFields returns the names fields shares with any verdict class
// beyond the store-wide shared ones. It is what makes "this class carries no verdict field"
// a checked property rather than a claim in a comment, enumerated from the verdict classes
// themselves so a field added there is covered.
func recordClassSharesVerdictFields(fields []string) []string {
	var shared []string
	for _, name := range fields {
		if contains(recordSharedFields, name) {
			continue
		}
		for _, class := range verdictClasses {
			if contains(storeRecordClasses[class], name) {
				shared = append(shared, name)
				break
			}
		}
	}
	return shared
}
