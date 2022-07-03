package config

import (
	"bytes"
	"encoding/base64"
	"io"
)

type QRCode string

type Pix struct {
	Key         string `fig:"key" yaml:"key"`
	Description string `fig:"description" yaml:"description"`
	QRCode      QRCode `fig:"qrcode" yaml:"qrcode"`
}

func (q QRCode) NewBuffer() (io.Reader, error) {
	dec, err := base64.StdEncoding.DecodeString(string(q))
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(dec), nil
}
