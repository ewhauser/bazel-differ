package internal

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"os"
	"path"
	"strings"
)

type BazelSourceFileTarget interface {
	Name() *string
	Digest() []byte
}

type bazelSourceFileTarget struct {
	name   *string
	digest []byte
}

func NewBazelSourceFileTarget(name string, digest []byte, workingDirectory string) (BazelSourceFileTarget, error) {
	finalDigest := bytes.NewBuffer([]byte{})
	if workingDirectory != "" && strings.HasPrefix(name, "//") {
		filenameSubstring := name[2:]
		filenamePath := strings.Replace(filenameSubstring, ":", "/", 1)
		sourceFile := path.Join(workingDirectory, filenamePath)
		if _, err := os.Stat(sourceFile); !errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does not exist
			contents, err := os.ReadFile(sourceFile)
			if err != nil {
				return nil, err
			}
			finalDigest.Write(contents)
		}
	}
	finalDigest.Write(digest)
	finalDigest.Write([]byte(name))
	checksum := sha256.Sum256(finalDigest.Bytes())
	return &bazelSourceFileTarget{
		name:   &name,
		digest: checksum[:],
	}, nil
}

func (b *bazelSourceFileTarget) Name() *string {
	return b.name
}

func (b *bazelSourceFileTarget) Digest() []byte {
	return b.digest
}
