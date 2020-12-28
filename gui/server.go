package gui

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Server :
type Server interface {
	Run(addr string)
}

// Timeout :
const Timeout = time.Duration(60 * time.Second)

// NewServer :
func NewServer() Server {
	s := &server{
		r: gin.Default(),
		s: newSession(Timeout),
	}
	return s
}
