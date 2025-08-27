package toolchain

import (
	"log"
	"os/exec"
)

// CheckToolchain prüft, ob die Cross-Compiler vorhanden sind
func CheckToolchain(prefix string) {
	_, err := exec.LookPath(prefix + "gcc")
	if err != nil {
		log.Fatalf("Toolchain nicht gefunden: %s", prefix+"gcc")
	}
	log.Println("Toolchain gefunden:", prefix)
}
