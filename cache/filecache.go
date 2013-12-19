package cache

import (
	// "strings"
	"io/ioutil"
	"path/filepath"
	"os"
)

func NewFilecacheEngine(path string) *FilecacheEngine {
	app := filepath.Dir(os.Args[0])
	cachepath := filepath.Join(app, path)
	ok, err := exists(cachepath)
	if err != nil {}
	if !ok {
		if err = os.Mkdir(cachepath, os.ModePerm); err != nil {}
	}

	return &FilecacheEngine{
		Path: cachepath,
	}
}

type FilecacheEngine struct {
	Path string
}

func (fc *FilecacheEngine) Get(key string) (*Item, error) {
	filename := filepath.Join(fc.Path, key)
	value, err := ioutil.ReadFile(filename)

	return &Item{
		Key: key,
		Value: value,
	}, err
}

func (fc *FilecacheEngine) Set(key string, value []byte) (err error) {
	filename := filepath.Join(fc.Path, key)
	err = ioutil.WriteFile(filename, value, os.ModePerm)
	return err
}


func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
            return true, nil
    }
    if os.IsNotExist(err) {
            return false, nil
    }
    return false, err
}
