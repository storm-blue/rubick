package common

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
)

func Download(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()
	return io.ReadAll(res.Body)
}

func Zip(files map[string]string) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	for key, resource := range files {
		fw, err := w.Create(key)
		if err != nil {
			return nil, err
		}
		_, err = fw.Write([]byte(resource))
		if err != nil {
			return nil, err
		}
	}

	if err := w.Flush(); err != nil {
		return nil, err
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Unzip(data []byte) (map[string]string, error) {
	// Open a zip archive for reading.
	reader := bytes.NewReader(data)

	r, err := zip.NewReader(reader, int64(len(data)))
	if err != nil {
		return nil, err
	}

	result := map[string]string{}

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		fileBytes, err := io.ReadAll(rc)
		if err != nil {
			return nil, err
		}
		result[f.Name] = string(fileBytes)

		_ = rc.Close()
	}

	return result, nil
}
