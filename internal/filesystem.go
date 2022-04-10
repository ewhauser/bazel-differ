package internal

import (
	"io/fs"
	"os"
)

var Filesystem = &filesystem{}

type filesystem struct {
}

func (f filesystem) Open(name string) (fs.File, error) {
	return os.Open(name)
}
