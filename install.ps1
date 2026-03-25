$ErrorActionPreference = "Stop"

$repo = "AntonZubritski/yamp"
$file = "yamp-windows-amd64.exe"
$url = "https://github.com/$repo/releases/latest/download/$file"
$installDir = "$env:LOCALAPPDATA\Microsoft\WinGet\Links"

Write-Host ""
Write-Host "  yamp installer for Windows" -ForegroundColor Cyan
Write-Host "  Downloading $file..." -ForegroundColor Gray
Write-Host ""

if (!(Test-Path $installDir)) { New-Item -ItemType Directory -Path $installDir -Force | Out-Null }

$tmp = Join-Path $env:TEMP "yamp.exe"
Invoke-WebRequest -Uri $url -OutFile $tmp -UseBasicParsing

Move-Item -Force $tmp "$installDir\yamp.exe"

Write-Host "  Installed to $installDir\yamp.exe" -ForegroundColor Green
Write-Host ""
Write-Host "  Open a NEW terminal and run: yamp" -ForegroundColor Yellow
Write-Host ""
