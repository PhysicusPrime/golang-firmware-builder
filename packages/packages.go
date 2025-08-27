package packages

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/PhysicusPrime/golang-firmware-builder/command"
)

func DownloadPackage(pkg, dest string) string {
	fmt.Printf("Downloading %s...\n", pkg)
	command.ProgressBar(pkg+" Download", 3)

	url := fmt.Sprintf("https://ftp.gnu.org/gnu/%s/%s-latest.tar.gz", pkg, pkg)
	tarball := filepath.Join(dest, pkg+".tar.gz")
	command.RunCommandLive("wget", "-O", tarball, url)

	srcDir := filepath.Join(dest, pkg)
	command.RunCommandLive("mkdir", "-p", srcDir)
	command.RunCommandLive("tar", "xzf", tarball, "-C", srcDir, "--strip-components=1")

	fmt.Println(pkg, "heruntergeladen:", srcDir)
	return srcDir
}

func BuildPackage(srcDir, rootfs, cross string) {
	os.Chdir(srcDir)
	exec.Command("./configure", "--host="+cross, "--prefix="+rootfs).Run()
	exec.Command("make").Run()
	exec.Command("make", "install").Run()
	log.Println("Package gebaut:", srcDir)
}

func DownloadOpkg(dest string) string {
	fmt.Println("Downloading Opkg...")
	command.ProgressBar("Opkg Download", 3)

	url := "https://downloads.openwrt.org/sources/opkg-0.6.1.tar.gz"
	tarball := filepath.Join(dest, "opkg.tar.gz")
	command.RunCommandLive("wget", "-O", tarball, url)

	srcDir := filepath.Join(dest, "opkg")
	command.RunCommandLive("mkdir", "-p", srcDir)
	command.RunCommandLive("tar", "xzf", tarball, "-C", srcDir, "--strip-components=1")

	fmt.Println("Opkg heruntergeladen:", srcDir)
	return srcDir
}

func BuildOpkg(srcDir, rootfs, cross string) {
	os.Chdir(srcDir)
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
