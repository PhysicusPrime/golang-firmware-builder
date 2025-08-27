package packages

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/PhysicusPrime/golang-firmware-builder/command"
	"github.com/PhysicusPrime/golang-firmware-builder/utils"
)

// Map mit festen Versionen der Pakete
var pkgVersions = map[string]string{
	"bash":     "5.2",
	"make":     "4.4",
	"cmake":    "3.27.6",
	"autoconf": "2.71",
	"automake": "1.17.3",
	"binutils": "2.41",
	"libtool":  "2.4.7",
	"fdisk":    "2.40.0",
	"parted":   "3.5",
}

// DownloadOpkg lädt Opkg herunter und entpackt es
func DownloadOpkg(dest string) string {
	fmt.Println("Downloading Opkg...")
	utils.ProgressBar("Opkg Download", 3)

	url := "https://downloads.openwrt.org/sources/opkg-0.6.1.tar.gz"
	tarball := filepath.Join(dest, "opkg.tar.gz")
	if err := command.RunCommandLive("wget", "-O", tarball, url); err != nil {
		log.Fatalf("Fehler beim Download von Opkg: %v", err)
	}

	srcDir := filepath.Join(dest, "opkg")
	if err := command.RunCommandLive("mkdir", "-p", srcDir); err != nil {
		log.Fatalf("Fehler beim Erstellen von %s: %v", srcDir, err)
	}
	if err := command.RunCommandLive("tar", "xzf", tarball, "-C", srcDir, "--strip-components=1"); err != nil {
		log.Fatalf("Fehler beim Entpacken von Opkg: %v", err)
	}

	fmt.Println("Opkg heruntergeladen:", srcDir)
	return srcDir
}

// BuildOpkg baut Opkg und installiert die Binärdatei ins RootFS
func BuildOpkg(srcDir, rootfs, cross string) {
	if err := os.Chdir(srcDir); err != nil {
		log.Fatalf("Fehler beim Wechseln ins Opkg-Verzeichnis: %v", err)
	}

	cmd := exec.Command("make", fmt.Sprintf("CC=%sgcc", cross), fmt.Sprintf("PREFIX=%s", rootfs))
	cmd.Env = append(os.Environ(), "CROSS_COMPILE="+cross)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Fehler beim Bauen von Opkg: %v", err)
	}

	binPath := filepath.Join(rootfs, "usr", "bin")
	os.MkdirAll(binPath, 0755)
	if err := exec.Command("cp", "opkg", binPath).Run(); err != nil {
		log.Fatalf("Fehler beim Kopieren von Opkg: %v", err)
	}

	fmt.Println("Opkg installiert")
}

// SetupOpkgConf erstellt die Standardkonfiguration für Opkg im RootFS
func SetupOpkgConf(rootfs string) {
	etcDir := filepath.Join(rootfs, "etc", "opkg")
	os.MkdirAll(etcDir, 0755)
	conf := `dest root /
lists_dir ext /var/lib/opkg
option overlay_root /overlay
option check_signature 0
`
	if err := os.WriteFile(filepath.Join(etcDir, "opkg.conf"), []byte(conf), 0644); err != nil {
		log.Fatalf("Fehler beim Schreiben von opkg.conf: %v", err)
	}

	os.MkdirAll(filepath.Join(rootfs, "var/lib/opkg"), 0755)
	os.MkdirAll(filepath.Join(rootfs, "overlay"), 0755)
}

// DownloadPackage lädt ein GNU-Paket herunter und entpackt es
func DownloadPackage(pkg, dest string) string {
	fmt.Printf("Downloading %s...\n", pkg)
	utils.ProgressBar(pkg+" Download", 3)

	version, ok := pkgVersions[pkg]
	if !ok {
		log.Fatalf("Keine Version für Paket %s definiert", pkg)
	}

	url := fmt.Sprintf("https://ftp.gnu.org/gnu/%s/%s-%s.tar.gz", pkg, pkg, version)
	tarball := filepath.Join(dest, pkg+".tar.gz")
	if err := command.RunCommandLive("wget", "-O", tarball, url); err != nil {
		log.Fatalf("Fehler beim Download von %s: %v", pkg, err)
	}

	srcDir := filepath.Join(dest, pkg)
	if err := command.RunCommandLive("mkdir", "-p", srcDir); err != nil {
		log.Fatalf("Fehler beim Erstellen von %s: %v", srcDir, err)
	}
	if err := command.RunCommandLive("tar", "xzf", tarball, "-C", srcDir, "--strip-components=1"); err != nil {
		log.Fatalf("Fehler beim Entpacken von %s: %v", pkg, err)
	}

	fmt.Println(pkg, "heruntergeladen:", srcDir)
	return srcDir
}

// BuildPackage konfiguriert, baut und installiert ein Paket ins RootFS
func BuildPackage(srcDir, rootfs, cross string) {
	if err := os.Chdir(srcDir); err != nil {
		log.Fatalf("Fehler beim Wechseln ins Paketverzeichnis: %v", err)
	}

	if err := exec.Command("./configure", "--host="+cross, "--prefix="+rootfs).Run(); err != nil {
		log.Fatalf("Fehler beim Konfigurieren von %s: %v", srcDir, err)
	}
	if err := exec.Command("make").Run(); err != nil {
		log.Fatalf("Fehler beim Bauen von %s: %v", srcDir, err)
	}
	if err := exec.Command("make", "install").Run(); err != nil {
		log.Fatalf("Fehler beim Installieren von %s: %v", srcDir, err)
	}

	log.Println("Package gebaut:", srcDir)
}
