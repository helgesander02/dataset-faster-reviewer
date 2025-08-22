package utils

import (
	"bytes"
	"encoding/base64"
	"image"
	"log"
	"os"
	"runtime"
	"sync"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
)

var maxWorkers = runtime.NumCPU()

func CompressImageSetToBase64(imagpathSet []string) []string {
	var base64Images []string
	
	sem := make(chan struct{}, maxWorkers)
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	base64Images = make([]string, len(imagpathSet))
	
	for i, imgPath := range imagpathSet {
		wg.Add(1)
		go func(index int, path string) {
			defer wg.Done()
			
			sem <- struct{}{}        
			defer func() { <-sem }() 
			
			base64Image, err := CompressImageToBase64(path)
			if err != nil {
				log.Printf("Failed to compress image %s: %v", path, err)
				base64Image = ""
			}
			
			mu.Lock()
			base64Images[index] = base64Image
			mu.Unlock()
		}(i, imgPath)
	}
	
	wg.Wait()
	return base64Images
}

func CompressImageToBase64(imgPath string) (string, error) {
	var (
		maxWidth int     = 400
		quality  float32 = 75
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

	resizedImg := resize.Resize(uint(maxWidth), 0, decodedImg, resize.Bilinear)
	
	decodedImg = nil
	runtime.GC()

	var buf bytes.Buffer
	opts := &webp.Options{Lossless: false, Quality: quality}
	if err := webp.Encode(&buf, resizedImg, opts); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
