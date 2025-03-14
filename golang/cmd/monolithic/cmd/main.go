package main

import (
	"mntreamer/monolithic/cmd/configuration"
	"mntreamer/monolithic/cmd/lib"
)

func main() {
	startLib()
}

func startLib() {
	ctnr := configuration.NewMonolithicContainer()
	lib.Start(ctnr)
}
