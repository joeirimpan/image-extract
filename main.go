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
		fmt.Println(string(p))
	}

	// Request ID
	{
		p := make([]byte, 15)
		_, err := reader.Read(p)
		if err != nil {
			return isFixed, err
		}
		fmt.Println(string(p))
	}

	// File Version
	{
		p := make([]byte, 4)
		_, err := reader.Read(p)
		if err != nil {
			return isFixed, err
		}
		fmt.Println(string(p))
	}

	// File creation date
	{
		p := make([]byte, 8)
		_, err := reader.Read(p)
		if err != nil {
			return isFixed, err
		}
		fmt.Println(string(p))
	}

	// File creation time
	{
		p := make([]byte, 6)
		_, err := reader.Read(p)
		if err != nil {
			return isFixed, err
		}
		fmt.Println(string(p))
	}

	// Number of check records
	{
		p := make([]byte, 6)
		_, err := reader.Read(p)
		if err != nil {
			return isFixed, err
		}
		fmt.Printf("Records: %s\n", string(p))

		numCheckRec, err := strconv.Atoi(string(p))
		if err != nil {
			return isFixed, err
		}
		fmt.Printf("Check Records: %d\n", numCheckRec)
	}

	// Record size
	{
		p := make([]byte, 4)
		_, err := reader.Read(p)
		if err != nil {
			return isFixed, err
		}
		fmt.Println(string(p))

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
		fmt.Println(string(p))
	}

	// Routing transit number
	{
		p := make([]byte, 9)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}
		fmt.Println(string(p))
	}

	// Account number
	{
		p := make([]byte, 20)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}
		fmt.Println(string(p))
	}

	// Check number
	{
		p := make([]byte, 15)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}
		fmt.Println(string(p))
	}

	// Amount
	{
		p := make([]byte, 10)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}
		fmt.Println(string(p))
	}

	// Seq number
	{
		p := make([]byte, 15)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}
		fmt.Println(string(p))
	}

	// Posted date
	{
		p := make([]byte, 8)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}
		fmt.Println(string(p))
	}

	// Number of images
	{
		p := make([]byte, 4)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}
		fmt.Println(string(p))

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
		fmt.Println(string(p))
	}

	// Image side
	{
		p := make([]byte, 1)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}
		fmt.Println(string(p))
	}

	// Number of records
	{
		p := make([]byte, 4)
		_, err := reader.Read(p)
		if err != nil {
			return numImages, err
		}
		fmt.Println(string(p))

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
		fmt.Println(string(p))
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

func readImageData(reader io.Reader, isFixed bool) error {
	var recLength int
	// Record length
	{
		p := make([]byte, 4)
		_, err := reader.Read(p)
		if err != nil {
			return err
		}
		fmt.Println(string(p))

		recLength, err = strconv.Atoi(string(p))
		if err != nil {
			return err
		}
	}

	// Image data
	{
		if isFixed {
			recLength = 246
		}

		p := make([]byte, recLength)
		_, err := reader.Read(p)
		if err != nil {
			return err
		}
		fmt.Println(p)
	}

	return nil
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
		fmt.Println(string(p))
	}

	// Request ID
	{
		p := make([]byte, 15)
		_, err := reader.Read(p)
		if err != nil {
			return err
		}
		fmt.Println(string(p))
	}

	// File version
	{
		p := make([]byte, 4)
		_, err := reader.Read(p)
		if err != nil {
			return err
		}
		fmt.Println(string(p))
	}

	// File creation date
	{
		p := make([]byte, 8)
		_, err := reader.Read(p)
		if err != nil {
			return err
		}
		fmt.Println(string(p))
	}

	// File creation time
	{
		p := make([]byte, 6)
		_, err := reader.Read(p)
		if err != nil {
			return err
		}
		fmt.Println(string(p))
	}

	// Number of detail records
	{
		p := make([]byte, 6)
		_, err := reader.Read(p)
		if err != nil {
			return err
		}
		fmt.Println(string(p))
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
	)
	reader := bufio.NewReader(fs)
	hdrBuf := make([]byte, 4)
	for {
		v, _ := reader.Read(hdrBuf)
		if v == 0 {
			break
		}

		switch string(hdrBuf) {
		case "1200":
			// file header
			isFixed, err = readFileHeader(reader)
			if err != nil {
				log.Fatalf("error while reading file header: %v", err)
			}

			log.Printf("isFixed: %v", isFixed)
		case "1201":
			// check index
			numImages, err := readCheckIndex(reader, isFixed)
			if err != nil {
				log.Fatalf("error while reading file header: %v", err)
			}

			log.Printf("number of images: %v", numImages)
		case "1202":
			// image header
			numImages, err := readImageHeader(reader, isFixed)
			if err != nil {
				log.Fatalf("error while reading file header: %v", err)
			}

			log.Printf("number of images: %v", numImages)
		case "1203":
			// image data
			err := readImageData(reader, isFixed)
			if err != nil {
				log.Fatalf("error while reading file header: %v", err)
			}
		case "1204":
			// check trailer
			err := readCheckTrailer(reader, isFixed)
			if err != nil {
				log.Fatalf("error while reading file header: %v", err)
			}
		case "1209":
			// file trailer
			err := readOutputTrailer(reader, isFixed)
			if err != nil {
				log.Fatalf("error while reading file header: %v", err)
			}
		}
	}
}
