package assessment

// Roles is the versioned performer-role vocabulary. Orchestration is last, so
// the existing roles keep their positions in every projection that indexes this
// list.
func Roles() []string {
	return []string{"implementation", "repair", "verification", "review", "diagnostic", "orchestration"}
}

// States is the versioned run and attempt state vocabulary.
func States() []string { return []string{"running", "succeeded", "failed", "cancelled", "incomplete"} }

const (
	// DeltaMode describes quantities contributed by a single event.
	DeltaMode = "delta"
	// CumulativeMode describes snapshots from an epoch's zero baseline.
	CumulativeMode = "cumulative"
)
