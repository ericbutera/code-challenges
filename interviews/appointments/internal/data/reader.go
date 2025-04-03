package data

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"
	"runtime"
)

var ErrInvalidPath = errors.New("invalid path")

func GetPath(datafile string) (string, error) {
	// https://hackandsla.sh/posts/2020-11-23-golang-test-fixtures/
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", ErrInvalidPath
	}
	return path.Join(path.Dir(filename), "..", "..", "testdata", datafile), nil
}

func FileToReadCloser(file string) (io.ReadCloser, error) {
	path, err := GetPath(file)
	if err != nil {
		return nil, err
	}
	return os.Open(path)
}

func GetFileJSONAs[T any](file string) (T, error) {
	var data T

	reader, err := FileToReadCloser(file)
	if err != nil {
		return data, err
	}
	defer reader.Close()

	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err != nil {
		return data, err
	}

	return data, nil
}
