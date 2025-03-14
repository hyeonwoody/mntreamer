package main

import (
	"mntreamer/platform/cmd/configuration"
	"mntreamer/platform/cmd/lib"
)

func main() {
	startLib()
}

func startLib() {
	ctnr := configuration.NewContainer()
	lib.Start(ctnr)
}
