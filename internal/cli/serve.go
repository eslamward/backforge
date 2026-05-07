package cli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func Serve() {
	fmt.Println("loading the server....")
	path := "./output/bin/app.exe"
	if runtime.GOOS != "windows" {
		fmt.Println("not windows")
		path = "./output/bin/app"
	}
	cmd := exec.Command(path)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
