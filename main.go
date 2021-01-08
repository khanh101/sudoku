package main

import (
	"fmt"
	"time"

	"github.com/khanhhhh/sudoku/gui"
)

func main() {
	seed := int(time.Now().UnixNano())
	portno := 8080
	s := gui.NewServer(seed)
	fmt.Printf("Server is up at: http://0.0.0.0:%d\n", portno)
	s.Run(fmt.Sprintf(":%d", portno))
	return
}
