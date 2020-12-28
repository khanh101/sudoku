package gui

import (
	"github.com/gin-gonic/gin"
)

// N : board size
const N = 3

type server struct {
	r *gin.Engine
	s *session
}

func (s *server) Run(addr string) {
	s.r.Run(addr)
}
