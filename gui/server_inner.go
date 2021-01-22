package gui

import (
	"math/rand"

	"github.com/gin-gonic/gin"
)

// N : board size
const N = 3

type server struct {
	r    *gin.Engine
	s    *session
	rand *rand.Rand
}

func (s *server) Handler() *gin.Engine {
	return s.r
}
