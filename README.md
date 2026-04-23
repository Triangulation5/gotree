# gotree (gtr)

<!--toc:start-->
- [gotree (gtr)](#gotree-gtr)
  - [Features](#features)
  - [Installtion](#installtion)
  - [Flags](#flags)
  - [Examples](#examples)
<!--toc:end-->

A simple Go implementation of the Unix `tree` command.
It recursively prints a directory structure in a readable ASCII format.

---

## Features

- Recursive directory traversal
- Deterministic sorted output
- Optional hidden file support
- Depth limiting
- Clean separation of scanner and printer logic
- **Icon sets**: Support for Nerd Font, Unicode, and ASCII icons, with runtime switching via `--icons` or `-i` flags.
- **Summary mode**: Collapses deep structures, aggregates ignored/hidden files, suppresses noise, and respects `.gitignore`.
- **Sort by extension**: Files can be sorted by their extension, with directories always appearing first, followed by extension groups, and then files without extensions.

---

## Installtion

Run from the project root:

```bash
# Installation process
git clone https://github.com/Triangulation5/gotree.git gotree
cd gotree

# Running the application
go run .\src\cmd	ree\main.go
```

Building a reusable executable:

```bash
go build -o gtr.exe .\src\cmd	ree\main.go
```

Link this executable to a command in your .bashrc or Microsoft.Powershell_profile.ps1

## Flags

| Flag                         | Description                                      |
|------------------------------|--------------------------------------------------|
| `-a`                         | Show hidden files (dotfiles)                     |
| `-L`                         | Maximum depth (0 = unlimited)                    |
| `-o`, `--order-by-extension` | Sort files by extension                          |
| `-m`, `--summary`            | Summary mode (collapses large/noise directories) |
| `--icons <mode>`             | Icon mode: nerd (default), unicode, ascii, none  |
| `-i <mode>`                  | Alias for --icons mode                           |

---

## Examples

Show hidden files:

```bash
go run .\src\cmd	ree\main.go -a .
```

Limit depth to 2:

```bash
go run .\src\cmd	ree\main.go -L 2 .
```

Sort by extension:

```bash
gotree -o .
```

Summary mode:

```bash
gotree -m .
```

Nerd icons (default):

```bash
gotree .
```

Unicode icons:

```bash
gotree --icons=unicode .
gotree -i unicode .
```

No icons:

```bash
gotree --icons=none .
gotree -i none .
