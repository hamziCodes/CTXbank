# CTXbank 1-Click Installer for Windows
# Usage: irm https://raw.githubusercontent.com/hamziCodes/CTXbank/main/scripts/install.ps1 | iex

$ErrorActionPreference = "Stop"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  CTXbank Fast Installer (Windows)" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

# 1. Detect Architecture
$arch = "amd64"
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
    $arch = "arm64"
}
Write-Host "[1/4] Detected architecture: $arch" -ForegroundColor Gray

# 2. Query GitHub for Latest Release
$repo = "hamziCodes/CTXbank"
$releasesUrl = "https://api.github.com/repos/$repo/releases/latest"
Write-Host "[2/4] Querying latest release from GitHub ($repo)..." -ForegroundColor Gray

try {
    $release = Invoke-RestMethod -Uri $releasesUrl -Headers @{"User-Agent"="CTXbank-Installer"}
    $tag = $release.tag_name
} catch {
    # Fallback to default v0.1.0 if API rate limit occurs
    $tag = "v0.1.0"
}

$assetName = "ctx-windows-$arch.zip"
$downloadUrl = "https://github.com/$repo/releases/download/$tag/$assetName"

# 3. Create Install Directory and Download
$installDir = Join-Path $HOME ".ctxbank\bin"
$tempZip = Join-Path $env:TEMP "ctxbank-$arch.zip"

if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}

Write-Host "[3/4] Downloading $assetName ($tag)..." -ForegroundColor Gray
Invoke-WebRequest -Uri $downloadUrl -OutFile $tempZip -UseBasicParsing

Write-Host "      Extracting to $installDir..." -ForegroundColor Gray
Expand-Archive -Path $tempZip -DestinationPath $installDir -Force
Remove-Item $tempZip -Force -ErrorAction SilentlyContinue

# 4. Add to User PATH
$userPath = [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::User)
if ($userPath -notlike "*$installDir*") {
    Write-Host "[4/4] Adding $installDir to User PATH..." -ForegroundColor Gray
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", [EnvironmentVariableTarget]::User)
    $env:Path = "$env:Path;$installDir"
} else {
    Write-Host "[4/4] $installDir is already in User PATH." -ForegroundColor Gray
}

Write-Host "`n+-------------------------------------------------------------+" -ForegroundColor Green
Write-Host "|  CTXbank installed successfully!                            |" -ForegroundColor Green
Write-Host "+-------------------------------------------------------------+" -ForegroundColor Green
Write-Host "`nRun 'ctx --help' in a new terminal window to get started." -ForegroundColor White
Write-Host "Location: $installDir\ctx.exe`n" -ForegroundColor Gray
