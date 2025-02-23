#
# Windows script for project building
# Enable PowerShell scripts invoking in Settings for run it
# from any ConHost instance
#

$ProjectPath = Get-Location
$OutputPath = "$(Get-Location)\static\build\main.wasm"
Set-Location $ProjectPath
$Env:GOOS = "js"
$Env:GOARCH = "wasm"
Write-Host "building..."

go build -o $OutputPath "$($ProjectPath)\cmd\main.go"

if ($?) {
    Write-Host "success!"
    Write-Host "check it: $($OutputPath)"
} else {
    Write-Host "fail: $($?)"
    exit 1
}