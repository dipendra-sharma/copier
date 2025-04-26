# Copier

A Go CLI tool to recursively copy a folder (with all subfolders and files) to another path, skipping common build/config directories, logging skipped items and errors.

## Usage

```
copier <source_path> <destination_path> [--log <logfile>]
```

- Skips: `.git`, `node_modules`, `build`, `dist`, `.cache`, `.idea`, `.vscode`, `target`, `venv`, `__pycache__`, `.DS_Store`, etc.
- Logs skipped directories/files and errors to `copy.log` (or a custom log file).
- Preserves permissions and copies symlinks as links.

## Build

```
go build -o copier ./cmd/copier
```

## Test

```
go test ./...
```

---

See `test/copier_test.go` for comprehensive test coverage.
