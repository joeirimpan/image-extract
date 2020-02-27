package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

// Various header types
const (
	TypeFileHeader   = "1200"
	TypeCheckIndex   = "1201"
	TypeImageHeader  = "1202"
	TypeImageData    = "1203"
	TypeCheckTrailer = "1204"
	TypeFileTrailer  = "1209"
)

type disbursement struct {
	isFixed bool

	reader io.Reader
}

func (d *disbursement) parseFileHeader() error {
	// File ID
	{
		p := make([]byte, 15)
		_, err := d.reader.Read(p)
		if err != nil {
			return err
		}
	}

	// Request ID
	{
		p := make([]byte, 15)
		_, err := d.reader.Read(p)
		if err != nil {
			return err
		}
	}

	// File Version
	{
		p := make([]byte, 4)
		_, err := d.reader.Read(p)
		if err != nil {
			return err
		}
	}

	// File creation date
	{
		p := make([]byte, 8)
		_, err := d.reader.Read(p)
		if err != nil {
			return err
		}
	}

	// File creation time
	{
		p := make([]byte, 6)
		_, err := d.reader.Read(p)
		if err != nil {
			return err
		}
	}

	// Number of check records
	{
		p := make([]byte, 6)
		_, err := d.reader.Read(p)
		if err != nil {
			return err
		}
	}

	// Record size
	{
		p := make([]byte, 4)
		_, err := d.reader.Read(p)
		if err != nil {
			return err
		}

		// Read filler
		switch string(p) {
		// Fixed
		case "0256":
			d.isFixed = true

			filler := make([]byte, 194)
			_, err := d.reader.Read(filler)
			if err != nil {
				return err
			}
		// Variable
		case "0090":
			filler := make([]byte, 28)
			_, err := d.reader.Read(filler)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *disbursement) parseCheckIndex() (int, string, error) {
	var (
		numImages int
		checkNo   string
	)

	// Bank number
	{
		p := make([]byte, 4)
		_, err := d.reader.Read(p)
		if err != nil {
			return numImages, checkNo, err
		}
	}

	// Routing transit number
	{
		p := make([]byte, 9)
		_, err := d.reader.Read(p)
		if err != nil {
			return numImages, checkNo, err
		}
	}

	// Account number
	{
		p := make([]byte, 20)
		_, err := d.reader.Read(p)
		if err != nil {
			return numImages, checkNo, err
		}
	}

	// Check number
	{
		p := make([]byte, 15)
		_, err := d.reader.Read(p)
		if err != nil {
			return numImages, checkNo, err
		}

		checkNo = string(p)
	}

	// Amount
	{
		p := make([]byte, 10)
		_, err := d.reader.Read(p)
		if err != nil {
			return numImages, checkNo, err
		}

	}

	// Seq number
	{
		p := make([]byte, 15)
		_, err := d.reader.Read(p)
		if err != nil {
			return numImages, checkNo, err
		}
	}

	// Posted date
	{
		p := make([]byte, 8)
		_, err := d.reader.Read(p)
		if err != nil {
			return numImages, checkNo, err
		}
	}

	// Number of images
	{
		p := make([]byte, 4)
		_, err := d.reader.Read(p)
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
		if d.isFixed {
			fillerLength = 167
		}

		p := make([]byte, fillerLength)
		_, err := d.reader.Read(p)
		if err != nil {
			return numImages, checkNo, err
		}
	}

	return numImages, checkNo, nil
}

func (d *disbursement) parseImageHeader() (int, string, error) {
	var (
		numImages int
		imgType   string
	)
	// Image type
	{
		p := make([]byte, 4)
		_, err := d.reader.Read(p)
		if err != nil {
			return numImages, imgType, err
		}

		imgType = string(p)
	}

	// Image side
	{
		p := make([]byte, 1)
		_, err := d.reader.Read(p)
		if err != nil {
			return numImages, imgType, err
		}
	}

	// Number of records
	{
		p := make([]byte, 4)
		_, err := d.reader.Read(p)
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
		_, err := d.reader.Read(p)
		if err != nil {
			return numImages, imgType, err
		}
	}

	// Filler
	{
		var fillerLength = 71
		if d.isFixed {
			fillerLength = 231
		}

		p := make([]byte, fillerLength)
		_, err := d.reader.Read(p)
		if err != nil {
			return numImages, imgType, err
		}
	}

	return numImages, imgType, nil
}

func (d *disbursement) parseImageData() ([]byte, error) {
	var (
		recLength int
		imageData []byte
	)
	// Record length
	{
		p := make([]byte, 4)
		_, err := d.reader.Read(p)
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
		if d.isFixed {
			recLength = 246
		}

		imageData = make([]byte, recLength)
		_, err := d.reader.Read(imageData)
		if err != nil {
			return imageData, err
		}
	}

	return imageData, nil
}

func (d *disbursement) parseCheckTrailer(isFixed bool) error {
	// Filler
	var fillerLength = 86
	if isFixed {
		fillerLength = 252
	}

	p := make([]byte, fillerLength)
	_, err := d.reader.Read(p)
	if err != nil {
		return err
	}

	return nil
}

func (d *disbursement) parseFileTrailer(isFixed bool) error {
	// File ID
	{
		p := make([]byte, 15)
		_, err := d.reader.Read(p)
		if err != nil {
			return err
		}
	}

	// Request ID
	{
		p := make([]byte, 15)
		_, err := d.reader.Read(p)
		if err != nil {
			return err
		}
	}

	// File version
	{
		p := make([]byte, 4)
		_, err := d.reader.Read(p)
		if err != nil {
			return err
		}
	}

	// File creation date
	{
		p := make([]byte, 8)
		_, err := d.reader.Read(p)
		if err != nil {
			return err
		}
	}

	// File creation time
	{
		p := make([]byte, 6)
		_, err := d.reader.Read(p)
		if err != nil {
			return err
		}
	}

	// Number of detail records
	{
		p := make([]byte, 6)
		_, err := d.reader.Read(p)
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
		_, err := d.reader.Read(p)
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
		imgCt   int
		isFixed bool
		checkNo string
		imgType string

		reader = bufio.NewReader(fs)
		disb   = &disbursement{
			reader: reader,
		}
		hdrBuf = make([]byte, 4)
	)

	for {
		_, err := reader.Read(hdrBuf)
		if err == io.EOF {
			break
		}

		switch string(hdrBuf) {
		case TypeFileHeader:
			err = disb.parseFileHeader()
			if err != nil {
				log.Fatalf("error while reading file header: %v", err)
			}

		case TypeCheckIndex:
			_, checkNo, err = disb.parseCheckIndex()
			if err != nil {
				log.Fatalf("error while reading check index: %v", err)
			}

		case TypeImageHeader:
			_, imgType, err = disb.parseImageHeader()
			if err != nil {
				log.Fatalf("error while reading image header: %v", err)
			}

		case TypeImageData:
			imgCt++
			image, err := disb.parseImageData()
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

		case TypeCheckTrailer:
			err := disb.parseCheckTrailer(isFixed)
			if err != nil {
				log.Fatalf("error while reading check trailer: %v", err)
			}

		case TypeFileTrailer:
			err := disb.parseFileTrailer(isFixed)
			if err != nil {
				log.Fatalf("error while reading file trailer: %v", err)
			}
		}
	}
}
