package main

import (
	"fmt"
	"time"
)

func main() {
	s := Subscribe(Fetch("thomas.sauvajon.tech"))

	// Close the subscription after some time
	time.AfterFunc(5*time.Second, func() {
		fmt.Println("closed:", s.Close())
	})

	// Print the stream
	for item := range s.Updates() {
		fmt.Println(item.Channel, item.GUID, item.Title)
	}

	// Show the running goroutines
	panic("show me the stacks")
}
