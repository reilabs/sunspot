package acir

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
)

func DecodeProgramBytecode(bytecode string) (reader io.Reader, err error) {
	// Decode the bytecode from base64 string to bytes
	data, err := base64.RawStdEncoding.DecodeString(bytecode)
	if err != nil {
		return nil, fmt.Errorf("failed to decode bytecode: %w", err)
	}

	// Decompress the bytecode using gzip
	reader, err = gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}

	return reader, err
}
