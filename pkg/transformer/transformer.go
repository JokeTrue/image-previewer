package transformer

import (
	"bytes"

	"github.com/disintegration/imaging"
)

type Transformer interface {
	Crop(img []byte, width, height int) ([]byte, error)
}

type transformer struct{}

// nolint
func NewTransformer() *transformer {
	return &transformer{}
}

func (t *transformer) Crop(img []byte, width, height int) ([]byte, error) {
	src, err := imaging.Decode(bytes.NewReader(img))
	if err != nil {
		return nil, err
	}
	src = imaging.Fill(src, width, height, imaging.Center, imaging.Lanczos)

	var buff bytes.Buffer
	err = imaging.Encode(&buff, src, imaging.JPEG)
	return buff.Bytes(), err
}
