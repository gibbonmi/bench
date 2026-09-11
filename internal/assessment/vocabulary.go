package assessment

// Roles is the versioned performer-role vocabulary.
func Roles() []string {
	return []string{"implementation", "repair", "verification", "review", "diagnostic"}
}

// States is the versioned run and attempt state vocabulary.
func States() []string { return []string{"running", "succeeded", "failed", "cancelled", "incomplete"} }

const (
	// DeltaMode describes quantities contributed by a single event.
	DeltaMode = "delta"
	// CumulativeMode describes snapshots from an epoch's zero baseline.
	CumulativeMode = "cumulative"
)
