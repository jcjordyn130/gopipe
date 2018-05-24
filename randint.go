// randint is a package that just outputs a lot of random numbers to stdout.
package main

import (
	"os"
	"math/rand"
	"fmt"
)

func main() {
	for {
		// Get some random, we do this to bypass integer overflows.
		random1 := rand.Int63()
		random2 := rand.Int63()
		random3 := rand.Int63()
		random4 := rand.Int63()
		random5 := rand.Int63()
		random6 := rand.Int63()
		random7 := rand.Int63()
		random8 := rand.Int63()
		random9 := rand.Int63()
		random10 := rand.Int63()

		// Convert it to a byte array.
		randbyte := []byte(fmt.Sprintf("%d%d%d%d%d%d%d%d%d%d\n", random1, random2, random3, random4, random5, random6, random7, random8, random9, random10))

		// Write it to stdout.
		_, err := os.Stdout.Write(randbyte)
		if err != nil {
			panic(err)
		}
	}
}
