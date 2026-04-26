Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$RootDir = Resolve-Path (Join-Path $PSScriptRoot "..")
$DistDir = Join-Path $RootDir "dist"
$BinDir = Join-Path $DistDir "bin"
$CompDir = Join-Path $DistDir "completions"

New-Item -ItemType Directory -Path $BinDir -Force | Out-Null
New-Item -ItemType Directory -Path $CompDir -Force | Out-Null

function Build-Binary {
    param(
        [Parameter(Mandatory = $true)][string]$GoOS,
        [Parameter(Mandatory = $true)][string]$GoArch
    )

    $ext = ""
    if ($GoOS -eq "windows") {
        $ext = ".exe"
    }

    $outPath = Join-Path $BinDir ("gotree_{0}_{1}{2}" -f $GoOS, $GoArch, $ext)
    Write-Host "Building $outPath"
    $env:GOOS = $GoOS
    $env:GOARCH = $GoArch
    go build -o $outPath ./src/cmd/gotree
}

Build-Binary -GoOS "linux" -GoArch "amd64"
Build-Binary -GoOS "darwin" -GoArch "amd64"
Build-Binary -GoOS "darwin" -GoArch "arm64"
Build-Binary -GoOS "windows" -GoArch "amd64"

Write-Host "Generating completions..."
go run ./src/cmd/gotree --completion bash | Out-File (Join-Path $CompDir "gotree.bash") -Encoding utf8
go run ./src/cmd/gotree --completion zsh | Out-File (Join-Path $CompDir "_gotree") -Encoding utf8
go run ./src/cmd/gotree --completion fish | Out-File (Join-Path $CompDir "gotree.fish") -Encoding utf8
go run ./src/cmd/gotree --completion powershell | Out-File (Join-Path $CompDir "gotree.ps1") -Encoding utf8

Write-Host "Done. Artifacts in $DistDir"
