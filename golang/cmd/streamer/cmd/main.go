package main

import (
	"mntreamer/streamer/cmd/configuration"
	"mntreamer/streamer/cmd/lib"

)

func main() {
	startLib()
}

func startLib() {
	ctnr := configuration.NewContainer()
	lib.Start(ctnr)
}
