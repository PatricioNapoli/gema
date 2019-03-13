package utils

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
)

func Gunzip(reader io.ReadCloser, to* []byte) error {
	reader, err := gzip.NewReader(reader)
	*to, err = ioutil.ReadAll(reader)
	err = reader.Close()
	return err
}


func GZip(from* []byte, to* []byte) error {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	_, err := gz.Write(*from)
	if err != nil {
		return err
	}

	if err = gz.Flush(); err != nil {
		return err
	}

	if err = gz.Close(); err != nil {
		return err
	}

	*to = b.Bytes()

	return nil
}
