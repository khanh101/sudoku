package main

import (
	"fmt"
	"time"

	"github.com/khanhhhh/sudoku/gui"
)

func main() {
	seed := int(time.Now().UnixNano())
	portno := 3000
	addr := fmt.Sprintf("0.0.0.0:%d", portno)
	s := gui.NewServer(seed)
	fmt.Printf("Server is up at: http://%s \n", addr)
	s.Run(addr)
	return
}
