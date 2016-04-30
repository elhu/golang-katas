package anagram

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
)

const dictPath string = "/usr/share/dict/words"

type byteArray []byte

func (b byteArray) Len() int {
	return len(b)
}

func (b byteArray) Less(i, j int) bool {
	return b[i] < b[j]
}

func (b byteArray) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

type wordsLength [][]byte

func (w wordsLength) Len() int {
	return len(w)
}

func (w wordsLength) Less(i, j int) bool {
	return len(w[i]) < len(w[j])
}

func (w wordsLength) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

type AnagramFinder struct {
	data  map[string][][]byte
	words [][]byte
	mux   sync.Mutex
}

func New() *AnagramFinder {
	finder := &AnagramFinder{
		data:  make(map[string][][]byte),
		words: make([][]byte, 1),
	}

	fh, err := os.Open(dictPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fh.Close()

	finder.loadDict(fh)

	sliceCount := countSlices(finder.words)
	jobs := make(chan [][]byte)
	res := make(chan map[string][][]byte)

	for i := 0; i < sliceCount; i++ {
		go finder.buildAnagramsWorker(jobs, res)
	}

	slice, prev := nextSlice(finder.words, 0)
	for i := 0; i < sliceCount; i++ {
		jobs <- slice
		slice, prev = nextSlice(finder.words, prev)
	}
	for i := 0; i < sliceCount; i++ {
		data := <-res
		finder.mux.Lock()
		for key, value := range data {
			finder.data[key] = value
		}
		finder.mux.Unlock()
	}
	return finder
}

func prepWord(word []byte) string {
	sortedBytes := bytes.ToLower(word)
	sort.Sort(byteArray(sortedBytes))
	return strings.Trim(string(sortedBytes[:]), "\n")
}

func (a *AnagramFinder) AnagramsFor(word []byte) [][]byte {
	return a.data[prepWord(word)]
}

func (a *AnagramFinder) buildAnagramsWorker(jobs chan [][]byte, res chan map[string][][]byte) {
	data := make(map[string][][]byte)
	for words := range jobs {
		for _, word := range words {
			sortedBytesAsString := prepWord(word)
			_, ok := data[sortedBytesAsString]
			if !ok {
				data[sortedBytesAsString] = make([][]byte, 0, 1)
			}
			data[sortedBytesAsString] = append(data[sortedBytesAsString], word)
		}
		res <- data
	}
}

func nextSlice(words [][]byte, prev int) ([][]byte, int) {
	start_len := len(words[prev])
	i := prev
	for i < len(words) && len(words[i]) == start_len {
		i++
	}
	return words[prev+1 : i], i
}

func countSlices(words [][]byte) int {
	i := 0
	currLen := 0
	for _, word := range words {
		if len(word) != currLen {
			i++
		}
		currLen = len(word)
	}
	return i
}

func (finder *AnagramFinder) loadDict(fh io.Reader) {
	reader := bufio.NewReader(fh)

	line, err := reader.ReadBytes('\n')
	for err != io.EOF {
		finder.words = append(finder.words, line[:len(line)-1])
		line, err = reader.ReadBytes('\n')
	}
	sort.Sort(wordsLength(finder.words))
}
