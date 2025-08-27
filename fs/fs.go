package fs

import (
	"log"
	"os"
)

func CreateDirs(dirs ...string) {
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Fehler beim Erstellen von %s: %v", dir, err)
		}
	}
}

func SetupRootFS(rootfs string) {
	CreateDirs(rootfs+"/bin", rootfs+"/sbin", rootfs+"/etc", rootfs+"/usr/bin", rootfs+"/usr/sbin")
}
