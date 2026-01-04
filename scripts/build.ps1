<#
.SYNOPSIS
    Build script for Lunno (Go project).
.DESCRIPTION
    Automates building Lunno for multiple OS/ARCH targets, embedding version info,
    generating checksums, and optionally running tests. Supports sequential or parallel builds.
.PARAMETER Configuration
    Build configuration: Debug or Release (default: Release). Determines optimization flags.
.PARAMETER Parallel
    Run builds for multiple targets in parallel.
.PARAMETER Targets
    Optional list of OS/ARCH targets (format: os/arch). Defaults to common Linux, macOS, Windows targets.
#>

[CmdletBinding()]
param (
    [ValidateSet("Debug", "Release")]
    [string]$Configuration = "Release",
    [switch]$RunTests,
    [switch]$Parallel,
    [string[]]$Targets
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Fail {
    param(
        [string]$Message,
        [int]$Code = 1
    )
    Write-Error $Message
    exit $Code
}

function Success {
    param(
        [string]$Message
    )
    Write-Host $Message
}

function VerboseLog {
    param(
        [string]$Message
    )
    Write-Verbose $Message
}

$StartTime = Get-Date

$ProjectRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$ProjectRoot = $ProjectRoot.ProviderPath
$SrcDir = Join-Path $ProjectRoot "cmd/lunno"
$DistDir = Join-Path $ProjectRoot "dist"

if (-not (Test-Path $DistDir)) { New-Item -ItemType Directory -Path $DistDir | Out-Null }

if (-not $Targets -or $Targets.Count -eq 0) {
    $Targets = @(
        "linux/amd64", "linux/arm64", "linux/arm",
        "darwin/amd64", "darwin/arm64",
        "windows/amd64", "windows/arm64"
    )
}

$Version = if ($env:VERSION) {
    $env:VERSION
}
else {
    "0.1.0"
}
$BuildType = if ($env:BUILD_TYPE) {
    $env:BUILD_TYPE
}
else {
    "alpha"
}
$FullVersion = "$Version-$BuildType"

try {
    $GitCommit = git rev-parse --short HEAD
}
catch {
    $GitCommit = "unknown"
}

$LdFlags = "-X 'lunno/internal/version.Version=$FullVersion'
            -X 'lunno/internal/version.GitCommit=$GitCommit'"

function Build-Target {
    param(
        [string]$Target
    )
    $Split = $Target.Split('/')
    $GOOS = $Split[0]
    $GOARCH = $Split[1]

    $BinaryName = "Lunno"
    $OutPath = Join-Path $DistDir "$BinaryName-$GOOS-$GOARCH"
    if ($GOOS -eq "windows") {
        $OutPath += ".exe"
    }

    Write-Host "[INFO] Building $Target -> $OutPath"
    $Env:GOOS = $GOOS
    $Env:GOARCH = $GOARCH

    $BuildArgs = @("build", "-ldflags", $LdFlags, "-o", $OutPath, $SrcDir)
    try {
        go @BuildArgs
    }
    catch {
        Fail "Build failed for $Target"
    }

    $ChecksumFile = "$OutPath-checksum.sha256"
    if (Get-Command shasum -ErrorAction SilentlyContinue) {
        shasum -a 256 $OutPath | Out-File -Encoding ascii $ChecksumFile
    }
    else {
        Get-FileHash -Algorithm SHA256 $OutPath | ForEach-Object {
            "$($_.Hash)  $OutPath"
        } | Out-File -Encoding ascii $ChecksumFile
    }

    return [PSCustomObject]@{
        Target   = $Target;
        Binary   = $OutPath;
        Checksum = (Get-Content $ChecksumFile)
    }
}

$Results = @()
if ($Parallel) {
    $Jobs = @()
    foreach ($t in $Targets) {
        $Jobs += Start-Job -ScriptBlock {
            param($t) Build-Target -Target $t
        } -ArgumentList $t
    }
    $Jobs | Wait-Job
    $Results = $Jobs | ForEach-Object {
        Receive-Job -Job $_
    }
    $Jobs | ForEach-Object {
        Remove-Job -Job $_ -Force
    }
}
else {
    foreach ($t in $Targets) {
        $Results += Build-Target -Target $t
    }
}

$FullVersion | Out-File -Encoding ascii (Join-Path $DistDir "FULL_VERSION.txt")

$EndTime = Get-Date
$Duration = [math]::Round(($EndTime - $StartTime).TotalSeconds, 2)
Write-Host ""
Write-Host "Build Summary for Lunno $FullVersion"
$Results | ForEach-Object {
    [PSCustomObject]@{
        Target   = $_.Target
        Binary   = Split-Path $_.Binary -Leaf
        Checksum = $_.Checksum.Split(" ")[0]
    }
} | Format-Table Target, Binary, Checksum -AutoSize
Write-Host ""
Success "Build completed in $Duration seconds. Binaries and checksums are in $DistDir"