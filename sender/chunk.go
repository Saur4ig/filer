package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

// chunk of file with needed data inside
type Chunk struct {
	ChunkID      int
	TotalChunks  int
	BaseFileName string
	ChunkName    string
}

// sends chunk/file to listener
func (c *Chunk) send(client *http.Client, listenerURL string) error {
	r, w := io.Pipe()

	req, err := http.NewRequest(http.MethodPut, listenerURL, r)
	if err != nil {
		return err
	}

	mpw := multipart.NewWriter(w)

	go func() {
		var part io.Writer
		defer w.Close()

		if part, err = mpw.CreateFormFile("file", c.ChunkName); err != nil {
			log.Fatal(err)
		}
		chunk, err := os.Open(c.ChunkName)
		if err != nil {
			log.Fatal(err)
		}
		defer chunk.Close()

		part = io.MultiWriter(part)
		if _, err = io.Copy(part, chunk); err != nil {
			log.Fatal(err)
		}
		if err = mpw.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	req.Header.Set("Content-Type", mpw.FormDataContentType())
	req.Header.Set("Chunk-Id", fmt.Sprintf("%d", c.ChunkID))
	req.Header.Set("Total-Chunks", fmt.Sprintf("%d", c.TotalChunks))
	req.Header.Set("Base-Filename", c.BaseFileName)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Printf("sent chunk #%d, for file %s, total number - %d\n", c.ChunkID, c.BaseFileName, c.TotalChunks)

	return nil
}
