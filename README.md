# gotree

A simple Go implementation of the Unix `tree` command.
It recursively prints a directory structure in a readable ASCII format.

---

## Features

- Recursive directory traversal
- Deterministic sorted output
- Optional hidden file support
- Depth limiting
- Clean separation of scanner and printer logic

---

## Installtion

Run from the project root:

```bash
# Installation process
git clone https://github.com/Triangulation5/gotree.git gotree
cd gotree

# Running the application
go run .\cmd\tree\main.go
```
## Flags

| Flag | Description                   |
|------|-------------------------------|
| `-a` | Show hidden files (dotfiles)  |
| `-L` | Maximum depth (0 = unlimited) |

---

## Examples

Show hidden files:

```bash
go run .\cmd\tree\main.go -a .
```

Limit depth to 2:

```bash
go run .\cmd\tree\main.go -L 2 .
```
