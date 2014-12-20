package main

import (
	"fmt"
	"github.com/adammck/sixaxis"
	"os"
	"time"
)

func main() {

	// open the device
	f, err := os.Open("/dev/input/event0")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// run synchronized in a goroutine
	sa := sixaxis.New(f)
	go sa.Run()

	// Dump the state at 10hz
	c := time.Tick(100 * time.Millisecond)
	for _ = range c {
		fmt.Println(sa)
	}
}
