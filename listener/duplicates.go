package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

// creates a file with all unique strings
func createUniqueFile(
	filename, dir string,
	totalNumber int,
	wait *sync.WaitGroup,
	errors chan error,
) {
	defer wait.Done()
	data, err := ioutil.ReadDir(dir)
	if err != nil {
		errors <- err
		return
	}

	database := newDatabase()

	var wg sync.WaitGroup
	wg.Add(totalNumber)
	for _, file := range data {
		if !file.IsDir() {
			go removeDuplicates(database, &wg, dir+"/"+file.Name(), errors)
		}
	}
	wg.Wait()

	err = saveToFile(filesUniqueDir+"/"+filename, database.all())
	if err != nil {
		errors <- err
	}
}

// removes all duplicates from a single file
func removeDuplicates(
	db Database,
	wg *sync.WaitGroup,
	path string,
	errors chan error,
) {
	defer wg.Done()

	file, err := os.Open(path)
	if err != nil {
		errors <- err
		return
	}
	defer file.Close()

	fileStats, err := file.Stat()
	if err != nil {
		errors <- err
		return
	}

	if fileStats.Size() > maxFileChunkSize {
		errors <- fmt.Errorf("file too big %d", fileStats.Size())
		return
	}

	// scan by line from file
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	// insert unique data to improvised database
	insertUnique(fileScanner, db)
}

func insertUnique(scanner *bufio.Scanner, db Database) {
	temp := make(map[string]struct{})

	for scanner.Scan() {
		text := scanner.Text()
		if _, ok := temp[text]; !ok {
			if db.isUnique(text) {
				db.add(text)
			}
			temp[text] = struct{}{}
		}
	}
}
