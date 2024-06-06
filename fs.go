package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

// fileExistsWith check whether file exists with given content.
func fileExistsWith(path string, content string) (bool, error) {
	// Attempt to open file for just reads.
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("error opening file: %w", err)
	}

	if file == nil {
		// file not found.
		return false, nil
	}

	// Ensure closed.
	defer file.Close()

	// Read file into memory.
	b, err := io.ReadAll(file)
	if err != nil {
		return false, fmt.Errorf("error reading file: %w", err)
	}

	// Check file content is a match.
	return string(b) == content, nil
}

// mkdirAll creates all necessary parent directories in path, with given perms.
func mkdirAll(path string, perm fs.FileMode) error {
	if dryrun {
		// Do nothing.
		return nil
	}

	// Create all directories in path.
	return os.MkdirAll(path, perm)
}

// writeJSON writes given in-memory structure as JSON bytes to file located at $path.
func writeJSON(path string, data any) (int, error) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return 0, fmt.Errorf("error marshaling json: %w", err)
	}
	return writeBytes(path, b)
}

// writeBytes writes given bytes to file located at $path.
func writeBytes(path string, data []byte) (int, error) {
	n, err := write(path, io.NopCloser(bytes.NewReader(data)))
	return int(n), err
}

// writeString writes given string to file located at $path.
func writeString(path string, data string) (int, error) {
	n, err := write(path, io.NopCloser(strings.NewReader(data)))
	return int(n), err
}

// write writes given read closer stream to file located at $path.
func write(path string, data io.ReadCloser) (int64, error) {
	if dryrun {
		// Do nothing.
		return 0, nil
	}

	// Attempt to open file at path on disk, creating if new.
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return 0, fmt.Errorf("error opening file %s: %w", path, err)
	}

	defer func() {
		// Ensure closed.
		_ = file.Close()
		_ = data.Close()

		if err != nil {
			// On error, remove the
			// failed-to-write file.
			_ = os.Remove(path)
		}
	}()

	// Stream data to file handle.
	n, err := file.ReadFrom(data)
	if err != nil {
		return n, fmt.Errorf("error writing data to %s: %w", path, err)
	}

	return n, nil
}
