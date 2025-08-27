package busybox

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/PhysicusPrime/golang-firmware-builder/command"
	"github.com/PhysicusPrime/golang-firmware-builder/utils"
)

func DownloadBusyBox(dest string) string {
	fmt.Println("Downloading BusyBox...")
	utils.ProgressBar("BusyBox Download", 3)

	url := "https://busybox.net/downloads/busybox-1.36.0.tar.bz2"
	tarball := filepath.Join(dest, "busybox.tar.bz2")
	if err := command.RunCommandLive("wget", "-O", tarball, url); err != nil {
		log.Fatalf("Fehler beim Download von BusyBox: %v", err)
	}

	srcDir := filepath.Join(dest, "busybox")
	if err := command.RunCommandLive("mkdir", "-p", srcDir); err != nil {
		log.Fatalf("Fehler beim Erstellen von %s: %v", srcDir, err)
	}
	if err := command.RunCommandLive("tar", "xjf", tarball, "-C", srcDir, "--strip-components=1"); err != nil {
		log.Fatalf("Fehler beim Entpacken von BusyBox: %v", err)
	}

	fmt.Println("BusyBox heruntergeladen:", srcDir)
	return srcDir
}

func PrepareBusyBox(srcDir, cross string) {
	if err := os.Chdir(srcDir); err != nil {
		log.Fatalf("Fehler beim Wechseln ins BusyBox-Verzeichnis: %v", err)
	}

	fmt.Println("Setze defconfig für BusyBox...")
	if err := command.RunCommandLive("make", "defconfig"); err != nil {
		log.Fatalf("Fehler beim make defconfig: %v", err)
	}

	configPath := filepath.Join(srcDir, ".config")
	b, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Fehler beim Lesen der .config: %v", err)
	}

	data := strings.ReplaceAll(string(b), "CONFIG_TC=y", "CONFIG_TC=n")
	if err := os.WriteFile(configPath, []byte(data), 0644); err != nil {
		log.Fatalf("Fehler beim Schreiben der .config: %v", err)
	}

	fmt.Println("BusyBox defconfig gepatcht (TC deaktiviert).")
}
