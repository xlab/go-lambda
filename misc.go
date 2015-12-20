package main

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jhoonb/archivex"
)

func gopathMounts(dir string) (string, []string) {
	gopath := "/go"
	srcDirs := build.Default.SrcDirs()
	for i := range srcDirs {
		gopath += fmt.Sprintf(":%s/path%d", dir, i)
	}
	return gopath, srcDirs
}

func packageName(pkgImport string) (name string) {
	for _, src := range build.Default.SrcDirs() {
		pkg, err := build.Import(pkgImport, src, build.ImportComment)
		if err != nil || pkg.IsCommand() || pkg.Goroot {
			continue
		}
		name = pkg.Name
		return
	}
	log.Fatalln("package not found in $GOPATH:", pkgImport)
	return
}

func packageDir(pkgImport string) (dir string) {
	for _, src := range build.Default.SrcDirs() {
		pkg, err := build.Import(pkgImport, src, build.ImportComment)
		if err != nil || pkg.IsCommand() || pkg.Goroot {
			continue
		}
		dir = pkg.Dir
		return
	}
	log.Fatalln("package not found in $GOPATH:", pkgImport)
	return
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

func getModuleSource(packageName, packageFunc string) []byte {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "import module\nimport json\n\n")
	fmt.Fprintf(buf, "handler = module.lambda_handler()\n\n")
	fmt.Fprintf(buf, "def %s(event, context):\n", packageFunc)
	fmt.Fprintln(buf, "    return handler(json.dumps(event), json.dumps(context))")
	return buf.Bytes()
}

func getTempDir() string {
	tag := make([]byte, 8)
	_, err := rand.Read(tag)
	if err != nil {
		return filepath.Join(os.TempDir(), "go-lambda")
	}
	return filepath.Join(os.TempDir(), fmt.Sprintf("go-lambda-%x", tag))
}

func buildModuleBridge(tempDir, pkg, packageFunc string) {
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		log.Fatalln(err)
	}
	defer os.RemoveAll(tempDir)

	err := ioutil.WriteFile(filepath.Join(tempDir, "module.c"), MustAsset("module/module.c"), 0644)
	if err != nil {
		log.Fatalln(err)
	}
	err = ioutil.WriteFile(filepath.Join(tempDir, "module.h"), MustAsset("module/module.h"), 0644)
	if err != nil {
		log.Fatalln(err)
	}
	goAsset := MustAsset("module/module.go")
	funcRef := fmt.Sprintf(`&lambda.%s`, strings.Title(packageFunc))
	goAsset = bytes.Replace(goAsset, []byte(`&lambda.Handler`), []byte(funcRef), -1)

	if isNative() {
		pkgRef := fmt.Sprintf(`lambda "%s"`, pkg)
		goAsset = bytes.Replace(goAsset, []byte(`lambda "in"`), []byte(pkgRef), -1)
	}

	err = ioutil.WriteFile(filepath.Join(tempDir, "module.go"), goAsset, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	runBuild(pkg, tempDir)
}
