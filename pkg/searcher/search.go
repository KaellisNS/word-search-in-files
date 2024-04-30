package searcher

import (
	"bufio"
	"os"
	"strings"
	"sync"
)

type Index struct {
	data  map[string][]string
	mutex sync.Mutex
}

func NewIndex() *Index {
	return &Index{data: make(map[string][]string)}
}

func (index *Index) IndexFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		for _, word := range words {
			index.mutex.Lock()
			index.data[word] = append(index.data[word], filename)
			index.mutex.Unlock()
		}
	}

	return scanner.Err()
}

func (index *Index) Search(keyword string) ([]string, bool) {
	index.mutex.Lock()
	defer index.mutex.Unlock()

	files, found := index.data[keyword]
	return files, found
}
