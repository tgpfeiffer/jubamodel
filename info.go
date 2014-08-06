package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"unsafe"
)

type Model struct {
	Path       string             `json:"path"`
	Header     Header             `json:"header"`
	SystemData []SystemDataHeader `json:"-"`
	UserData   []UserDataHeader   `json:"-"`
}

type Header struct {
	Magic          string `json:"Magic"`
	FormatVersion  uint64 `json:"format_version"`
	JubatusVersion string `json:"jubatus_version"`
	CRC32          uint32 `json:"crc32"`
	SystemDataSize uint64 `json:"system_data_size"`
	UserDataSize   uint64 `json:"user_data_size"`

	Raw []byte `json:"-"`
}

type SystemDataHeader struct {
	Version   string                 `json:"version"`
	Timestamp int64                  `json:"timestamp"`
	Type      string                 `json:"type"`
	ID        string                 `json:"id"`
	Config    map[string]interface{} `json:"config"`
}

type UserDataHeader struct {
	Version uint64 `json:"version"`
}

func Info(paths []string) ([]*Model, error) {
	res := []*Model{}
	for _, p := range paths {
		m, err := info(p)
		if err != nil { // TODO: Add option to ignore all files which are not jubatus models.
			fmt.Fprintln(os.Stderr, "Cannot read a model file:", err)
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

func info(path string) (*Model, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	type BinaryHeader struct {
		Magic                        [8]byte
		FormatVersion                uint64
		Major, Minor, Maintenance    uint32
		CRC32                        uint32
		SystemDataSize, UserDataSize uint64
	}

	headerBuffer := make([]byte, unsafe.Sizeof(BinaryHeader{}))
	if n, err := f.Read(headerBuffer); err != nil {
		return nil, err
	} else if n < len(headerBuffer) {
		return nil, fmt.Errorf("the file is too small")
	}

	bh := BinaryHeader{}
	if err := binary.Read(bytes.NewReader(headerBuffer), binary.BigEndian, &bh); err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
	}
	m := &Model{
		Path: absPath,
	}
	header := &m.Header
	header.Magic = string(bh.Magic[:])
	header.FormatVersion = bh.FormatVersion
	header.JubatusVersion = fmt.Sprint(bh.Major, ".", bh.Minor, ".", bh.Maintenance)
	header.CRC32 = bh.CRC32
	header.SystemDataSize = bh.SystemDataSize
	header.UserDataSize = bh.UserDataSize
	header.Raw = headerBuffer

	// TODO: read containers
	return m, nil
}