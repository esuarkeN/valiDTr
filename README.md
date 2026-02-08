# valiDTr — Git commit verification CLI

**valiDTr** verifies Git commits by checking:

1. The commit is GPG-signed and the signature is valid.
2. The signer key and developer email are allowlisted under a time-aware policy.

This is intended for CI enforcement (PRs + pushes) and for local verification.

## ✅ GitHub Actions: quick setup (recommended)

This section walks through the setup to implement the valiDTr workflow.

### 1) Add the allowlist config

Create `.validtr/config.yml` in your repo and list allowed developers and their key IDs:

```yaml
developers:
  - email: dev@example.com
    name: Dev Example
    keys:
      - id: 1234ABCD5678EFGH
      - id: 90AB12CD34EF56GH
```

Notes:
- `email` is required; `name` is optional.
- `keys[].id` must match the output of `git show -s --format=%GK <commit>`.
- Removing a developer or key from the config and running `sync --mode=reconcile` will mark it removed/revoked (history is preserved).

### 2) Add public keys

Export each developer public key and place them under `.validtr/pubkeys/`:

```sh
gpg --armor --export <KEY_FPR_OR_ID> > .validtr/pubkeys/dev1.asc
```

Directory layout:

```text
.validtr/
  config.yml
  pubkeys/
    dev1.asc
    dev2.asc
```

### 3) Add the workflow

Create `.github/workflows/validtr.yml`:

```yaml
name: valiDTr

on:
  pull_request:
  push:
    branches: ["main"]

permissions:
  contents: read

jobs:
  validtr:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6.0.2
        with:
          fetch-depth: 0

      - uses: esuarkeN/valiDTr/.github/actions/validtr@v1
        with:
          config: .validtr/config.yml
          pubkeys: .validtr/pubkeys/*.asc
          policy: current
          email_mode: committer
```

### 4) Choose policy and email mode

- **policy**
  - `current`: evaluate the allowlist *now* (strict enforcement)
  - `historical`: evaluate the allowlist at commit time (preserves older commits for audits; commit timestamps are forgeable)
- **email_mode**
  - `committer` (default): checks committer email
  - `author`: checks author email

Security note:
- In the GitHub Action, PR/push admission checks always use `current` (even if `historical` is requested) to avoid commit-date backdating bypasses.

You can also control these via env vars: `VALIDTR_POLICY`, `VALIDTR_EMAIL_MODE`.

## ✅ How it works (overview)

valiDTr verifies commits in two steps:

1. **GPG signature check**: The commit must be signed and the signature must be valid.
2. **Policy check**: The developer email must be active and the signing key must be active for that developer under the selected policy.

Internally, valiDTr keeps a **SQLite database** of developers/keys with timestamps. The DB is typically built from YAML during CI (`valiDTr sync`).

## 🔒 Workflow safety (no secret leakage)

- The workflow only uses `.validtr/config.yml` and `.validtr/pubkeys/*.asc` (public keys).
- For `pull_request` events, the action reads allowlist files from the PR **base commit**, not from PR-modified files.
- No private keys or secrets are required.
- The GitHub Action uses a temporary DB + `GNUPGHOME` on the runner and does not upload artifacts by default.

## 🛠️ CLI

All DB-backed commands accept:
- `--db <path>` (or `VALIDTR_DB`) to choose the SQLite DB path.
- `--policy current|historical` (or `VALIDTR_POLICY`) for verification.
- `--email-mode committer|author` (or `VALIDTR_EMAIL_MODE`) for verification.

Commands:

- `valiDTr sync --config <path> --mode reset|reconcile`
  - `reset`: overwrite DB from config (fast/simple).
  - `reconcile`: preserve history by revoking/removing instead of deleting.
- `valiDTr verify <commit>`
- `valiDTr verify-range <from> <to>`
- `valiDTr add-dev <email> [--name <name>] [--since now|epoch|RFC3339]`
- `valiDTr remove-dev <email>`
- `valiDTr add-key <developer-email> <key-id> [--since now|epoch|RFC3339]`
- `valiDTr revoke-key <developer-email> <key-id>`
- `valiDTr list-devs`
- `valiDTr list-keys [--email <developer-email>]`
- `valiDTr version`

## 📌 Installation

### Go

```sh
go install github.com/esuarkeN/valiDTr@latest
```

Notes:
- This project uses `github.com/mattn/go-sqlite3`, so **CGO + a C compiler** are required.
- `git` and `gpg` must be available for signature verification.

### Build from source

```sh
go test ./...
go build -o valiDTr .
```

## Contributors

esuarkeN (Eric Krause)
