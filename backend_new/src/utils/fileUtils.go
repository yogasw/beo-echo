package utils

import (
	"io"
	"os"
	"path/filepath"
)

// EnsureDirectoryExists creates a directory if it doesn't exist
func EnsureDirectoryExists(dirPath string) error {
	return os.MkdirAll(dirPath, os.ModePerm)
}

// FileExists checks if a file exists at the given path
func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// CopyFile copies a file from source to destination
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Make sure the destination directory exists
	if err := EnsureDirectoryExists(filepath.Dir(dst)); err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// DeleteFile deletes the file at the given path
func DeleteFile(path string) error {
	return os.Remove(path)
}
