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

func readCheckIndex(reader io.Reader, isFixed bool) (int, string, error) {
	var (
		numImages int
		checkNo   string
	)

	// Bank number
	{
		p := make([]byte, 4)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, checkNo, err
		}
	}

	// Routing transit number
	{
		p := make([]byte, 9)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, checkNo, err
		}
	}

	// Account number
	{
		p := make([]byte, 20)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, checkNo, err
		}
	}

	// Check number
	{
		p := make([]byte, 15)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, checkNo, err
		}

		checkNo = string(p)
	}

	// Amount
	{
		p := make([]byte, 10)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, checkNo, err
		}

	}

	// Seq number
	{
		p := make([]byte, 15)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, checkNo, err
		}
	}

	// Posted date
	{
		p := make([]byte, 8)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, checkNo, err
		}
	}

	// Number of images
	{
		p := make([]byte, 4)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, checkNo, err
		}

		numImages, err = strconv.Atoi(string(p))
		if err != nil {
			return numImages, checkNo, err
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
			return numImages, checkNo, err
		}
	}

	return numImages, checkNo, nil
}

func readImageHeader(reader io.Reader, isFixed bool) (int, string, error) {
	var (
		numImages int
		imgType   string
	)
	// Image type
	{
		p := make([]byte, 4)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, imgType, err
		}

		imgType = string(p)
	}

	// Image side
	{
		p := make([]byte, 1)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, imgType, err
		}
	}

	// Number of records
	{
		p := make([]byte, 4)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, imgType, err
		}

		numImages, err = strconv.Atoi(string(p))
		if err != nil {
			return numImages, imgType, err
		}
	}

	// Image data record length
	{
		p := make([]byte, 6)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, imgType, err
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
			return numImages, imgType, err
		}
	}

	return numImages, imgType, nil
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

// TODO: Make sure that we read the correct number of records as specified in the headers.
func main() {
	if len(os.Args) < 2 {
		log.Fatalf("missing required arguments")
	}

	fs, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("error reading dat file")
	}
	defer fs.Close()

	var (
		isFixed bool
		checkNo string
		imgType string
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
		// file header
		case "1200":
			isFixed, err = readFileHeader(reader)
			if err != nil {
				log.Fatalf("error while reading file header: %v", err)
			}

		// check index
		case "1201":
			_, checkNo, err = readCheckIndex(reader, isFixed)
			if err != nil {
				log.Fatalf("error while reading check index: %v", err)
			}

		// image header
		case "1202":
			_, imgType, err = readImageHeader(reader, isFixed)
			if err != nil {
				log.Fatalf("error while reading image header: %v", err)
			}

		// image data
		case "1203":
			imgCt++
			image, err := readImageData(reader, isFixed)
			if err != nil {
				log.Fatalf("error while reading image data: %v", err)
			}

			file, err := os.Create(fmt.Sprintf("%s.%s", checkNo, imgType))
			if err != nil {
				log.Fatalf("error while creating image file: %v", err)
			}

			if _, err := file.Write(image); err != nil {
				log.Fatalf("error while writing to image file: %v", err)
			}

			file.Close()

		// check trailer
		case "1204":
			err := readCheckTrailer(reader, isFixed)
			if err != nil {
				log.Fatalf("error while reading check trailer: %v", err)
			}

		// file trailer
		case "1209":
			err := readOutputTrailer(reader, isFixed)
			if err != nil {
				log.Fatalf("error while reading file trailer: %v", err)
			}
		}
	}
}
