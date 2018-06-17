// gopipe implements a buffer around pipes.
// example: cat some_huge_file_over_the_network | go run gopipe.go | some_consumer
package main

import (
	"syscall"
	"os"
)

// byteCopy() returns a copy of the byte array given.
func byteCopy(array []byte) ([]byte) {
	// Make a new array of the same size.
	newArray := make([]byte, len(array))

	// Copy the data.
	copy(newArray, array)

	return newArray
}

// readStdin() reads from the producer.
func readStdin(data chan[]byte) {
	// Make the buffer.
	// We use the size of a page because it results in more speed.
	buf := make([]byte, os.Getpagesize())

	for {
		// Read stdin.
		_, err := os.Stdin.Read(buf)
		if err != nil {
			// The producer closed the pipe, close the channel and die.
			if err.Error() == "EOF" {
				close(data)
				break
			}

			// We recieved an error, panic time
			panic(err)
		}

		// Write the bytes we read to the channel.
		data <- byteCopy(buf)
	}
}

// writeStdout() writes to the consumer.
func writeStdout(data chan[]byte) {
	for {
		// Get the next element.
		buf := <- data

		// The producer closed the pipe, die.
		if buf == nil {
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
	// We use channels here so the consumer doesn't have to wait for all the data.
	// Plus we get built in buffering.
	// This is useful when streaming over the network for example.
	data := make(chan []byte, 8096)

	// Start reading stdin and stdout.
	go readStdin(data)
	go writeStdout(data)

	// Wait forever, if we need to die one of the goroutines will do it.
	syscall.Pause()
}
