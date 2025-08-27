#!/bin/bash
set -e

# Projektname und Modul
MODULE="github.com/PhysicusPrime/golang-firmware-builder"

echo "[*] Erstelle go.mod..."
go mod init $MODULE || echo "go.mod existiert bereits"
go mod tidy

echo "[*] Korrigiere Paketnamen in Unterordnern..."
# Unterordner und gewünschte Paketnamen
declare -A FOLDERS=(
    ["busybox"]="busybox"
    ["command"]="command"
    ["fs"]="fs"
    ["packages"]="packages"
    ["toolchain"]="toolchain"
    ["utils"]="utils"
)

for folder in "${!FOLDERS[@]}"; do
    pkg=${FOLDERS[$folder]}
    if [ -d "$folder" ]; then
        find "$folder" -name "*.go" | while read file; do
            echo "Setze package $pkg in $file"
            sed -i "1s/^package .*/package $pkg/" "$file"
        done
    fi
done

echo "[*] Passe Imports in main.go an..."
MAIN_FILE="main.go"
for folder in "${!FOLDERS[@]}"; do
    pkg=${FOLDERS[$folder]}
    # Entferne alte Imports falls vorhanden
    sed -i "/\"$pkg\"/d" $MAIN_FILE
done

# Füge korrekte Imports hinzu
IMPORTS=$(for folder in "${!FOLDERS[@]}"; do
    pkg=${FOLDERS[$folder]}
    echo "    \"$MODULE/$pkg\""
done)

# Ersetze Platzhalter oder füge nach import ( falls nötig )
sed -i "/import (/a $IMPORTS" $MAIN_FILE

echo "[*] Fertig! Du kannst jetzt 'go mod tidy' und 'go build' ausführen."

