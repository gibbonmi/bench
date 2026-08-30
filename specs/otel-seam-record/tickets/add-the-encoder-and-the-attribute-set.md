# 3. Add the record package with its encoder and attribute set

Blocked by: add-the-pinned-otel-dependencies.md
Line: opus / medium
Rows: OT1, OT2, OT3, OT4, OT22
Writes: internal/otelrecord/ (new), internal/otelrecord/encode.go (new), internal/otelrecord/encode_test.go (new), internal/otelrecord/attributes.go (new)

## What to build

The package `internal/otelrecord` is new. It holds the one encoder that turns a
span into a single OTLP-JSON line. The line parses as one JSON object under the
key `resourceSpans`, so a collector receiver ingests the record.

The encoder holds the four OTLP-JSON deviations. It writes each trace id and
each span id as a lowercase hex string. It writes the span kind and the status
code as integers. It writes `startTimeUnixNano` and `endTimeUnixNano` as quoted
decimal strings. It writes every field name in lowerCamelCase.

The encoder is hand-written, and it mirrors the struct schema of the upstream
internal `otlpjson` package. Every official encoder needs
`google.golang.org/protobuf`, and the reviewer excluded that footprint.

The package also declares the attribute set that every later ticket reads: the
seam name, the subject id, the outcome, and the measures. Row OT22 is
review-owned, and this declaration is the artifact the review grades each new
attribute against. No mechanical check enforces it.

## Acceptance

- [ ] OT1: the encoder returns one line, and that line parses as a JSON object with the key `resourceSpans`.
- [ ] OT2: the encoded fixture span carries its trace id and its span id as lowercase hex strings.
- [ ] OT3: the encoded fixture span carries the span kind and the status code as integers.
- [ ] OT4: the encoded fixture span carries `startTimeUnixNano` as a quoted decimal string.
- [ ] the gate `test` phase stays green for `internal/otelrecord`.
