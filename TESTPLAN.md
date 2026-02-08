# Test Plan

## Goals

- Prevent regressions in time-aware allowlisting semantics:
  - Developers are allowed only during their active windows.
  - Keys are allowed only during their active windows.
- Ensure YAML sync behavior matches intended modes:
  - `reset`: overwrite DB from config (fast, simple).
  - `reconcile`: preserve history by revoking/removing instead of deleting.
- Keep `go test ./...` fast and deterministic; keep “real world” toolchain coverage in the existing GitHub Actions self-test.

## What We Test

- `db` (unit-ish, DB-backed)
  - Developer lifecycle: add, remove, re-add; boundary behavior at timestamps.
  - Key lifecycle: add, revoke, re-add; verifies “historical” correctness for old commits after re-add.
  - YAML sync:
    - `SyncFromYAMLReset` inserts “active now” entries with epoch timestamps.
    - `SyncFromYAMLReconcile` revokes removed keys and removes missing developers while preserving history.

- `git` (integration, lightweight)
  - Uses a temporary git repo to test helpers like `CommitEmail`, `CommitTimestamp`, `CommitsInRange`, `RootCommit`.
  - Signature helpers are only tested for the unsigned/failure path here.

- End-to-end (workflow)
  - `.github/workflows/validtr-selftest.yml` remains the authoritative E2E check (git + gpg + CLI + policy behavior).

## How To Run

```sh
go test ./...
```

Notes:
- `db` tests require CGO because this project uses `github.com/mattn/go-sqlite3`.
- `git` tests require the `git` executable; they will `t.Skip()` if it’s not available.

