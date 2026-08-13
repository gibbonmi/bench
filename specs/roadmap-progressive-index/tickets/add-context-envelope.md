# Add the terminal help envelope to every --context mode

Blocked by: none
Writes: internal/roadmap

## What to build

Every `bench roadmap --context` invocation (default and `--full`) emits one
TOON document ending with the terminal `help[N]{cmd,why}` block, so the whole
stdout decodes as a single document. No conformance pin binds the context
block set yet, so this lands green alone; the schema advertisement stays 3
until land-index-doctrine.md moves it with the doctrine anchors. Disclosure
content here is minimal (the block may be empty); the exact
request-complete-value wording lands with join-axi-set.md.

## Acceptance

- [ ] `--context` and `--context --full` stdout each decode as one TOON
      document whose final block is a schema-correct `help` table
      (covers PI16).
- [ ] `--context` stdout carries no bytes outside its TOON blocks — no
      preamble, separator, or trailing line — so the complete output is the
      document, not a document embedded in noise.
