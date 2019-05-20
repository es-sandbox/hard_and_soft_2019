package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
)

// JPG |640 * 480|

const (
	imageFilepath    = "example.jpg"
	rawImageFilepath = "example.jpg.base64"
	finalImage       = "final.jpg"
)

func decodeImage() {
	rawImage, err := ioutil.ReadFile(imageFilepath)
	if err != nil {
		log.Fatalf("can't get raw image: %v\n", err)
	}
	rawImageReader := bytes.NewReader(rawImage)

	rawImageBase64 := base64.StdEncoding.EncodeToString(rawImage)
	if err := ioutil.WriteFile(rawImageFilepath, []byte(rawImageBase64), 0600); err != nil {
		log.Fatalf("can't dump raw image: %v\n", err)
	}

	image, err := jpeg.Decode(rawImageReader)
	if err != nil {
		log.Fatalf("can't decode raw image: %v\n", err)
	}

	_ = image
}

func encodeImage() {
	rawImageBase64, err := ioutil.ReadFile(rawImageFilepath)
	if err != nil {
		log.Fatalf("can't read base64-encoded raw image")
	}

	rawImage, err := base64.StdEncoding.DecodeString(string(rawImageBase64))
	if err != nil {
		log.Fatal(err)
	}
	rawImageReader := bytes.NewReader(rawImage)

	image, err := jpeg.Decode(rawImageReader)
	if err != nil {
		log.Fatalf("can't decode raw image: %v\n", err)
	}

	file, err := os.OpenFile(finalImage, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("can't open file: %v\n", err)
	}
	err = jpeg.Encode(
		file,
		image,
		&jpeg.Options{Quality: jpeg.DefaultQuality},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	encode := flag.Bool("encode", false, "")
	decode := flag.Bool("decode", false, "")
	flag.Parse()

	if *encode {
		encodeImage()
	}

	if *decode {
		decodeImage()
	}
}
