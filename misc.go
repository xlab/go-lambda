package main

import (
	"archive/zip"
	"bytes"
	"log"
	"os"
	"os/exec"

	"github.com/jhoonb/archivex"
)

func packageName(pkg string) string {
	path := getExecPath("go")
	cmd := exec.Cmd{
		Path:   path,
		Args:   []string{path, "list", "-f", "{{.Name}}", pkg},
		Stderr: os.Stderr,
	}
	name, err := cmd.Output()
	if err != nil {
		log.Fatalln(err)
	}
	return string(bytes.TrimSpace(name))
}

func packageDir(pkg string) string {
	path := getExecPath("go")
	cmd := exec.Cmd{
		Path:   path,
		Args:   []string{path, "list", "-f", "{{.Dir}}", pkg},
		Stderr: os.Stderr,
	}
	name, err := cmd.Output()
	if err != nil {
		log.Fatalln(err)
	}
	return string(bytes.TrimSpace(name))
}

func getExecPath(name string) string {
	out, err := exec.Command("which", name).Output()
	if err != nil {
		log.Fatalf("executable file %s not found in $PATH", name)
		return ""
	}
	return string(bytes.TrimSpace(out))
}

func makeZip(main []byte, mainPath, libPath string, other ...string) []byte {
	buf := new(bytes.Buffer)
	zipper := &archivex.ZipFile{
		Name:   "source.zip",
		Writer: zip.NewWriter(buf),
	}
	if err := zipper.Add(mainPath, main); err != nil {
		log.Fatalln(err)
	}
	if err := zipper.AddFile(libPath); err != nil {
		log.Fatalf("error adding file '%s' to zip archive: %v", main, err)
	}
	for _, name := range other {
		info, err := os.Stat(name)
		if err != nil {
			log.Fatalln(err)
		}
		if info.IsDir() {
			err = zipper.AddAll(name, false)
		} else {
			err = zipper.AddFile(name)
		}
		if err != nil {
			log.Fatalln(err)
		}
	}
	zipper.Close()
	return buf.Bytes()
}
