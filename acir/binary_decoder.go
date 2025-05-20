package acir

import (
	"encoding/binary"
	"io"
)

func DecodeUint8(reader io.Reader) (uint8, error) {
	var value uint8
	err := binary.Read(reader, binary.LittleEndian, &value)
	return value, err
}

func DecodeUint16(reader io.Reader) (uint16, error) {
	var value uint16
	err := binary.Read(reader, binary.LittleEndian, &value)
	return value, err
}

func DecodeUint32(reader io.Reader) (uint32, error) {
	var value uint32
	err := binary.Read(reader, binary.LittleEndian, &value)
	return value, err
}

func DecodeUint64(reader io.Reader) (uint64, error) {
	var value uint64
	err := binary.Read(reader, binary.LittleEndian, &value)
	return value, err
}
