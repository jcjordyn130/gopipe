// gopipe implements a buffer around pipes.
// example: cat some_huge_file_over_the_network | go run gopipe.go | some_consumer
package main

import (
	"syscall"
	"os"
	"bytes"
)

// bytecopy returns a copy of a bytearray with all the null bytes removed.
func bytecopy(array []byte) ([]byte) {
	// Make a new array of the same size.
	newArray := make([]byte, len(array))

	// Copy the data.
	copy(newArray, array)

	// Remove all the null bytes.
	newArray = bytes.Trim(newArray, "\x00")

	// Return it.
	return newArray
}

// readStdin() reads from the producer.
func readStdin(data chan[]byte) {
	// Make the buffer.
	buf := make([]byte, os.Getpagesize())

	for {
		// Read stdin.
		_, err := os.Stdin.Read(buf)
		if err != nil {
			// The producer closed the pipe, signal the consumer and die.
			if err.Error() == "EOF" {
				data <- []byte{}
				break
			}

			// We recieved an error, panic time
			panic(err)
		}

		// Write the bytes we read to the channel.
		data <- bytecopy(buf)
	}
}

// writeStdout() writes to the consumer.
func writeStdout(data chan[]byte) {
	for {
		// Get the next element.
		buf := <- data

		// The producer closed the pipe, die.
		if len(buf) == 0 {
			os.Stdout.Close()
			os.Exit(0)
		}

		// Write the buffer to stdout.
		_, err := os.Stdout.Write(buf)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	// Make the channel.
	data := make(chan []byte, 8096)

	// Run stuff.
	go readStdin(data)
	go writeStdout(data)

	// Wait forever, if we need to die one of the goroutines will do it.
	syscall.Pause()
}
