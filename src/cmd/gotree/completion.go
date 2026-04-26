package main

import (
	"fmt"
	"strings"
)

var completionFlags = []string{
	"-h", "--help",
	"-v", "--version",
	"-a", "--all",
	"-L",
	"-d", "--directories",
	"-s", "--sorted",
	"-o", "--order-by-extension",
	"-m", "--summary",
	"--sort",
	"--reverse",
	"--include",
	"--ignore",
	"--du",
	"--json",
	"--yaml",
	"--theme",
	"-i", "--icons",
	"--completion",
}

func completionScript(shell string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(shell)) {
	case "bash":
		return bashCompletion(), nil
	case "zsh":
		return zshCompletion(), nil
	case "fish":
		return fishCompletion(), nil
	case "powershell":
		return powershellCompletion(), nil
	default:
		return "", fmt.Errorf("unsupported shell %q (expected: bash, zsh, fish, powershell)", shell)
	}
}

func bashCompletion() string {
	return `# bash completion for gotree
_gotree_completion() {
  local cur prev opts
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"
  opts="` + strings.Join(completionFlags, " ") + `"

  case "${prev}" in
    --sort) COMPREPLY=( $(compgen -W "name ext size mtime" -- "${cur}") ); return 0 ;;
    --theme) COMPREPLY=( $(compgen -W "auto color mono" -- "${cur}") ); return 0 ;;
    --icons|-i) COMPREPLY=( $(compgen -W "nerd unicode ascii none" -- "${cur}") ); return 0 ;;
    --completion) COMPREPLY=( $(compgen -W "bash zsh fish powershell" -- "${cur}") ); return 0 ;;
  esac

  COMPREPLY=( $(compgen -W "${opts}" -- "${cur}") )
  return 0
}
complete -F _gotree_completion gotree
complete -F _gotree_completion gtr
`
}

func zshCompletion() string {
	return `#compdef gotree gtr
_gotree_completion() {
  _arguments \
    '(-h --help)'{-h,--help}'[show help]' \
    '(-v --version)'{-v,--version}'[show version]' \
    '(-a --all)'{-a,--all}'[show hidden files]' \
    '(-L)'{-L}'[max depth]:depth:' \
    '(-d --directories)'{-d,--directories}'[list directories only]' \
    '(-s --sorted)'{-s,--sorted}'[sort by name]' \
    '(-o --order-by-extension)'{-o,--order-by-extension}'[sort by extension]' \
    '(-m --summary)'{-m,--summary}'[summary mode]' \
    '--sort[sort mode]:mode:(name ext size mtime)' \
    '--reverse[reverse sort order]' \
    '--include[include glob]:pattern:' \
    '--ignore[ignore glob]:pattern:' \
    '--du[show recursive sizes]' \
    '--json[output JSON]' \
    '--yaml[output YAML]' \
    '--theme[color theme]:mode:(auto color mono)' \
    '(-i --icons)'{-i,--icons}'[icon mode]:mode:(nerd unicode ascii none)' \
    '--completion[emit completion script]:shell:(bash zsh fish powershell)' \
    '*:path:_files'
}
_gotree_completion "$@"
`
}

func fishCompletion() string {
	return `complete -c gotree -s h -l help -d "show help"
complete -c gotree -s v -l version -d "show version"
complete -c gotree -s a -l all -d "show hidden files"
complete -c gotree -s L -d "max depth" -r
complete -c gotree -s d -l directories -d "list directories only"
complete -c gotree -s s -l sorted -d "sort by name"
complete -c gotree -s o -l order-by-extension -d "sort by extension"
complete -c gotree -s m -l summary -d "summary mode"
complete -c gotree -l sort -d "sort mode" -a "name ext size mtime"
complete -c gotree -l reverse -d "reverse sort order"
complete -c gotree -l include -d "include glob" -r
complete -c gotree -l ignore -d "ignore glob" -r
complete -c gotree -l du -d "show recursive sizes"
complete -c gotree -l json -d "output JSON"
complete -c gotree -l yaml -d "output YAML"
complete -c gotree -l theme -d "theme mode" -a "auto color mono"
complete -c gotree -s i -l icons -d "icon mode" -a "nerd unicode ascii none"
complete -c gotree -l completion -d "emit completion script" -a "bash zsh fish powershell"
complete -c gtr -w gotree
`
}

func powershellCompletion() string {
	return `Register-ArgumentCompleter -CommandName gotree,gtr -ScriptBlock {
  param($commandName, $wordToComplete, $cursorPosition)
  $flags = @(
    "-h","--help","-v","--version","-a","--all","-L","-d","--directories",
    "-s","--sorted","-o","--order-by-extension","-m","--summary","--sort",
    "--reverse","--include","--ignore","--du","--json","--yaml","--theme",
    "-i","--icons","--completion"
  )
  foreach ($f in $flags) {
    if ($f -like "$wordToComplete*") {
      [System.Management.Automation.CompletionResult]::new($f, $f, 'ParameterName', $f)
    }
  }
}
`
}
