package busybox

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	utils "command-line-argumentsC:\\Users\\hexzhen3x7\\GoLang-Buildsystem\\utils\\progress.go"
	"github.com/PhysicusPrime/golang-firmware-builder/command"
)

func DownloadBusyBox(dest string) string {
	fmt.Println("Downloading BusyBox...")
	utils.ProgressBar("BusyBox Download", 3)

	url := "https://busybox.net/downloads/busybox-1.36.0.tar.bz2"
	tarball := filepath.Join(dest, "busybox.tar.bz2")
	command.RunCommandLive("wget", "-O", tarball, url)

	srcDir := filepath.Join(dest, "busybox")
	command.RunCommandLive("mkdir", "-p", srcDir)
	command.RunCommandLive("tar", "xjf", tarball, "-C", srcDir, "--strip-components=1")

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
