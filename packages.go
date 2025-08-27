package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"rpi4-firmware-builder/command"
)

// BuildPackage kompiliert ein GNU-Paket mit configure/make
func BuildPackage(srcDir, rootfs, cross string) {
	if err := os.Chdir(srcDir); err != nil {
		log.Fatalf("Fehler beim Wechseln in %s: %v", srcDir, err)
	}
	if err := command.RunCommandLive("./configure", "--host="+cross, "--prefix="+rootfs); err != nil {
		log.Fatalf("Fehler beim Configure: %v", err)
	}
	if err := command.RunCommandLive("make"); err != nil {
		log.Fatalf("Fehler beim Make: %v", err)
	}
	if err := command.RunCommandLive("make", "install"); err != nil {
		log.Fatalf("Fehler beim Make Install: %v", err)
	}
	log.Println("Package gebaut:", srcDir)
}

// BuildOpkg kompiliert Opkg und kopiert Binary ins RootFS
func BuildOpkg(srcDir, rootfs, cross string) {
	if err := os.Chdir(srcDir); err != nil {
		log.Fatalf("Fehler beim Wechseln in %s: %v", srcDir, err)
	}

	cmd := exec.Command("make", fmt.Sprintf("CC=%sgcc", cross), fmt.Sprintf("PREFIX=%s", rootfs))
	cmd.Env = append(os.Environ(), "CROSS_COMPILE="+cross)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	binPath := filepath.Join(rootfs, "usr", "bin")
	os.MkdirAll(binPath, 0755)
	exec.Command("cp", "opkg", binPath).Run()
	fmt.Println("Opkg installiert")
}

// SetupOpkgConf erstellt die minimalen Konfig-Dateien für Opkg
func SetupOpkgConf(rootfs string) {
	etcDir := filepath.Join(rootfs, "etc", "opkg")
	os.MkdirAll(etcDir, 0755)
	conf := `dest root /
lists_dir ext /var/lib/opkg
option overlay_root /overlay
option check_signature 0
`
	os.WriteFile(filepath.Join(etcDir, "opkg.conf"), []byte(conf), 0644)
	os.MkdirAll(filepath.Join(rootfs, "var/lib/opkg"), 0755)
	os.MkdirAll(filepath.Join(rootfs, "overlay"), 0755)
}
