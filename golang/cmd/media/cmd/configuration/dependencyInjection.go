package configuration

import (
	"mntreamer/shared/configuration"
)

func NewContainer() configuration.IContainer {
	return NewMonolithicContainer(nil)
}
