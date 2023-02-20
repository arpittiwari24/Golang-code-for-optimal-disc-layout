package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

type IndexEntry struct {
	key       int64
	valuePos  int64
	valueSize int32
}

type Record struct {
	key   int64
	value string
}

func main() {
	// create a new B+ tree index
	index := make([]IndexEntry, 0)

	// open file handles for the index data and record data
	indexFile, err := os.Create("index.txt")
	if err != nil {
		panic(err)
	}
	recordFile, err := os.Create("record.txt")
	if err != nil {
		panic(err)
	}

	// add some sample records to the index
	index = append(index, IndexEntry{1, 0, int32(len("record 1"))})
	index = append(index, IndexEntry{2, int64(len("record 1")), int32(len("record 2"))})
	index = append(index, IndexEntry{3, int64(len("record 1") + len("record 2")), int32(len("record 3"))})

	// write the index data to the index file
	for _, entry := range index {
		err := binary.Write(indexFile, binary.LittleEndian, &entry)
		if err != nil {
			panic(err)
		}
	}

	// write the record data to the record file
	for _, entry := range index {
		record := Record{entry.key, fmt.Sprintf("record %d", entry.key)}
		_, err := fmt.Fprintln(recordFile, record.value)
		if err != nil {
			panic(err)
		}
	}

	// perform a lookup operation on the index
	key := int64(2)
	var value string
	for _, entry := range index {
		if entry.key == key {
			value = readRecord(recordFile, entry.valuePos, int64(entry.valueSize))
			break
		}
	}
	if value != "" {
		fmt.Println("Found record:", value)
	} else {
		fmt.Println("Record not found")
	}

	// close the file handles
	err = indexFile.Close()
	if err != nil {
		panic(err)
	}
	err = recordFile.Close()
	if err != nil {
		panic(err)
	}
}

func readRecord(file *os.File, pos int64, size int64) string {
	buffer := make([]byte, size)
	_, err := file.ReadAt(buffer, pos)
	if err != nil {
		panic(err)
	}
	return string(buffer)
}
