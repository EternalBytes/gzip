package gzip

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
)

type gzipgo struct {
	Compressed   string
	Decompressed string
}

// NewGzip initializes a Gzipgo object
func NewGzip() *gzipgo {
	return &gzipgo{}
}

// Compress, compress a string and returns a Gzipgo pointer and an error if there is one
func (gz *gzipgo) Compress(str string) (*gzipgo, error) {
	chanErr := make(chan error)
	var gzp *gzip.Writer
	var zpd bytes.Buffer

	go func() {
		gzp = gzip.NewWriter(&zpd)
		_, err := gzp.Write([]byte(str))
		if err != nil {
			chanErr <- err
			return
		}

		err = gzp.Close()
		if err != nil {
			chanErr <- err
			return
		}

		chanErr <- nil
	}()
	if err := <-chanErr; err != nil {
		return nil, err
	}

	gz.Compressed = base64.StdEncoding.EncodeToString(zpd.Bytes())
	return gz, nil
}

// Decompress, decompresses a string and returns a Gzipgo pointer and an error if there is one
func (gz *gzipgo) Decompress(base64Str string) (*gzipgo, error) {
	chanErr := make(chan error)
	var strUn bytes.Buffer
	var str bytes.Buffer
	decoded, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}

	_, err = str.Write(decoded)
	if err != nil {
		return nil, err
	}

	go func() {
		gzp, err := gzip.NewReader(&str)
		if err != nil {
			chanErr <- err
			return
		}

		err = gzp.Close()
		if err != nil {
			chanErr <- err
			return
		}

		_, err = strUn.ReadFrom(gzp)
		if err != nil {
			chanErr <- err
			return
		}
		chanErr <- nil
	}()

	if err := <-chanErr; err != nil {
		return nil, err
	}

	gz.Decompressed = strUn.String()
	return gz, nil
}

// GetComp, returns the compressed string
func (gz *gzipgo) GetComp() string {
	return gz.Compressed
}

// GetDecomp, returns the decompressed string
func (gz *gzipgo) GetDecomp() string {
	return gz.Decompressed
}
