package gui

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/khanhhhh/sudoku/sudoku"
)

// Server :
type Server interface {
	Handler() *gin.Engine
}

// Timeout :
const Timeout = time.Duration(60 * time.Second)

// NewServer :
func NewServer(seed int) Server {
	s := &server{
		r:    gin.Default(),
		s:    newSession(Timeout),
		rand: rand.New(rand.NewSource(int64(seed))),
	}
	// init (3*3 x 3*3) board
	sudoku.ReduceBase(N)

	s.r.Static("/", "./gui/static/")
	s.r.POST("/api/new", func(c *gin.Context) {
		boardView := BoardView{}
		if err := c.BindJSON(&boardView); err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}

		board, ok := sudoku.FromString(N, boardView.Board)
		intKey := s.rand.Int()
		if !ok {
			board = sudoku.Generate(N, intKey)
		}
		key := KeyView{
			Key: strconv.Itoa(intKey),
		}
		game, ok := sudoku.NewGame(N, board)
		if !ok {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		s.s.set(key.Key, game)
		c.JSON(http.StatusOK, key)
	})
	s.r.POST("/api/login", func(c *gin.Context) {
		key := KeyView{}
		if err := c.BindJSON(&key); err == nil {
			value := s.s.get(key.Key)
			if value != nil {
				c.JSON(http.StatusOK, key)
				return
			}
		}
		c.JSON(http.StatusBadRequest, nil)
	})
	s.r.POST("/api/view", func(c *gin.Context) {
		key := KeyView{}
		if err := c.BindJSON(&key); err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		value := s.s.get(key.Key)
		if value == nil {
			c.JSON(http.StatusNotFound, nil)
			return
		}
		game := value.(sudoku.Game)
		if game == nil {
			panic("wrong")
		}
		c.JSON(http.StatusOK, game.View())
		return
	})
	s.r.POST("/api/point", func(c *gin.Context) {
		point := PointView{}
		if err := c.BindJSON(&point); err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		value := s.s.get(point.Key)
		if value == nil {
			c.JSON(http.StatusNotFound, nil)
			return
		}
		game := value.(sudoku.Game)
		if game == nil {
			panic("wrong")
		}
		game.Point(sudoku.CellView{
			Row: point.Row,
			Col: point.Col,
		})
		c.JSON(http.StatusOK, nil)
	})
	s.r.POST("/api/place", func(c *gin.Context) {
		place := PlaceView{}
		if err := c.BindJSON(&place); err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		value := s.s.get(place.Key)
		if value == nil {
			c.JSON(http.StatusNotFound, nil)
			return
		}
		game := value.(sudoku.Game)
		if game == nil {
			panic("wrong")
		}
		game.Place(place.Val)
		c.JSON(http.StatusOK, nil)
	})
	s.r.POST("/api/undo", func(c *gin.Context) {
		key := KeyView{}
		if err := c.BindJSON(&key); err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		value := s.s.get(key.Key)
		if value == nil {
			c.JSON(http.StatusNotFound, nil)
			return
		}
		game := value.(sudoku.Game)
		if game == nil {
			panic("wrong")
		}
		game.Undo()
		c.JSON(http.StatusOK, nil)
	})
	s.r.POST("/api/implication", func(c *gin.Context) {
		key := KeyView{}
		if err := c.BindJSON(&key); err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		value := s.s.get(key.Key)
		if value == nil {
			c.JSON(http.StatusNotFound, nil)
			return
		}
		game := value.(sudoku.Game)
		if game == nil {
			panic("wrong")
		}
		game.Implication()
		c.JSON(http.StatusOK, nil)
	})
	s.r.POST("/api/access", func(c *gin.Context) {
		key := KeyView{}
		if err := c.BindJSON(&key); err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		s.s.get(key.Key)
	})
	s.r.POST("/api/global_stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"number of active users": s.s.numActiveKey(),
		})
	})
	return s
}
