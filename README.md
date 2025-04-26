# Copier

A robust Go CLI tool to recursively copy a folder (with all subfolders and files) to another path, with advanced skipping, logging, and cross-platform support (macOS, Linux, Windows).

## Features
- **Recursively copies** directories and files, preserving permissions and symlinks.
- **Configurable skipping** using a `.copyignore` file (supports glob patterns, like `.gitignore`).
- **Default skips** for common build, cache, and dependency directories (see below).
- **Logs** all skipped paths and errors to a log file (default: `copy.log`).
- **Continues** on error (does not abort entire copy).
- **CLI flags** for custom log and ignore file paths.

## Usage

```sh
copier <source_path> <destination_path> [--log <logfile>] [--ignore <ignorefile>]
```

- `--log <logfile>`: Path to log file (default: `copy.log` in current directory)
- `--ignore <ignorefile>`: Path to ignore file (default: `.copyignore` in source directory)

### Example
```sh
go build -o copier ./cmd/copier
./copier ~/myproject /Volumes/backup/myproject --log backup.log --ignore .copyignore
```

## Example `.copyignore`
```
.git/
build/
dist/
.cache/
node_modules/
__pycache__/
.venv/
venv/
.mypy_cache/
.pytest_cache/
```
Supports glob patterns (wildcards), e.g. `*.log`, `temp*/`, etc.

## Logging
- All skipped paths and errors are written to the log file for audit and troubleshooting.
- Skips are logged whether from `.copyignore` or default rules.

## Running Tests
```sh
go test -v ./...
```

## Contributing
Pull requests and issues are welcome! Please ensure new features are covered by tests.

---
MIT License. Built with care by [Dipendra Sharma](https://github.com/dipendra-sharma).
