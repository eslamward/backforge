package main

import (
	"fmt"
	"os"

	"github.com/eslamward/backforge/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: backforge <command>")
		return
	}

	switch os.Args[1] {

	case "build":
		cli.Build()
	case "serve":
		cli.Serve()
	default:
		fmt.Println("Unknown command")
		fmt.Println("Use >> build to build the server")
		fmt.Println("Use >> serve to run the server")

	}
}
