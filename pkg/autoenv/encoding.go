package autoenv

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Marshal encodes the object into the gzenv format
func Marshal(obj interface{}) (string, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	zlibData := bytes.NewBuffer([]byte{})
	w := zlib.NewWriter(zlibData)
	w.Write(jsonData)
	w.Close()

	base64Data := base64.URLEncoding.EncodeToString(zlibData.Bytes())
	return base64Data, nil
}

// Unmarshal restores the gzenv format back into a Go object
func Unmarshal(gzenv string, obj interface{}) error {
	gzenv = strings.TrimSpace(gzenv)

	data, err := base64.URLEncoding.DecodeString(gzenv)
	if err != nil {
		return fmt.Errorf("unmarshal() base64 decoding: %v", err)
	}

	zlibReader := bytes.NewReader(data)
	w, err := zlib.NewReader(zlibReader)
	if err != nil {
		return fmt.Errorf("unmarshal() zlib opening: %v", err)
	}

	envData := bytes.NewBuffer([]byte{})
	_, err = io.Copy(envData, w)
	if err != nil {
		return fmt.Errorf("unmarshal() zlib decoding: %v", err)
	}
	w.Close()

	err = json.Unmarshal(envData.Bytes(), &obj)
	if err != nil {
		return fmt.Errorf("unmarshal() json parsing: %v", err)
	}

	return nil
}
