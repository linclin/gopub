package gokits

import (
	"os"
	"os/exec"
	"path/filepath"
)

// checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func FileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// returns true if given path is a file,
// or returns false when it's a directory or does not exist.
func IsFile(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

// returns true if given path is a directory,
// or returns false when it's a file or does not exist.
func IsDir(dir string) bool {
	f, e := os.Stat(dir)
	if e != nil {
		return false
	}
	return f.IsDir()
}

// returns file size in bytes and possible error.
func FileSize(file string) (int64, error) {
	f, err := os.Stat(file)
	if err != nil {
		return 0, err
	}
	return f.Size(), nil
}

// returns the complete directory of the current execution file
func GetProcPwd() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(filepath.Dir(file))
	return path
}
