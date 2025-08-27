package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// CreateDirs erstellt alle benötigten Verzeichnisse
func CreateDirs(dirs ...string) {
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			log.Fatalf("Fehler beim Erstellen von %s: %v", d, err)
		}
	}
}

// SetupRootFS erstellt die grundlegende RootFS-Struktur
func SetupRootFS(rootfs string) {
	dirs := []string{
		"bin", "sbin", "etc", "proc", "sys", "dev", "tmp",
		"usr/bin", "usr/sbin", "var", "home",
	}

	for _, d := range dirs {
		os.MkdirAll(filepath.Join(rootfs, d), 0755)
	}

	createDevice(filepath.Join(rootfs, "dev", "null"), 1, 3)
	createDevice(filepath.Join(rootfs, "dev", "zero"), 1, 5)
	createDevice(filepath.Join(rootfs, "dev", "tty"), 5, 0)
	createDevice(filepath.Join(rootfs, "dev", "console"), 5, 1)

	os.Chmod(filepath.Join(rootfs, "tmp"), 01777)
}

func createDevice(path string, major, minor int) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(path), 0755)
		exec.Command("sudo", "mknod", "-m", "666", path, "c",
			fmt.Sprintf("%d", major), fmt.Sprintf("%d", minor)).Run()
	}
}
