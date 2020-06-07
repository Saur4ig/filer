package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

// info about big file piece
type chunk struct {
	// name of a big file
	baseFilename string
	// total number of chunks
	totalChunks int
	// current chunk number
	ID int
	// name template for temp chunk file
	chunkTemplate string
}

// get info of a file(chunk) and validate it
func getChunkAndValidate(req *http.Request) (*chunk, error) {
	baseFileName := req.Header.Get("Base-Filename")
	if baseFileName == "" || baseFileName == " " {
		return nil, fmt.Errorf("base filename missing")
	}

	chunkIDStr := req.Header.Get("Chunk-Id")
	if chunkIDStr == "" || chunkIDStr == " " {
		return nil, fmt.Errorf("chunk id missing")
	}

	totalStr := req.Header.Get("Total-Chunks")
	if totalStr == "" || totalStr == " " {
		return nil, fmt.Errorf("total number of chunks missing")
	}

	total, err := strconv.Atoi(totalStr)
	if err != nil {
		return nil, err
	}

	chunkID, err := strconv.Atoi(chunkIDStr)
	if err != nil {
		return nil, err
	}

	template := "chunk-*.txt"
	if total == chunkID {
		template = "*-" + baseFileName
	}

	return &chunk{
		baseFilename:  baseFileName,
		totalChunks:   total,
		ID:            chunkID,
		chunkTemplate: template,
	}, nil
}

// checks if all chunks downloaded
func (c *chunk) isAllDownloaded(dir string) (bool, int, error) {
	counter := 0
	data, err := ioutil.ReadDir(dir)
	if err != nil {
		return false, 0, err
	}
	for _, file := range data {
		if !file.IsDir() {
			counter++
		}
	}

	return counter == c.totalChunks, counter, nil
}

// merge chunks to a file and clean everything
func (c *chunk) mergeAndRemove(dir string, wg *sync.WaitGroup) {
	wg.Wait()
	err := c.merge(dir)
	if err != nil {
		log.Printf("merge file err - %s", err.Error())
		return
	}
	err = removeTempDir(dir)
	if err != nil {
		log.Printf("remove dir err - %s", err.Error())
		return
	}
}

// merge all chunks into one file
func (c *chunk) merge(dir string) error {
	newFileName, err := c.createNewFileName(filesDir)
	if err != nil {
		return err
	}

	newFile, err := os.OpenFile(newFileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	var position int64

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() {
			position, err = pasteChunkToFile(newFile, file, position, dir)
			if err != nil {
				return err
			}
		}
	}
	return newFile.Close()
}
