package toolchain

import (
	"log"
	"os/exec"
)

func CheckToolchain(prefix string) {
	_, err := exec.LookPath(prefix + "gcc")
	if err != nil {
		log.Fatalf("Toolchain nicht gefunden: %s", prefix+"gcc")
	}
	log.Println("Toolchain gefunden:", prefix)
}
