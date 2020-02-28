package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
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

	reader *bufio.Reader
}

func (d *disbursement) parseFileHeader() error {
	var err error
	// File ID(15)
	// Request ID(15)
	// File Version(4)
	// File creation date(8)
	// File creation time(6)
	// Number of check records(6)
	// 54
	if _, err = d.reader.Discard(54); err != nil {
		return err
	}

	// Record size
	p := make([]byte, 4)
	if _, err = d.reader.Read(p); err != nil {
		return err
	}

	// Read filler
	var skipBytes int
	switch string(p) {
	// Fixed
	case "0256":
		d.isFixed = true
		skipBytes = 194
	// Variable
	case "0090":
		skipBytes = 28
	}

	if _, err := d.reader.Discard(skipBytes); err != nil {
		return err
	}

	return nil
}

func (d *disbursement) parseCheckIndex() (string, error) {
	var checkNo string

	// Bank number(4)
	// Routing transit number(9)
	// Account number(20)
	// 33
	if _, err := d.reader.Discard(33); err != nil {
		return checkNo, err
	}

	// Check number
	p := make([]byte, 15)
	if _, err := d.reader.Read(p); err != nil {
		return checkNo, err
	}

	checkNo = string(p)

	// Amount(10)
	// Seq number(15)
	// Posted date(8)
	// Number of images(4)
	// 37
	var skipBytes = 37
	// Filler
	if d.isFixed {
		skipBytes += 167
	} else {
		skipBytes++
	}

	if _, err := d.reader.Discard(skipBytes); err != nil {
		return checkNo, err
	}

	return checkNo, nil
}

func (d *disbursement) parseImageHeader() (string, error) {
	var imgType string
	// Image type
	p := make([]byte, 4)
	if _, err := d.reader.Read(p); err != nil {
		return imgType, err
	}

	imgType = string(p)

	// Image side(1)
	// Number of records(4)
	// Image data record length(6)
	// 11
	var skipBytes = 11

	// Filler
	if d.isFixed {
		skipBytes += 231
	} else {
		skipBytes += 71
	}
	if _, err := d.reader.Discard(skipBytes); err != nil {
		return imgType, err
	}

	return imgType, nil
}

func (d *disbursement) parseImageData() ([]byte, error) {
	var (
		recLength int
		imageData []byte

		err error
	)
	// Record length
	p := make([]byte, 4)
	if _, err := d.reader.Read(p); err != nil {
		return imageData, err
	}

	recLength, err = strconv.Atoi(string(p))
	if err != nil {
		return imageData, err
	}

	// Image data
	if d.isFixed {
		recLength = 246
	}

	imageData = make([]byte, recLength)
	if _, err := d.reader.Read(imageData); err != nil {
		return imageData, err
	}

	return imageData, nil
}

func (d *disbursement) parseCheckTrailer(isFixed bool) error {
	// Filler
	var fillerLength = 86
	if isFixed {
		fillerLength = 252
	}

	if _, err := d.reader.Discard(fillerLength); err != nil {
		return err
	}

	return nil
}

func (d *disbursement) parseFileTrailer(isFixed bool) error {
	// File ID(15)
	// Request ID(15)
	// File version(4)
	// File creation date(8)
	// File creation time(6)
	// Number of detail records(6)
	// 54
	var skipBytes = 54
	// Filler
	if isFixed {
		skipBytes += 198
	} else {
		skipBytes += 32
	}

	if _, err := d.reader.Discard(skipBytes); err != nil {
		return err
	}

	return nil
}

// Run runs the image extracter program
func Run(fileName string) {
	fs, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("error reading dat file")
	}
	defer fs.Close()

	// Setup file writer queue
	wg := &sync.WaitGroup{}
	fQueue := make(chan fileInfo)
	go writer(wg, fQueue)

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
			checkNo, err = disb.parseCheckIndex()
			if err != nil {
				log.Fatalf("error while reading check index: %v", err)
			}

		case TypeImageHeader:
			imgType, err = disb.parseImageHeader()
			if err != nil {
				log.Fatalf("error while reading image header: %v", err)
			}

		case TypeImageData:
			imgCt++
			image, err := disb.parseImageData()
			if err != nil {
				log.Fatalf("error while reading image data: %v", err)
			}

			// Push to write queue
			wg.Add(1)
			fQueue <- fileInfo{
				name: fmt.Sprintf("%s.%s", checkNo, imgType),
				data: image,
			}

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

	wg.Wait()
}

// TODO: Make sure that we read the correct number of records as specified in the headers.
func main() {
	if len(os.Args) < 2 {
		log.Fatalf("missing required arguments")
	}

	Run(os.Args[1])
}
