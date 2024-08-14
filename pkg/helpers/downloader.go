package helpers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type downloader struct {
	url string
}

func NewDownloader(url string) *downloader {
	return &downloader{url}
}

// DownloadToFile downloads to a temporary file then rename to the specified location.
func (d *downloader) DownloadToFile(file string) error {
	resp, err := http.Get(d.url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to download (code %d)", resp.StatusCode)
	}
	defer func() {
		cerr := resp.Body.Close()
		if err == nil {
			err = cerr
		}
	}()

	// Create a temp file
	tmpFile, err := os.CreateTemp(filepath.Dir(file), "")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	// Download with progress information
	reader := &progressPrinterReader{Reader: resp.Body, Size: resp.ContentLength}
	_, err = io.Copy(tmpFile, reader)
	if err != nil {
		return err
	}

	// Move temp file to final destination
	return os.Rename(tmpFile.Name(), file)
}

type progressPrinterReader struct {
	io.Reader
	Size int64

	total int64
}

func (pp *progressPrinterReader) display(msg string) {
	fmt.Printf("\r%s\r%s", strings.Repeat(" ", 25), msg)
}

func (pp *progressPrinterReader) progress(n int, err error) {
	if err == io.EOF {
		pp.display("Download complete\n")
		return
	}

	pp.total += int64(n)
	pct := float64(pp.total) / float64(pp.Size) * 100
	pp.display(fmt.Sprintf("Downloading... %.0f%%", pct))

}

func (pp *progressPrinterReader) Read(p []byte) (int, error) {
	n, err := pp.Reader.Read(p)
	pp.progress(n, err)
	return n, err
}
