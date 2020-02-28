package main

import (
	"log"
	"os"
	"sync"
)

// fileInfo describes the file name and the data it holds
type fileInfo struct {
	name string
	data []byte
}

// writer listens on the file queue and create files
func writer(wg *sync.WaitGroup, q <-chan fileInfo) {
	for fInfo := range q {
		fInfo.writeImage()
		wg.Done()
	}
}

// writeImage writes image to the file
func (f *fileInfo) writeImage() {
	file, err := os.Create(f.name)
	if err != nil {
		log.Fatalf("error while creating image file: %v", err)
	}
	defer file.Close()

	if _, err := file.Write(f.data); err != nil {
		log.Fatalf("error while writing to image file: %v", err)
	}
}
