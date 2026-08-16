# Publish the complete hostile-input release note

Blocked by: harden-every-skill-file-reader.md, enumerate-skill-directories-literally.md, diagnose-orphans-before-adapter-suppression.md, harden-every-reference-file-reader.md, validate-payload-once-for-every-consumer.md, require-complete-leading-frontmatter.md, refuse-control-runes-before-rendering.md, require-exactly-one-marker-span.md, clean-up-failed-replacements.md, clean-up-sigint-before-replacement.md, distinguish-missing-git-at-dispatch.md
Writes: CHANGELOG.md

## What to build

Add one concise typed entry under the changelog's Unreleased section describing the
complete user-visible skills-index hostile-input refusal, safe replacement cleanup,
and missing-Git recovery outcome. Every HI1-HI14 behavior-bearing sibling blocks this
ticket, so the note cannot claim partially landed hardening.

This is an explicitly user-facing completion unit, not a behavior implementation
dependency: it changes only the shipped release-note surface after all observable
behavior is green, and no behavior ticket waits on it.

## Acceptance

- [ ] The Unreleased changelog contains exactly one concise typed entry that describes
  the complete skills-index hostile-input hardening and operator recovery outcome.
- [ ] The entry lands only after every named behavior blocker and neither duplicates
  implementation detail nor advertises an unlanded HI1-HI14 behavior.
