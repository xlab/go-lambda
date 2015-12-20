package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func isNativeGopy() bool {
	return runtime.GOOS == "linux" && runtime.GOARCH == "amd64"
}

func runGopy(pkg string) {
	if isNativeGopy() {
		nativeGopy(pkg)
		return
	}
	dockerGopy(pkg)
}

func nativeGopy(pkg string) {
	if *debug {
		log.Println("gopy bind native", pkg)
	}
	path := getExecPath("gopy")
	args := []string{path, "bind", pkg}
	if *debug {
		log.Println(strings.Join(args, " "))
	}
	cmd := exec.Cmd{
		Path:   path,
		Args:   args,
		Stderr: os.Stderr,
	}
	if *debug {
		cmd.Stdout = os.Stdout
	}
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}

func mountPackageDir(pkg, dst string) string {
	return fmt.Sprintf("%s:%s", packageDir(pkg), dst)
}

func dockerGopy(pkg string) {
	if *debug {
		log.Println("gopy bind via Docker", pkg)
	}
	pwd, _ := os.Getwd()
	dockerPath := getExecPath("docker")

	gopath, mounts := gopathMounts("/go")
	args := []string{dockerPath,
		"run", "-a", "stdout", "-a", "stderr", "--rm", "-e", "GOPATH=" + gopath,
		"-v", mountPackageDir(pkg, "/go/src/in"), "-v", fmt.Sprintf("%s:/out", pwd),
	}
	for i, src := range mounts {
		args = append(args, "-v", fmt.Sprintf("%s:/go/path%d/src", src, i))
	}
	args = append(args, "xlab/gopy", "app", "bind", "-output", "/out", "in")
	if *debug {
		log.Println(strings.Join(args, " "))
	}

	cmd := exec.Cmd{
		Path:   dockerPath,
		Args:   args,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	// if *debug {
	// 	cmd.Stdout = os.Stdout
	// }
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}

func getDockerVersion(path string) (string, error) {
	out, err := exec.Command(path, "-v").Output()
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(out)), nil
}
