# fix_imports.ps1
# Setzt alle Paketnamen und korrigiert Imports für Go-Modul
# Nutze: ./fix_imports.ps1

$Module = "github.com/PhysicusPrime/golang-firmware-builder"

# Unterordner und gewünschte Paketnamen
$folders = @{
    "busybox" = "busybox"
    "command" = "command"
    "fs" = "fs"
    "packages" = "packages"
    "toolchain" = "toolchain"
    "utils" = "utils"
}

# 1. Paketnamen setzen
Write-Host "[*] Setze Paketnamen in Unterordnern..."
foreach ($folder in $folders.Keys) {
    $pkg = $folders[$folder]
    if (Test-Path $folder) {
        Get-ChildItem -Path $folder -Filter *.go | ForEach-Object {
            (Get-Content $_.FullName) |
                ForEach-Object { $_ -replace '^package .*', "package $pkg" } |
                Set-Content $_.FullName
            Write-Host "Setze package $pkg in $($_.FullName)"
        }
    }
}

# 2. main.go Imports korrigieren
$mainFile = "main.go"
if (Test-Path $mainFile) {
    Write-Host "[*] Korrigiere Imports in main.go..."
    
    # Alte Imports der Unterordner entfernen
    foreach ($folder in $folders.Keys) {
        $pkg = $folders[$folder]
        (Get-Content $mainFile) |
            ForEach-Object { $_ -replace ".*$pkg.*", "" } |
            Set-Content $mainFile
    }

    # Korrekte Imports hinzufügen
    $imports = $folders.Values | ForEach-Object { "    `"$Module/$_`"" }
    $importBlock = @("import (") + $imports + @(")")
    
    # Entferne existierenden import Block und setze neuen
    $content = Get-Content $mainFile -Raw
    $content = $content -replace "import\s*\([\s\S]*?\)", ($importBlock -join "`n")
    Set-Content $mainFile $content
}

# 3. go.mod erstellen, falls nicht vorhanden
if (-Not (Test-Path "go.mod")) {
    Write-Host "[*] Erstelle go.mod..."
    go mod init $Module
}
Write-Host "[*] Fertig! Du kannst jetzt 'go mod tidy' und 'go build' ausführen."
