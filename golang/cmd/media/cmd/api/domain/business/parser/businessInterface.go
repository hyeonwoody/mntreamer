package business

import (
	"io"
	"mntreamer/media/cmd/model"
)

type IBusiness interface {
	Decode(reader io.Reader) (interface{}, error)
	Encode(interface{}) model.IBuffer
}
