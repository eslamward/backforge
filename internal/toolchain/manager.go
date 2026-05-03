package toolchain

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func getGoURL() string {
	base := "https://go.dev/dl/go1.25.9"
	os := runtime.GOOS
	arch := runtime.GOARCH

	switch os {
	case "windows":
		return fmt.Sprintf("%s.%s-%s.zip", base, os, arch)
	case "linux", "darwin":
		return fmt.Sprintf("%s.%s-%s.tar.gz", base, os, arch)
	default:
		return ""
	}
}
func EnsureGo() (string, error) {

	if Exists() {
		fmt.Println("Go already cached")
		return GoPath(), nil
	}

	// create cache folder
	os.MkdirAll(".cache/go", os.ModePerm)

	zipPath := ".cache/go.zip"

	// download
	err := DownloadGo(getGoURL(), zipPath)
	if err != nil {
		return "", err
	}

	// unzip
	err = Unzip(zipPath, ".cache")
	if err != nil {
		return "", err
	}

	fmt.Println("Go installed in cache ")

	fmt.Println("GoPath", GoPath())
	return GoPath(), nil
}

func RunGoBuild(projectDir string) error {

	goPath, err := EnsureGo()
	if err != nil {
		fmt.Println("err-go", err)
		return err
	}

	err = run(goPath, projectDir, "mod", "init", "backforge")

	if err != nil {
		fmt.Println(err)
	}
	err = run(goPath, projectDir, "get", "-u", "github.com/go-chi/chi/v5")
	if err != nil {
		fmt.Println(err)
	}
	err = run(goPath, projectDir, "get", "github.com/go-chi/cors")
	if err != nil {
		fmt.Println(err)
	}
	err = run(goPath, projectDir, "get", "modernc.org/sqlite")
	if err != nil {
		fmt.Println(err)
	}

	err = run(goPath, projectDir, "fmt", "./...")
	if err != nil {
		fmt.Println(err)
	}

	buildPath := "./bin/app.exe"
	os := runtime.GOOS
	if os != "windows" {
		buildPath = "./bin/app"
	}

	err = run(goPath, projectDir, "build", "-o", buildPath, "./cmd")
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func run(goPath, projectDir string, args ...string) error {
	cmd := exec.Command(goPath, args...)

	cmd.Dir = projectDir

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()

}
