package toolchain

import (
	"fmt"
	"log"
	"os/exec"
)

func CheckToolchain(cross string) {
	if err := exec.Command(cross+"gcc", "--version").Run(); err != nil {
		log.Fatalf("Toolchain %s nicht gefunden", cross)
	}
	fmt.Println("Toolchain gefunden:", cross)
}
