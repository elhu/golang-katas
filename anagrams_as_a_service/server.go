package main

import (
	"bufio"
	"fmt"
	"github.com/elhu/katas/anagrams_as_a_service/anagram"
	"net"
	"os"
)

func handleClient(conn net.Conn, finder *anagram.AnagramFinder) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			conn.Close()
			return
		}
		newMessage := finder.AnagramsFor([]byte(message))
		for _, word := range newMessage {
			conn.Write(word)
			conn.Write([]byte("\n"))
		}
	}
}

func main() {
	fmt.Print("Loading anagrams...")
	anagramFinder := anagram.New()
	fmt.Println(" done!")

	listener, err := net.Listen("tcp", "localhost:4567")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Socket listening")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("New client accepted")
		go handleClient(conn, anagramFinder)
	}
}
