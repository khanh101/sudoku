package main

import (
	"fmt"
	"github.com/khanh-nguyen-code/sudoku/gui"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	seed := int(time.Now().UnixNano())
	portno := 3000
	addr := fmt.Sprintf("0.0.0.0:%d", portno)
	s := gui.NewServer(seed)
	fmt.Printf("Server is up at: http://%s/sudoku/ \n", addr)
	_ = s.Handler().Run(addr)
	return
}
