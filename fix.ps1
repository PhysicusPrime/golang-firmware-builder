# PowerShell Skript zum Setup des Go Projekts
$ErrorActionPreference = "Stop"

$RepoURL = "https://github.com/PhysicusPrime/golang-firmware-builder.git"
$ModuleName = "github.com/PhysicusPrime/golang-firmware-builder"
$Branch = "main"

# 1️⃣ Git Repo initialisieren
if (-not (Test-Path ".git")) {
    Write-Host "[*] Initialisiere Git Repository..."
    git init
    git remote add origin $RepoURL
}

# 2️⃣ go.mod erstellen / initialisieren
if (-not (Test-Path "go.mod")) {
    Write-Host "[*] Initialisiere go.mod..."
    go mod init $ModuleName
} else {
    Write-Host "[*] go.mod existiert bereits."
}

# 3️⃣ Alle Go-Files prüfen und Imports korrigieren

# Alle Go-Files prüfen und Imports korrigieren
Write-Host "[*] Korrigiere Import-Pfade in allen .go Dateien..."
Get-ChildItem -Recurse -Filter *.go | ForEach-Object {
    $file = $_.FullName
    $content = Get-Content $file

    # Jede Ersetzung einzeln
    $content = $content -replace '"utils"', ('"' + $ModuleName + '/utils"')
    $content = $content -replace '"command"', ('"' + $ModuleName + '/command"')
    
    # Datei zurückschreiben
    Set-Content -Path $file -Value $content
}


# 4️⃣ Git Status anzeigen
Write-Host "[*] Git Status:"
git status

# 5️⃣ Dateien zum Commit vorbereiten
git add .

# 6️⃣ Commit erstellen
git commit -m "Initial Go project setup with module and import paths"

# 7️⃣ Auf main pushen
git branch -M $Branch
git push -u origin $Branch

Write-Host "[*] Fertig! Repository auf $Branch gepusht."
