package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	time.Sleep(1 * time.Second)
	_, err := os.Stdout.Write([]byte("hello, world.\n"))
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
}
