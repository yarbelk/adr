package ioutils

import (
	"os"
)

// ReadDir reads the directory named by dirname and returns
// a list of directory entries sorted by filename.
// Returns an unsorted list of files
func ReadDir(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	return list, nil
}
