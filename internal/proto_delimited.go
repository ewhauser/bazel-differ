package internal

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
)

type ProtoDelimitedReader struct {
	buf *bufio.Reader
}

func NewReader(r io.Reader) *ProtoDelimitedReader {
	return &ProtoDelimitedReader{
		buf: bufio.NewReader(r),
	}
}

func (r *ProtoDelimitedReader) Next() ([]byte, error) {
	size, err := binary.ReadUvarint(r.buf)
	if err != nil {
		return nil, err
	}
	data := make([]byte, size)
	if _, err := io.ReadFull(r.buf, data); err != nil {
		return nil, err
	}
	return data, nil
}

func (r *ProtoDelimitedReader) ReadTargets() ([]*Target, error) {
	var targets []*Target
	var err error
	for buf, err := r.Next(); err == nil; buf, err = r.Next() {
		var target Target
		if err := proto.Unmarshal(buf, &target); err != nil {
			return nil, fmt.Errorf("failed to unmarshal Target message: %w", err)
		}
		targets = append(targets, &target)
	}
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("error while reading stdout from bazel command: %w", err)
	}
	return targets, nil
}
