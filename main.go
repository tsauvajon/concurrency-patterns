package main

import (
	"fmt"
	"time"
)

func main() {
	m := Merge(
		Subscribe(Fetch("thomas.sauvajon.tech")),
		Subscribe(Fetch("linkedin.com/in/tsauvajon")),
		Subscribe(Fetch("github.com/tsauvajon")),
	)

	// Close the subscription after some time
	time.AfterFunc(2*time.Second, func() {
		fmt.Println("closed:", m.Close())
	})

	// Print the stream
	for item := range m.Updates() {
		fmt.Println(item.Channel, item.GUID, item.Title)
	}

	// Show the running goroutines
	panic("show me the stacks")
}
