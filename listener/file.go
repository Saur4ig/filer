package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// insert chunk to the main file
func pasteChunkToFile(newFile *os.File, chunkFile os.FileInfo, position int64, dir string) (int64, error) {
	chunk, err := os.Open(filepath.Join(dir, chunkFile.Name()))
	if err != nil {
		return 0, err
	}
	defer chunk.Close()

	size := chunkFile.Size()
	buffer := make([]byte, size)
	position = position + size

	reader := bufio.NewReader(chunk)
	_, err = reader.Read(buffer)
	if err != nil {
		return 0, err
	}
	_, err = newFile.Write(buffer)
	if err != nil {
		return 0, err
	}

	err = newFile.Sync()
	if err != nil {
		return 0, err
	}

	// nolint
	buffer = nil

	return position, nil
}

// create name and named file
func (c *chunk) createNewFileName(dir string) (string, error) {
	fileName := filepath.Join(dir, "new_"+c.baseFilename)
	_, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

// remove temporary directory, the home of chunks
func removeTempDir(dir string) error {
	return os.RemoveAll(dir)
}

// save all strings to file
func saveToFile(path string, values []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, value := range values {
		_, err = fmt.Fprintln(f, value)
		if err != nil {
			return err
		}
	}
	return nil
}
