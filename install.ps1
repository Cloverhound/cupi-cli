#Requires -Version 5.1
[CmdletBinding()]
param(
    [string]$InstallDir = "$env:LOCALAPPDATA\Programs\cupi-cli"
)

$Repo = "Cloverhound/cupi-cli"
$Binary = "cupi.exe"
$ErrorActionPreference = "Stop"

# Detect architecture
$Arch = if ([System.Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
} else {
    Write-Error "32-bit Windows is not supported."
    exit 1
}

# Fetch latest release
Write-Host "Fetching latest release..."
$Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
$Tag = $Release.tag_name
$Version = $Tag.TrimStart("v")

Write-Host "Installing cupi $Tag (windows/$Arch)..."

$Archive = "cupi_${Version}_windows_${Arch}.zip"
$Url = "https://github.com/$Repo/releases/download/$Tag/$Archive"
$TmpDir = [System.IO.Path]::GetTempPath() + [System.Guid]::NewGuid().ToString()
New-Item -ItemType Directory -Path $TmpDir | Out-Null

try {
    $ArchivePath = Join-Path $TmpDir $Archive
    Invoke-WebRequest -Uri $Url -OutFile $ArchivePath
    Expand-Archive -Path $ArchivePath -DestinationPath $TmpDir -Force

    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    Copy-Item -Path (Join-Path $TmpDir $Binary) -Destination $InstallDir -Force

    Write-Host "Installed to $InstallDir\$Binary"
} finally {
    Remove-Item -Recurse -Force $TmpDir -ErrorAction SilentlyContinue
}

# Add to PATH if not present
$UserPath = [System.Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [System.Environment]::SetEnvironmentVariable("PATH", "$UserPath;$InstallDir", "User")
    Write-Host "Added $InstallDir to your user PATH."
    Write-Host "Restart your terminal for the PATH change to take effect."
}

Write-Host ""
Write-Host "Run 'cupi --help' to get started."
