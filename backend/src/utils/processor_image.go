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

func CompressImageSetToBase64(imagpathSet []string) []string {
	var base64Images []string
	for _, imgPath := range imagpathSet {
		base64Image, err := CompressImageToBase64(imgPath)
		if err != nil {
			log.Printf("Failed to compress image %s: %v", imgPath, err)
			continue
		}
		base64Images = append(base64Images, base64Image)
	}
	return base64Images
}

func CompressImageToBase64(imgPath string) (string, error) {
	var (
		maxWidth int     = 150
		quality  float32 = 90
	)
	log.Printf("Processing image: %s (maxWidth:%v, quality:%v)", imgPath, maxWidth, quality)

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

	resizedImg := resize.Resize(uint(maxWidth), 0, decodedImg, resize.Lanczos3)

	var buf bytes.Buffer
	opts := &webp.Options{Lossless: false, Quality: quality}
	if err := webp.Encode(&buf, resizedImg, opts); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
