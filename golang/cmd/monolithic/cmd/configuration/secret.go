package configuration

import (
	mntreamerConfiguration "mntreamer/shared/configuration"
)

type Variable struct {
	Database *mntreamerConfiguration.Database
	Api      *mntreamerConfiguration.Api
	Frontend *mntreamerConfiguration.Api
}

func NewVariable() *Variable {
	return &Variable{
		Database: &mntreamerConfiguration.Database{
			Uri:      "127.0.0.1:11001",
			Username: "root",
			Password: "root",
		},
		Api: &mntreamerConfiguration.Api{
			Ip:   "localhost",
			Port: 11000,
		},
		Frontend: &mntreamerConfiguration.Api{
			Ip:   "localhost",
			Port: 11002,
		},
	}
}
