# gotree (gtr)

<!--toc:start-->
- [gotree (gtr)](#gotree-gtr)
  - [Release](#release)
  - [Features](#features)
  - [Run and Build](#run-and-build)
  - [Flags](#flags)
  - [Examples](#examples)
  - [Testing](#testing)
    - [Go tests](#go-tests)
    - [Python CLI integration tests](#python-cli-integration-tests)
  - [Packaging](#packaging)
  - [Production Readiness](#production-readiness)
<!--toc:end-->

`gotree` is a Go implementation of the Unix `tree` command with filtering, sorting, summaries, icon modes, and structured output.

## Release

Current release: **v2.0.0** (`gotree-v2.0.0`).

## Features

- Recursive directory traversal
- Hidden file support (`-a`, `--all`)
- Depth limiting (`-L`)
- Directory-only mode (`-d`, `--directories`)
- Include/ignore glob filtering (`--include`, `--ignore`), where include matches override ignore matches
- Sorting (`--sort name|ext|size|mtime`, `--reverse`)
- Backward-compatible aliases (`-s`, `--sorted`, `-o`, `--order-by-extension`)
- Recursive size rollups (`--du`)
- Icon sets (`--icons`, `-i`) and color theme control (`--theme auto|color|mono`)
- Summary mode (`-m`, `--summary`)
- Structured output (`--json`, `--yaml`)
- Shell completion generation (`--completion bash|zsh|fish|powershell`)

## Run and Build

From project root:

```bash
go run ./src/cmd/gotree
```

Build executable:

```bash
go build -o gotree ./src/cmd/gotree
```

## Flags

| Flag | Description |
|---|---|
| `-a`, `--all` | Show hidden files |
| `-L` | Maximum depth (0 = unlimited) |
| `-d`, `--directories` | Show directories only |
| `--include <glob>` | Include glob pattern (repeatable; comma-separated supported) |
| `--ignore <glob>` | Ignore glob pattern (repeatable; comma-separated supported) |
| `--sort <name\|ext\|size\|mtime>` | Sort mode |
| `--reverse` | Reverse sort order |
| `-s`, `--sorted` | Alias for `--sort name` |
| `-o`, `--order-by-extension` | Alias for `--sort ext` |
| `--du` | Show recursive size totals in tree and summary |
| `--theme <auto\|color\|mono>` | Color behavior (`auto` respects `NO_COLOR`) |
| `-i`, `--icons <mode>` | Icon mode: `nerd`, `unicode`, `ascii`, `none` |
| `-m`, `--summary` | Summary mode |
| `--json` | Output JSON payload |
| `--yaml` | Output YAML payload |
| `--completion <shell>` | Print completion script |
| `-v`, `--version` | Show version |

## Examples

```bash
# Default tree
gotree .

# Include only Go files
gotree --include "*.go" .

# Ignore build output
gotree --ignore "dist/*" .

# Sort by size descending
gotree --sort size --reverse .

# Structured output
gotree --json .
gotree --yaml .

# No ANSI colors
gotree --theme mono .
```

## Testing

### Go tests

```bash
go test ./src/...
go vet ./src/...
```

### Python CLI integration tests

Install dependencies:

```bash
python -m pip install -r ./tests/python/requirements.txt
```

Run tests:

```bash
python -m pytest ./tests/python
```

If your workspace is in OneDrive and pytest temp cleanup fails on Windows, run:

```bash
python -m pytest ./tests/python --basetemp C:/Users/<you>/.codex/memories/pytest-tmp
```

The Python suite runs the CLI as a black box (`go run ./src/cmd/gotree ...`) and validates behavior like structured output, filtering precedence, and completion output.

## Packaging

Build release artifacts into `dist/`:

```bash
# Unix-like shells
bash ./scripts/release.sh

# PowerShell
./scripts/release.ps1
```

## Production Readiness

For this release, readiness means all of the following pass:

- `go build ./src/cmd/gotree`
- `go test ./src/...`
- `go vet ./src/...`
- CLI smoke checks for text output, JSON/YAML output, completion generation, and filter/sort behavior
