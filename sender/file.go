package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"time"
)

const (
	// 5mb
	fileChunkSize = 5 * (1 << 20)
)

// sends file to a listener
func sender(fileToSend, listenerURL string) error {
	file, err := os.Open(fileToSend)
	if err != nil {
		return err
	}
	defer file.Close()

	fileStats, err := file.Stat()
	if err != nil {
		return err
	}

	client := &http.Client{}

	// if file more than 5mb -> chunk it to several files, and send by parts
	if fileStats.Size() >= fileChunkSize {
		return chunkAndSend(file, fileStats, client, listenerURL)
	} else {
		c := Chunk{
			ChunkID:      1,
			TotalChunks:  1,
			BaseFileName: file.Name(),
			ChunkName:    file.Name(),
		}
		err := c.send(client, listenerURL)
		if err != nil {
			return err
		}
	}

	return nil
}

// splits file to a parts, sends every part to listener and clean-ups after sending
func chunkAndSend(file *os.File, fileStats os.FileInfo, client *http.Client, listenerURL string) error {
	// timestamp for temp dir
	t := time.Now()

	// ideally generate unique id for every folder, coz of same name and type of file
	path := "send_" + file.Name() + "_dir" + "_" + t.Format("20060102150405")

	filesToSend, err := SplitFile(file, fileStats, path)
	if err != nil {
		return err
	}

	for _, c := range filesToSend {
		err = c.send(client, listenerURL)
		if err != nil {
			return err
		}
	}
	return removeTempDir(path)
}

// splits one big file into a bunch of small one
// creates a dir for all chunks of this file
func SplitFile(file *os.File, fileStats os.FileInfo, path string) ([]Chunk, error) {
	if os.MkdirAll(path, 0777) != nil {
		return nil, fmt.Errorf("failed to create folder %s", path)
	}

	// get total number of chunks
	total := int(math.Ceil(float64(fileStats.Size()) / float64(fileChunkSize)))
	filesToSend := make([]Chunk, total)

	for i := 0; i < total; i++ {
		c := Chunk{
			ChunkID:      i,
			TotalChunks:  total,
			BaseFileName: file.Name(),
		}
		chunkSize := int(math.Min(fileChunkSize, float64(fileStats.Size()-int64(i*fileChunkSize))))

		// chunk part buffer
		buf := make([]byte, chunkSize)
		_, err := file.Read(buf)
		if err != nil {
			return nil, err
		}

		// creates a file
		chunkName := fmt.Sprintf("%s/chunk_%s_%d", path, file.Name(), i)
		_, err = os.Create(chunkName)
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile(chunkName, buf, os.ModePerm)
		if err != nil {
			return nil, err
		}
		c.ChunkName = chunkName
		filesToSend[i] = c
		fmt.Println("Chunked to : ", chunkName)
	}

	return filesToSend, nil
}

// removes temp dir with all chunks inside
func removeTempDir(dir string) error {
	return os.RemoveAll(dir)
}
