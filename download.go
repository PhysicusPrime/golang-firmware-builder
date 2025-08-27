package main

import (
	"fmt"
	"path/filepath"

	"rpi4-firmware-builder/command"
	"rpi4-firmware-builder/utils"
)

// DownloadBusyBox lädt BusyBox herunter und entpackt
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
