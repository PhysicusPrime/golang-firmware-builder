package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"rpi4-firmware-builder/command"
	"rpi4-firmware-builder/utils"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	base := filepath.Join(home, "work")
	buildDir := filepath.Join(base, "build")
	rootfsDir := filepath.Join(buildDir, "rootfs")
	bootfsDir := filepath.Join(buildDir, "bootfs")
	downloadsDir := filepath.Join(buildDir, "downloads")

	fmt.Println("Erstelle Arbeitsverzeichnisse...")
	CreateDirs(buildDir, rootfsDir, bootfsDir, downloadsDir)

	fmt.Println("Erstelle RootFS-Struktur...")
	utils.ProgressBar("RootFS Struktur", 2)
	SetupRootFS(rootfsDir)

	cross := "aarch64-linux-gnu-"
	fmt.Println("Prüfe Toolchain...")
	utils.ProgressBar("Toolchain prüfen", 1)
	CheckToolchain(cross)

	// BusyBox
	fmt.Println("Download BusyBox...")
	utils.ProgressBar("BusyBox Download", 3)
	busyboxSrc := DownloadBusyBox(downloadsDir)

	fmt.Println("Prepare BusyBox...")
	utils.ProgressBar("BusyBox defconfig patchen", 2)
	PrepareBusyBox(busyboxSrc, cross)

	fmt.Println("Build BusyBox...")
	utils.ProgressBar("BusyBox bauen", 5)
	if err := command.RunCommandLive("make", "ARCH=arm64", "CROSS_COMPILE="+cross); err != nil {
		log.Fatal(err)
	}
	if err := command.RunCommandLive("make", "ARCH=arm64", "CROSS_COMPILE="+cross, "CONFIG_PREFIX="+rootfsDir, "install"); err != nil {
		log.Fatal(err)
	}

	// Andere Pakete
	pkgs := []string{"bash", "make", "cmake", "autoconf", "automake", "binutils", "libtool", "fdisk", "parted"}
	for _, pkg := range pkgs {
		fmt.Printf("Download & Build %s...\n", pkg)
		utils.ProgressBar(pkg+" Download", 2)
		src := DownloadPackage(pkg, downloadsDir)
		utils.ProgressBar(pkg+" Build", 4)
		BuildPackage(src, rootfsDir, cross)
	}

	// Opkg
	fmt.Println("Download & Build Opkg...")
	utils.ProgressBar("Opkg Download", 2)
	opkgSrc := DownloadOpkg(downloadsDir)
	utils.ProgressBar("Opkg Build", 4)
	BuildOpkg(opkgSrc, rootfsDir, cross)
	SetupOpkgConf(rootfsDir)

	fmt.Println("Firmware Build abgeschlossen! RootFS enthält jetzt BusyBox, Pakete und Opkg.")
}
