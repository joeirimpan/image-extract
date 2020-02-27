package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func readFileHeader(reader io.Reader) (bool, error) {
	var isFixed bool

	// File ID
	{
		p := make([]byte, 15)
		_, err := reader.Read(p)
		if err != nil {
			return isFixed, err
		}
	}

	// Request ID
	{
		p := make([]byte, 15)
		_, err := reader.Read(p)
		if err != nil {
			return isFixed, err
		}
	}

	// File Version
	{
		p := make([]byte, 4)
		_, err := reader.Read(p)
		if err != nil {
			return isFixed, err
		}
	}

	// File creation date
	{
		p := make([]byte, 8)
		_, err := reader.Read(p)
		if err != nil {
			return isFixed, err
		}
	}

	// File creation time
	{
		p := make([]byte, 6)
		_, err := reader.Read(p)
		if err != nil {
			return isFixed, err
		}
	}

	// Number of check records
	{
		p := make([]byte, 6)
		_, err := reader.Read(p)
		if err != nil {
			return isFixed, err
		}
	}

	// Record size
	{
		p := make([]byte, 4)
		_, err := reader.Read(p)
		if err != nil {
			return isFixed, err
		}

		// Read filler
		switch string(p) {
		// Fixed
		case "0256":
			isFixed = true

			filler := make([]byte, 194)
			_, err := reader.Read(filler)
			if err != nil {
				return isFixed, err
			}
		// Variable
		case "0090":
			filler := make([]byte, 28)
			_, err := reader.Read(filler)
			if err != nil {
				return isFixed, err
			}
		}
	}

	return isFixed, nil
}

func readCheckIndex(reader io.Reader, isFixed bool) (int, error) {
	var numImages int
	// Bank number
	{
		p := make([]byte, 4)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}
	}

	// Routing transit number
	{
		p := make([]byte, 9)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}
	}

	// Account number
	{
		p := make([]byte, 20)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}
	}

	// Check number
	{
		p := make([]byte, 15)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}
	}

	// Amount
	{
		p := make([]byte, 10)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}

	}

	// Seq number
	{
		p := make([]byte, 15)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}
	}

	// Posted date
	{
		p := make([]byte, 8)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}
	}

	// Number of images
	{
		p := make([]byte, 4)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}

		numImages, err = strconv.Atoi(string(p))
		if err != nil {
			return numImages, err
		}
	}

	// Filler
	{
		var fillerLength = 1
		if isFixed {
			fillerLength = 167
		}

		p := make([]byte, fillerLength)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}
	}

	return numImages, nil
}

func readImageHeader(reader io.Reader, isFixed bool) (int, error) {
	var numImages int
	// Image type
	{
		p := make([]byte, 4)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}
	}

	// Image side
	{
		p := make([]byte, 1)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}
	}

	// Number of records
	{
		p := make([]byte, 4)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}

		numImages, err = strconv.Atoi(string(p))
		if err != nil {
			return numImages, err
		}
	}

	// Image data record length
	{
		p := make([]byte, 6)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}
	}

	// Filler
	{
		var fillerLength = 71
		if isFixed {
			fillerLength = 231
		}

		p := make([]byte, fillerLength)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}
	}

	return numImages, nil
}

func readImageData(reader io.Reader, isFixed bool) ([]byte, error) {
	var (
		recLength int
		imageData []byte
	)
	// Record length
	{
		p := make([]byte, 4)
		_, err := reader.Read(p)
		if err != nil {
			return imageData, err
		}

		recLength, err = strconv.Atoi(string(p))
		if err != nil {
			return imageData, err
		}
	}

	// Image data
	{
		if isFixed {
			recLength = 246
		}

		imageData = make([]byte, recLength)
		_, err := reader.Read(imageData)
		if err != nil {
			return imageData, err
		}
	}

	return imageData, nil
}

func readCheckTrailer(reader io.Reader, isFixed bool) error {
	// Filler
	var fillerLength = 86
	if isFixed {
		fillerLength = 252
	}

	p := make([]byte, fillerLength)
	_, err := reader.Read(p)
	if err != nil {
		return err
	}

	return nil
}

func readOutputTrailer(reader io.Reader, isFixed bool) error {
	// File ID
	{
		p := make([]byte, 15)
		_, err := reader.Read(p)
		if err != nil {
			return err
		}
	}

	// Request ID
	{
		p := make([]byte, 15)
		_, err := reader.Read(p)
		if err != nil {
			return err
		}
	}

	// File version
	{
		p := make([]byte, 4)
		_, err := reader.Read(p)
		if err != nil {
			return err
		}
	}

	// File creation date
	{
		p := make([]byte, 8)
		_, err := reader.Read(p)
		if err != nil {
			return err
		}
	}

	// File creation time
	{
		p := make([]byte, 6)
		_, err := reader.Read(p)
		if err != nil {
			return err
		}
	}

	// Number of detail records
	{
		p := make([]byte, 6)
		_, err := reader.Read(p)
		if err != nil {
			return err
		}
	}

	// Filler
	{
		var fillerLength = 32
		if isFixed {
			fillerLength = 198
		}

		p := make([]byte, fillerLength)
		_, err := reader.Read(p)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	fs, err := os.Open("./data.dat")
	if err != nil {
		log.Fatalf("error reading dat file")
	}
	defer fs.Close()

	var (
		isFixed bool
		imgCt   int
	)
	reader := bufio.NewReader(fs)
	hdrBuf := make([]byte, 4)
	for {
		_, err := reader.Read(hdrBuf)
		if err == io.EOF {
			break
		}

		switch string(hdrBuf) {
		case "1200":
			// file header
			isFixed, err = readFileHeader(reader)
			if err != nil {
				log.Fatalf("error while reading file header: %v", err)
			}
		case "1201":
			// check index
			_, err := readCheckIndex(reader, isFixed)
			if err != nil {
				log.Fatalf("error while reading check index: %v", err)
			}
		case "1202":
			// image header
			_, err := readImageHeader(reader, isFixed)
			if err != nil {
				log.Fatalf("error while reading image header: %v", err)
			}
		case "1203":
			imgCt++
			// image data
			image, err := readImageData(reader, isFixed)
			if err != nil {
				log.Fatalf("error while reading image data: %v", err)
			}

			file, err := os.Create(fmt.Sprintf("image_%d", imgCt))
			if err != nil {
				log.Fatalf("error while creating image file: %v", err)
			}
			defer file.Close()

			if _, err := file.Write(image); err != nil {
				log.Fatalf("error while writing to image file: %v", err)
			}
		case "1204":
			// check trailer
			err := readCheckTrailer(reader, isFixed)
			if err != nil {
				log.Fatalf("error while reading check trailer: %v", err)
			}
		case "1209":
			// file trailer
			err := readOutputTrailer(reader, isFixed)
			if err != nil {
				log.Fatalf("error while reading file trailer: %v", err)
			}
		}
	}
}
