package main

import (
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
)

// file saving
func ProcessHandler(w http.ResponseWriter, req *http.Request) {
	// get file info and validate it
	c, err := getChunkAndValidate(req)
	if err != nil {
		ErrorResponse(w, err)
		log.Printf("validation err - %s", err.Error())
		return
	}

	baseFileName := c.baseFilename

	log.Printf("Chunk #%d, for file %s recieved, total - %d", c.ID, baseFileName, c.totalChunks)

	reader, err := req.MultipartReader()
	if err != nil {
		ErrorResponse(w, err)
		log.Printf("failed start reader, err - %s", err.Error())
		return
	}

	// setup a dir, if it is single file - use files dir directly, if it is a chunk - use/create temp dir
	dir, err := getDir(c)
	if err != nil {
		ErrorResponse(w, err)
		log.Printf("get file dir err - %s", err.Error())
		return
	}

	// get all part of multipart file
	err = copyPart(reader, c, dir)
	if err != nil {
		ErrorResponse(w, err)
		log.Printf("copy part err - %s", err.Error())
		return
	}

	// send ok response
	OKResponse(w)

	downloaded, counter, err := c.isAllDownloaded(dir)
	if err != nil {
		log.Printf("check downloaded err - %s", err.Error())
		return
	}

	// if all chunks downloaded - merge them and create file with all uniques
	if downloaded {
		go createUniqueAndMerge(baseFileName, dir, counter, c)
	}
}

func createUniqueAndMerge(baseFileName, dir string, counter int, c *chunk) {
	var wg sync.WaitGroup
	wg.Add(1)
	errors := make(chan error)
	go createUniqueFile(baseFileName, dir, counter, &wg, errors)
	go c.mergeAndRemove(dir, &wg)
	select {
	case err := <-errors:
		if err != nil {
			log.Printf("err - %s", err.Error())
		}
	default:
		return
	}
}

// copy all parts to a single file
func copyPart(
	reader *multipart.Reader,
	chunk *chunk,
	dir string,
) error {
	var temp *os.File
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		temp, err = ioutil.TempFile(dir, chunk.chunkTemplate)
		if err != nil {
			return err
		}

		_, err = io.Copy(temp, part)
		if err != nil {
			return err
		}
	}
	return temp.Close()
}

// get dir for current chunk
func getDir(chunk *chunk) (string, error) {
	// for single file - save directly to files folder
	if chunk.ID == chunk.totalChunks {
		return filesDir, nil
	}

	dir := chunk.baseFilename + "_dir"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0777)
		if err != nil {
			return "", err
		}
	}
	return dir, nil
}
