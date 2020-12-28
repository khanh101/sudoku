package gui

import (
	"github.com/gin-gonic/gin"
)

// Server :
type Server interface {
	Run(addr string)
}

// Timeout :
const Timeout = 60

// NewServer :
func NewServer() Server {
	s := &server{
		r: gin.Default(),
		s: newSession(Timeout),
	}
	return s
}
