package helpers

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

func close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func makeTemporaryFile() (f *os.File, err error) {
	tmpFile, err := ioutil.TempFile("", dadPrefix)

	if err != nil {
		return nil, err
	}

	return tmpFile, nil
}
