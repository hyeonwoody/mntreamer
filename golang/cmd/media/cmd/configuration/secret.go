package configuration

import (
	mntreamerConfiguration "mntreamer/shared/configuration"
)

type Variable struct {
	Database *mntreamerConfiguration.Database
	Api      *mntreamerConfiguration.Api
	RootPath string
}

func NewVariable() *Variable {
	return &Variable{
		Database: &mntreamerConfiguration.Database{
			Uri:      "127.0.0.1:11001",
			Username: "root",
			Password: "root",
		},
		RootPath: "/zzz/mntreamer",
	}
}
