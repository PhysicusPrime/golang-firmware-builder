package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"rpi4-firmware-builder/command"
)

// PrepareBusyBox setzt defconfig und deaktiviert TC
func PrepareBusyBox(srcDir, cross string) {
	if err := os.Chdir(srcDir); err != nil {
		log.Fatalf("Fehler beim Wechseln ins BusyBox-Verzeichnis: %v", err)
	}

	fmt.Println("Setze defconfig für BusyBox...")
	if err := command.RunCommandLive("make", "defconfig"); err != nil {
		log.Fatalf("Fehler beim make defconfig: %v", err)
	}

	configPath := filepath.Join(srcDir, ".config")
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Fehler beim Lesen der .config: %v", err)
	}

	data := strings.ReplaceAll(string(b), "CONFIG_TC=y", "CONFIG_TC=n")
	if err := ioutil.WriteFile(configPath, []byte(data), 0644); err != nil {
		log.Fatalf("Fehler beim Schreiben der .config: %v", err)
	}

	fmt.Println("BusyBox defconfig gepatcht (TC deaktiviert).")
}

// BuildBusyBox kompiliert und installiert BusyBox ins RootFS
func BuildBusyBox(srcDir, rootfs, cross string) {
	if err := os.Chdir(srcDir); err != nil {
		log.Fatalf("Fehler beim Wechseln ins BusyBox-Verzeichnis: %v", err)
	}

	fmt.Println("BusyBox bauen...")
	if err := command.RunCommandLive("make", "ARCH=arm64", "CROSS_COMPILE="+cross); err != nil {
		log.Fatalf("Fehler beim BusyBox-Build: %v", err)
	}

	fmt.Println("BusyBox installieren...")
	if err := command.RunCommandLive("make", "ARCH=arm64", "CROSS_COMPILE="+cross, "CONFIG_PREFIX="+rootfs, "install"); err != nil {
		log.Fatalf("Fehler beim Installieren von BusyBox: %v", err)
	}

	fmt.Println("BusyBox erfolgreich installiert in:", rootfs)
}
