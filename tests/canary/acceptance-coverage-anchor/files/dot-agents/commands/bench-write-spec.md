# /bench-write-spec fixture

This canary fixture intentionally omits the feature-build coverage map anchor. The
require_anchor gate must fire when this command no longer names the map it is
supposed to write.

It still carries the other two anchors so only the missing one bites:

- red signal
- why it catches the failure
