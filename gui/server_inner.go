package gui

import "github.com/gin-gonic/gin"

type server struct {
	r *gin.Engine
	pool 
}

func (s *server) Run(addr string) {
	s.r.Run(addr)
}
