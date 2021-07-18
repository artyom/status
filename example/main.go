package main

import (
	"fmt"
	"time"

	"github.com/artyom/status"
)

func main() {
	fmt.Println("start")        // These will be printed regardless of whether
	defer fmt.Println("finish") // stdout is terminal or not.

	line := new(status.Line)
	defer line.Done() // Try replacing with: defer line.Print("done!\n")
	const total = 15
	for i := 0; i < total; i++ {
		word := "odd"
		if i%2 == 0 {
			word = "even"
		}
		// This is only printed if stdout is connected to the terminal.
		// Try redirecting program's output to the file and see how it works.
		line.Printf("step %d (%s), %d left to do", i, word, total-i)
		time.Sleep(300 * time.Millisecond)
	}
}
