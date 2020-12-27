package gui

import (
	"github.com/gin-gonic/gin"
)

// Server :
type Server interface {
	Run(addr string)
}

// NewServer :
func NewServer() Server {
	s := &server{
		r: gin.Default(),
	}
	s.r.Static("/", "./gui/static/")
	s.r.POST("/api/view", func(c *gin.Context) {
		c.ShouldBindJSON()
	})
}
