package chargeevidence

// ChargeRow is one charge descriptor: one build row, or one row per review axis.
type ChargeRow struct {
	Axis, Ticket, Access string
}

// SharedRow binds one generated review evidence kind to its source.
type SharedRow struct {
	Kind, Source string
}

// CompletionRow carries the existing review completion facts.
type CompletionRow struct {
	Record, SourceDigest, PlanDigest, RecordState, Detail string
}

// Metadata is the typed body of the required metadata source. Every source cell holds a
// manifest source identifier, and each repeated list has one row per member.
type Metadata struct {
	Charge     []ChargeRow
	Fence      []string
	Writes     []string
	Coverage   []string
	Checks     []string
	Returns    []string
	Shared     []SharedRow
	Completion []CompletionRow
}

// EncodeMetadata renders the canonical metadata source body.
func EncodeMetadata(m Metadata) ([]byte, error) {
	values := map[string][][]any{
		blockFence:    column(m.Fence),
		blockWrites:   column(m.Writes),
		blockCoverage: column(m.Coverage),
		blockChecks:   column(m.Checks),
		blockReturns:  column(m.Returns),
	}
	for _, r := range m.Charge {
		values[blockCharge] = append(values[blockCharge], []any{r.Axis, r.Ticket, r.Access})
	}
	for _, r := range m.Shared {
		values[blockShared] = append(values[blockShared], []any{r.Kind, r.Source})
	}
	for _, r := range m.Completion {
		values[blockCompletion] = append(values[blockCompletion], []any{r.Record, r.SourceDigest, r.PlanDigest, r.RecordState, r.Detail})
	}
	return encodeBlocks(MetadataBlocks, values)
}

// DecodeMetadata strictly decodes a canonical metadata source body.
func DecodeMetadata(data []byte) (Metadata, error) {
	values, err := decodeBlocks(MetadataBlocks, data)
	if err != nil {
		return Metadata{}, err
	}
	m := Metadata{
		Fence:    texts(values[blockFence]),
		Writes:   texts(values[blockWrites]),
		Coverage: texts(values[blockCoverage]),
		Checks:   texts(values[blockChecks]),
		Returns:  texts(values[blockReturns]),
	}
	for _, r := range values[blockCharge] {
		m.Charge = append(m.Charge, ChargeRow{r[0].(string), r[1].(string), r[2].(string)})
	}
	for _, r := range values[blockShared] {
		m.Shared = append(m.Shared, SharedRow{r[0].(string), r[1].(string)})
	}
	for _, r := range values[blockCompletion] {
		m.Completion = append(m.Completion, CompletionRow{r[0].(string), r[1].(string), r[2].(string), r[3].(string), r[4].(string)})
	}
	return m, nil
}

// references lists every source identifier the metadata names. An empty cell names no
// source: a review charge row selects no ticket, and the reader must not read that blank
// as a reference to a source the pack never declared.
func (m Metadata) references() []string {
	var refs []string
	add := func(values ...string) {
		for _, value := range values {
			if value != "" {
				refs = append(refs, value)
			}
		}
	}
	add(m.Checks...)
	add(m.Returns...)
	for _, r := range m.Charge {
		add(r.Ticket)
	}
	for _, r := range m.Shared {
		add(r.Source)
	}
	return refs
}

func column(values []string) [][]any {
	rows := make([][]any, len(values))
	for i, value := range values {
		rows[i] = []any{value}
	}
	return rows
}

func texts(rows [][]any) []string {
	values := make([]string, len(rows))
	for i, row := range rows {
		values[i] = row[0].(string)
	}
	return values
}
