package main

import "mntreamer/monitor/cmd/configuration"

func main() {
	startLib()
}

func startLib() {
	ctnr := configuration.NewMonolithicContainer()
	lib.Start(ctnr)
}
