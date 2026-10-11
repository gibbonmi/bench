package spec

// Metadata returns the Status: and Roadmap: values of spec content. It applies the parser of
// Facts, so spec bytes that a caller reads at a commit give the same values.
func Metadata(content []byte) (status, roadmapID string) { return metadata(content) }
