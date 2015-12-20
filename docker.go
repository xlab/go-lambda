package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func mountPackageDir(pkg, dst string) string {
	return fmt.Sprintf("%s:%s", packageDir(pkg), dst)
}

func getDockerVersion(path string) (string, error) {
	out, err := exec.Command(path, "-v").Output()
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(out)), nil
}

func isNative() bool {
	return runtime.GOOS == "linux" && runtime.GOARCH == "amd64"
}

func runGopy(pkg string) {
	if isNative() {
		nativeGopy(pkg)
		return
	}
	dockerGopy(pkg)
}

func runBuild(pkg, modulePath string) {
	if isNative() {
		nativeBuild(modulePath)
		return
	}
	dockerBuild(pkg, modulePath)
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
	args = append(args, "gopy/gopy", "app", "bind", "-output", "/out", "in")
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

func nativeBuild(modulePath string) {
	if *debug {
		log.Println("go build native", modulePath)
	}

	path := getExecPath("go")
	pwd, _ := os.Getwd()
	args := []string{path, "build", "-buildmode", "c-shared", "-o", filepath.Join(pwd, "module.so")}
	if *debug {
		log.Println(strings.Join(args, " "))
	}
	cmd := exec.Cmd{
		Path:   path,
		Args:   args,
		Dir:    modulePath,
		Stderr: os.Stderr,
	}
	if *debug {
		cmd.Stdout = os.Stdout
	}
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
}

func dockerBuild(pkg, modulePath string) {
	if *debug {
		log.Println("go build via Docker", pkg)
	}
	pwd, _ := os.Getwd()
	dockerPath := getExecPath("docker")

	gopath, mounts := gopathMounts("/go")
	args := []string{dockerPath,
		"run", "-a", "stdout", "-a", "stderr", "--rm", "-e", "GOPATH=" + gopath,
		"-v", mountPackageDir(pkg, "/go/src/in"), "-v", fmt.Sprintf("%s:/out", pwd),
		"-v", fmt.Sprintf("%s:/go/src/module", modulePath),
	}
	for i, src := range mounts {
		args = append(args, "-v", fmt.Sprintf("%s:/go/path%d/src", src, i))
	}
	args = append(args, "xlab/go-lambda", "go", "build", "-buildmode", "c-shared", "-o", "/out/module.so", "module")
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
