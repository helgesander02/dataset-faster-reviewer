package utils

import (
	"bytes"
	"encoding/base64"
	"image"
	"log"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
)

type ImageProcessor struct {
	maxWidth int
	quality  float32
}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{
		maxWidth: 150,
		quality:  75,
	}
}

func (ip *ImageProcessor) CompressToBase64(imgPath string) (string, error) {
	log.Println("Processing image:", imgPath)

	file, err := os.Open(imgPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	decodedImg, format, err := image.Decode(file)
	if err != nil {
		log.Printf("Failed to decode image (format: %s): %v\n", format, err)
		return "", err
	}

	resizedImg := resize.Resize(uint(ip.maxWidth), 0, decodedImg, resize.Lanczos3)

	var buf bytes.Buffer
	opts := &webp.Options{Lossless: false, Quality: ip.quality}
	if err := webp.Encode(&buf, resizedImg, opts); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
