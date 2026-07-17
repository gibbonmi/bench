# npm publication research

Researched 2026-07-16 for FT83 against npm CLI 11 documentation and read-only
registry/CLI probes. This records registry facts; the reviewer-owned release policy
is ticket #6 in the parent map.

## Findings

1. A published `package@version` is immutable and can never be reused, including
   after unpublish. `npm publish` stores both SHA-1 and SHA-512 SRI for the uploaded
   tarball and rejects an existing name/version. A retry may therefore skip an
   existing version only when its registry integrity equals the locally approved
   tarball; a mismatch must stop the release. [`npm publish`](https://docs.npmjs.com/cli/v11/commands/npm-publish/),
   [npm unpublish policy](https://docs.npmjs.com/policies/unpublish/)

2. `npm stage publish` creates a reviewable package that is not live until a
   maintainer approves it with 2FA. The staged tarball can be viewed and downloaded,
   and its tag is fixed at staging time. Staging requires npm 11.15.0+, Node 22.14.0+,
   write access, an existing registry package, and 2FA; it cannot create a brand-new
   package identity. [`npm stage`](https://docs.npmjs.com/cli/v11/commands/npm-stage/),
   [staged publishing](https://docs.npmjs.com/staged-publishing/)

3. A CI trusted publisher can be restricted to stage-only OIDC permission. OIDC can
   submit a stage but cannot approve, reject, view, list, or download it; approval is
   deliberately interactive proof-of-presence. Each package accepts only one trusted
   publisher configuration. [`npm trusted publishers`](https://docs.npmjs.com/trusted-publishers/)

4. Direct publish assigns `latest` unless `--tag` supplies a non-default tag. Tags
   can later be added or removed, so promotion can move `latest` only after all live
   packages verify. A non-default tag prevents default installation but does not make
   the version private: an exact version or that tag remains installable.
   [`npm publish`](https://docs.npmjs.com/cli/v11/commands/npm-publish/),
   [`npm dist-tag`](https://docs.npmjs.com/cli/v11/commands/npm-dist-tag/)

5. npm has no atomic multi-package publish or approval operation. This is an
   inference from the documented package-scoped publish, stage-ID approval, and
   dist-tag commands. Bench must sequence platform packages before the wrapper and
   verify each registry integrity before advancing.

6. Unpublish is not rollback: eligibility is conditional, deletion is irreversible,
   and the name/version remains burned. Once a version is live, rollback means remove
   or restore its distribution tags, deprecate the version with a replacement or
   recovery message, and publish a new version. Automation must never rely on
   unpublish. [`npm unpublish policy`](https://docs.npmjs.com/policies/unpublish/),
   [npm deprecation](https://docs.npmjs.com/deprecating-and-undeprecating-packages-or-package-versions/)

7. `npm view <package>@<version> dist.integrity dist.shasum --json` is a read-only
   registry query, while `npm pack --dry-run --json` exposes the prospective file
   inventory and local integrity without publishing. An E404 cannot distinguish an
   available public name from an inaccessible identity and therefore cannot prove
   publish authority. [`npm view`](https://docs.npmjs.com/cli/v11/commands/npm-view/),
   [`npm pack`](https://docs.npmjs.com/cli/v11/commands/npm-pack/)

## Consequence for the release design

There are two paths:

- **First publication:** because none of the five identities exists and npm cannot
  stage a new identity, direct-publish platform packages under one version-specific
  non-default tag, verify every live integrity, direct-publish and verify the wrapper
  last, then promote `latest`. A partial retry treats matching live bytes as complete
  and any mismatch as terminal. Name ownership and publish authority remain an
  explicit manual precondition.
- **Subsequent releases:** pin npm 11.15+ and stage every package under a
  version-specific non-default tag using stage-only trusted publishing. A maintainer
  downloads/reviews the staged bytes, approves platform packages first, Bench waits
  until their registry integrities match, then the maintainer approves the wrapper.
  Bench verifies it and promotes `latest` only after the complete set agrees.

Before wrapper approval, rejected stages are safe cleanup. After any approval, failed
releases use tag restoration/removal plus deprecation; they never unpublish.

## Read-only probes

The local environment reported Node `v25.8.1`, npm `11.11.0`, and `Unknown command:
"stage"`, consistent with staging's npm 11.15 minimum. `npm pack --dry-run --json
--ignore-scripts .` emitted the current package's file inventory, SHA-1, and SHA-512
SRI without a tarball. An unauthenticated query of `npm@11.18.0` returned both
registry digests. All five intended `0.2.0` identities (`redbench` and the four
`@redbench/*` platform packages) returned E404; that is absence from the public view,
not proof of availability or authority.

Repeat the registry probe without ambient npm credentials:

```sh
NPM_CONFIG_USERCONFIG=/dev/null npm view \
  '<package>@<version>' version dist.integrity dist.shasum --json
```

Compare an already-built tarball byte-for-byte with the registry's SHA-512 SRI:

```sh
local_integrity="$(node -e '
  const crypto = require("crypto"), fs = require("fs");
  const bytes = fs.readFileSync(process.argv[1]);
  process.stdout.write("sha512-" + crypto.createHash("sha512").update(bytes).digest("base64"));
' './artifact.tgz')"
remote_integrity="$(NPM_CONFIG_USERCONFIG=/dev/null npm view \
  '<package>@<version>' dist.integrity)"
test "$local_integrity" = "$remote_integrity"
```
