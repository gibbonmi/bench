package landing

// Deliverable is the reviewed deliverable that the landing closes: the staged spec, the
// tickets-only folder, or nothing on a spec-less landing.
func (r ReviewedRequest) Deliverable() string {
	if r.SpecPath != "" {
		return r.SpecPath
	}
	return r.ClosePath
}
