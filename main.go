package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/PhysicusPrime/golang-firmware-builder/busybox"
	"github.com/PhysicusPrime/golang-firmware-builder/command"
	"github.com/PhysicusPrime/golang-firmware-builder/fs"
	"github.com/PhysicusPrime/golang-firmware-builder/packages"
	"github.com/PhysicusPrime/golang-firmware-builder/toolchain"
	"github.com/PhysicusPrime/golang-firmware-builder/utils"
)

// DownloadBusyBox lädt BusyBox herunter und entpackt sie
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

// DownloadPackage lädt ein GNU-Paket herunter
func DownloadPackage(pkg, dest string) string {
	fmt.Printf("Downloading %s...\n", pkg)
	utils.ProgressBar(pkg+" Download", 3)

	url := fmt.Sprintf("https://ftp.gnu.org/gnu/%s/%s-latest.tar.gz", pkg, pkg)
	tarball := filepath.Join(dest, pkg+".tar.gz")
	command.RunCommandLive("wget", "-O", tarball, url)

	srcDir := filepath.Join(dest, pkg)
	command.RunCommandLive("mkdir", "-p", srcDir)
	command.RunCommandLive("tar", "xzf", tarball, "-C", srcDir, "--strip-components=1")

	fmt.Println(pkg, "heruntergeladen:", srcDir)
	return srcDir
}

// DownloadOpkg lädt Opkg herunter
func DownloadOpkg(dest string) string {
	fmt.Println("Downloading Opkg...")
	utils.ProgressBar("Opkg Download", 3)

	url := "https://downloads.openwrt.org/sources/opkg-0.6.1.tar.gz"
	tarball := filepath.Join(dest, "opkg.tar.gz")
	command.RunCommandLive("wget", "-O", tarball, url)

	srcDir := filepath.Join(dest, "opkg")
	command.RunCommandLive("mkdir", "-p", srcDir)
	command.RunCommandLive("tar", "xzf", tarball, "-C", srcDir, "--strip-components=1")

	fmt.Println("Opkg heruntergeladen:", srcDir)
	return srcDir
}

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
	fs.CreateDirs(buildDir, rootfsDir, bootfsDir, downloadsDir)

	fmt.Println("Erstelle RootFS-Struktur...")
	utils.ProgressBar("RootFS Struktur", 2)
	fs.SetupRootFS(rootfsDir)

	cross := "aarch64-linux-gnu-"
	fmt.Println("Prüfe Toolchain...")
	utils.ProgressBar("Toolchain prüfen", 1)
	toolchain.CheckToolchain(cross)

	// BusyBox
	fmt.Println("Download BusyBox...")
	utils.ProgressBar("BusyBox Download", 3)
	busyboxSrc := busybox.DownloadBusyBox(downloadsDir)

	fmt.Println("Prepare BusyBox...")
	utils.ProgressBar("BusyBox defconfig patchen", 2)
	busybox.PrepareBusyBox(busyboxSrc, cross)

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
		src := packages.DownloadPackage(pkg, downloadsDir)
		utils.ProgressBar(pkg+" Build", 4)
		packages.BuildPackage(src, rootfsDir, cross)
	}

	// Opkg
	fmt.Println("Download & Build Opkg...")
	utils.ProgressBar("Opkg Download", 2)
	opkgSrc := packages.DownloadOpkg(downloadsDir)
	utils.ProgressBar("Opkg Build", 4)
	packages.BuildOpkg(opkgSrc, rootfsDir, cross)
	packages.SetupOpkgConf(rootfsDir)

	fmt.Println("Firmware Build abgeschlossen! RootFS enthält jetzt BusyBox, Pakete und Opkg.")
}
