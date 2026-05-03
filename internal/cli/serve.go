package cli

import (
	"fmt"
	"os"
	"os/exec"
)

func Serve() {
	fmt.Println("loading the server....")
	path := "./output/bin/app.exe"
	cmd := exec.Command(path)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
