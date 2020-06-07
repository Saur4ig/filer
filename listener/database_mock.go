package main

import (
	"sync"
)

type Database interface {
	isUnique(str string) bool
	add(str string)
	all() []string
}

type cell struct {
	sync.Mutex
	data map[string]struct{}
}

type storage struct {
	sync.Mutex
	bucket map[string]*cell
}

type clickhouse struct {
	storage *storage
}

// check if string is unique in database
func (c *clickhouse) isUnique(str string) bool {
	key := getFirst(str)
	c.storage.Lock()
	defer c.storage.Unlock()
	if _, ok := c.storage.bucket[key]; ok {
		c.storage.bucket[key].Lock()
		defer c.storage.bucket[key].Unlock()
		if _, found := c.storage.bucket[key].data[str]; found {
			return false
		}
	}
	return true
}

// add a string to database
// if first letter of a string(key) is unique -> create a new map(cell)
func (c *clickhouse) add(str string) {
	key := getFirst(str)
	c.storage.Lock()
	defer c.storage.Unlock()
	if _, ok := c.storage.bucket[key]; !ok {
		c.storage.bucket[key] = &cell{
			data: make(map[string]struct{}),
		}
	}
	c.storage.bucket[key].Lock()
	c.storage.bucket[key].data[str] = struct{}{}
	c.storage.bucket[key].Unlock()
}

// get all unique from db
// ideally it should be something like a stream, but for test task - ok
func (c *clickhouse) all() []string {
	res := make([]string, 0)
	for _, cell := range c.storage.bucket {
		for key := range cell.data {
			res = append(res, key)
		}
	}
	return res
}

func newDatabase() Database {
	return &clickhouse{
		storage: &storage{
			bucket: make(map[string]*cell),
		},
	}
}

// returns first letter of string
func getFirst(str string) string {
	if str == "" {
		return str
	}
	return string([]rune(str)[0])
}
