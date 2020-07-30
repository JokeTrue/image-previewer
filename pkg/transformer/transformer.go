package cropper

import (
	"bytes"
	"image/jpeg"

	"github.com/disintegration/imaging"
)

type Transformer interface {
	Crop(img []byte, width, height int) ([]byte, error)
}

type Cropper struct{}

func NewCropper() *Cropper {
	return &Cropper{}
}

func (t *Cropper) Crop(img []byte, width, height int) ([]byte, error) {
	src, err := imaging.Decode(bytes.NewReader(img))
	if err != nil {
		return nil, err
	}
	src = imaging.Fill(src, width, height, imaging.Center, imaging.Lanczos)

	var buff bytes.Buffer
	err = jpeg.Encode(&buff, src, nil)
	return buff.Bytes(), err
}
