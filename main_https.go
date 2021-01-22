package main

import (
	"log"
	"time"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/khanhhhh/sudoku/gui"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	seed := int(time.Now().UnixNano())
	s := gui.NewServer(seed)

	log.Fatal(autotls.Run(s.Handler(), "personal.khanh3am.org"))
}
