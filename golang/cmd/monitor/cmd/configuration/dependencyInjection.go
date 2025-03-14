package configuration

import (
	"mntreamer/shared/configuration"
)

func NewContainer() configuration.IContainer {
	ctnr := NewMonolithicContainer(nil)
	return ctnr
}
