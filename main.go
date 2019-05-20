package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// JPG |640 * 480|

const (
	imageFilepath    = "images/example.jpg"
	rawImageFilepath = "images/example.jpg.base64"
	finalImage       = "images/final.jpg"
)

var httpAddress *string

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

func enableHttp() {
	fmt.Printf("Start HTTP Server\n")

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintf(w, "Welcome to my website!"); err != nil {
			log.Println(err)
			return
		}
	})

	http.HandleFunc("/api/v1/image", func (w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			rawImageBase64, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
				return
			}

			rawImage, err := base64.StdEncoding.DecodeString(string(rawImageBase64))
			if err != nil {
				log.Fatal(err)
				return
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
	})

	fs := http.FileServer(http.Dir("images/"))
	http.Handle("/images/", http.StripPrefix("/images/", fs))

	if err := http.ListenAndServe(*httpAddress, nil); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}

func enableHttpClient() {
	url := "http://localhost:8080/api/v1/image"
	contentType := "image/jpeg"

	rawImage, err := ioutil.ReadFile(imageFilepath)
	if err != nil {
		log.Fatalf("can't get raw image: %v\n", err)
	}

	rawImageBase64 := base64.StdEncoding.EncodeToString(rawImage)
	rawImageBase64Reader := bytes.NewReader([]byte(rawImageBase64))

	resp, err := http.Post(url, contentType, rawImageBase64Reader)
	if err != nil {
		log.Fatal(err)
	}

	_ = resp
}

func main() {
	encode := flag.Bool("encode", false, "")
	decode := flag.Bool("decode", false, "")
	enableHttpFlag := flag.Bool("http", false, "")
	enableHttpClientFlag := flag.Bool("http-client", false, "")
	httpAddress = flag.String("http-address", "0.0.0.0:8080", "")

	flag.Parse()

	if *encode {
		encodeImage()
	}

	if *decode {
		decodeImage()
	}

	if *enableHttpFlag {
		enableHttp()
	}

	if *enableHttpClientFlag {
		enableHttpClient()
	}
}
