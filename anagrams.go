package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
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

func loadDict(fh io.Reader, words *[][]byte) {
	reader := bufio.NewReader(fh)

	line, err := reader.ReadBytes('\n')
	for err != io.EOF {
		*words = append(*words, bytes.ToLower(line[:len(line)-1]))
		line, err = reader.ReadBytes('\n')
	}
}

func findAnagrams(words [][]byte) map[string][][]byte {
	anagrams := make(map[string][][]byte)
	sortWord := func(word []byte) string {
		sortedBytes := make([]byte, len(word))
		copy(sortedBytes, word)
		sort.Sort(byteArray(sortedBytes))
		return string(sortedBytes[:])
	}

	for _, word := range words {
		sortedBytesAsString := sortWord(word)
		_, ok := anagrams[sortedBytesAsString]
		if !ok {
			anagrams[sortedBytesAsString] = make([][]byte, 0, 1)
		}
		anagrams[sortedBytesAsString] = append(anagrams[sortedBytesAsString], word)
	}
	return anagrams
}

func displayAnagrams(anagrams map[string][][]byte) {
	for _, words := range anagrams {
		if len(words) > 1 {
			var buffer bytes.Buffer
			buffer.WriteString(string(words[0][:]))
			buffer.WriteString(":")
			for _, word := range words[1:] {
				buffer.WriteString(" ")
				buffer.WriteString(string(word[:]))
			}
			fmt.Println(buffer.String())
		}
	}
}

func nextSlice(words [][]byte, prev int) ([][]byte, int) {
	start_len := len(words[prev+1])
	i := prev
	for i < len(words) && len(words[i]) == start_len {
		i++
	}
	return words[prev+1 : i], i
}

func main() {
	fh, err := os.Open(dictPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer fh.Close()

	words := make([][]byte, 1000)
	loadDict(fh, &words)
	sort.Sort(wordsLength(words))
	slice, prev := nextSlice(words, 0)

	for prev < len(words) {
		go displayAnagrams(findAnagrams(slice))
		slice, prev = nextSlice(words, prev)
	}
}
