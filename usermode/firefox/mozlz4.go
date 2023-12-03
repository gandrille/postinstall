package firefox

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pierrec/lz4"
)

// COPY/PASTE from
// https://github.com/frioux/leatherman/blob/master/pkg/mozlz4/mozlz4_test.go
// Only error management updates

const magicHeader = "mozLz40\x00"

// NewDecompressReader returns an io.Reader that decompresses the data from r.
func NewDecompressReader(r io.Reader) (io.Reader, error) {
	header := make([]byte, len(magicHeader))
	_, err := r.Read(header)
	if err != nil {
		return nil, errors.New("Couldn't read header: " + err.Error())
	}
	if string(header) != magicHeader {
		return nil, errors.New("mozLz4 header not found")
	}

	var size uint32
	err = binary.Read(r, binary.LittleEndian, &size)
	if err != nil {
		return nil, errors.New("Couldn't read size: " + err.Error())
	}

	src, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.New("Couldn't read compressed data: " + err.Error())
	}

	out := make([]byte, size)
	sz, err := lz4.UncompressBlock(src, out)

	if err != nil {
		return nil, errors.New("Couldn't decompress data: " + err.Error())
	}
	// This could maybe be a warning or ignored entirely
	if sz != int(size) {
		return nil, errors.New("Header size expected " + fmt.Sprint(size) + ", but got " + fmt.Sprint(sz))
	}

	return bytes.NewReader(out), nil
}

// Compress data
func Compress(src io.Reader, dst io.Writer, intendedSize int) error {
	_, err := dst.Write([]byte(magicHeader))
	if err != nil {
		return errors.New("couldn't Write header: " + err.Error())
	}
	b, err := ioutil.ReadAll(src)
	if err != nil {
		return errors.New("couldn't ReadAll to Compress: " + err.Error())
	}

	err = binary.Write(dst, binary.LittleEndian, uint32(intendedSize))
	if err != nil {
		return errors.New("couldn't encode length: " + err.Error())
	}
	dstBytes := make([]byte, 10*len(b))
	sz, err := lz4.CompressBlockHC(b, dstBytes, -1)
	if err != nil {
		return errors.New("couldn't CompressBlock: " + err.Error())
	}
	if sz == 0 {
		return errors.New("data incompressible")
	}
	_, err = dst.Write(dstBytes[:sz])
	if err != nil {
		return errors.New("couldn't Write compressed data: " + err.Error())
	}

	return nil
}
