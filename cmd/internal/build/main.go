package main

import (
	"go/build"
	"os"
	"os/exec"

	"github.com/tim06/sing-box/cmd/internal/build_shared"
	"github.com/tim06/sing-box/log"
)

func main() {
	build_shared.FindSDK()

	if os.Getenv("GOPATH") == "" {
		os.Setenv("GOPATH", build.Default.GOPATH)
	}

	command := exec.Command(os.Args[1], os.Args[2:]...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		log.Fatal(err)
	}
}
