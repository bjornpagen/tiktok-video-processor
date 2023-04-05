package storer

import (
	"io"
	"os"
	"path/filepath"
)

type Storer interface {
	Store(path string) (access string, err error)
	Get(access string) (tmpfile string, err error)
}

type LocalStorer struct {
	Path string
}

func NewLocalStorer(path string) *LocalStorer {
	return &LocalStorer{
		Path: path,
	}
}

func (ls *LocalStorer) Store(file string) (string, error) {
	// Make sure path exists
	if _, err := os.Stat(ls.Path); os.IsNotExist(err) {
		err = os.MkdirAll(ls.Path, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	// Copy the input file to the destination directory
	dstPath := filepath.Join(ls.Path, filepath.Base(file))
	err := copyFile(file, dstPath)
	if err != nil {
		return "", err
	}
	return dstPath, nil
}

func (ls *LocalStorer) Get(access string) (string, error) {
	// As it's a local storer, we just return the same access string as the input file
	return access, nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	err = destinationFile.Sync()
	if err != nil {
		return err
	}

	sourceFileStat, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	err = os.Chmod(dst, sourceFileStat.Mode())
	if err != nil {
		return err
	}

	return nil
}
