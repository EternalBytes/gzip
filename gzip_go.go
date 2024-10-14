package gzip_go

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
)

type Gzipgo struct {
	Compressed   string
	Decompressed string
}

func NewGzip() *Gzipgo {
	return &Gzipgo{}
}

func (gz *Gzipgo) Compress(str string) error {
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
		return err
	}

	gz.Compressed = base64.StdEncoding.EncodeToString(zpd.Bytes())
	return nil
}

func (gz *Gzipgo) Decompress(base64Str string) error {
	chanErr := make(chan error)
	var strUn bytes.Buffer
	var str bytes.Buffer
	decoded, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return err
	}

	_, err = str.Write(decoded)
	if err != nil {
		return err
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
		return err
	}

	gz.Decompressed = strUn.String()
	return nil
}
