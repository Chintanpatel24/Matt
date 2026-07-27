# ==============================================================================
#  Matt File Manager Installer (PowerShell)
#  Repo: https://github.com/Chintanpatel24/Matt
# ==============================================================================

$ErrorActionPreference = "Stop"

Write-Host "Installing Matt Black Terminal File Manager..." -ForegroundColor Cyan

$InstallDir = "$env:LOCALAPPDATA\Matt"
if (!(Test-Path -Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}

$BinaryPath = "$InstallDir\matt.exe"

if (Get-Command go -ErrorAction SilentlyContinue) {
    Write-Host "-> Go compiler detected. Building Matt from source..." -ForegroundColor Yellow
    $TmpDir = [System.IO.Path]::GetTempFileName() + "_dir"
    git clone --depth 1 https://github.com/Chintanpatel24/Matt.git $TmpDir
    Set-Location $TmpDir
    go build -o $BinaryPath ./cmd/matt
    Set-Location $env:USERPROFILE
    Remove-Item -Recurse -Force $TmpDir
} else {
    Write-Host "-> Downloading prebuilt Windows binary..." -ForegroundColor Yellow
    $DownloadUrl = "https://github.com/Chintanpatel24/Matt/releases/latest/download/matt-windows-amd64.exe"
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $BinaryPath
}

# Add to User PATH if not present
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    Write-Host "-> Adding $InstallDir to User PATH..." -ForegroundColor Yellow
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
}

Write-Host "`n✓ Matt successfully installed to $BinaryPath" -ForegroundColor Green
Write-Host "Restart your terminal and run 'matt' to start!" -ForegroundColor Cyan
