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
	Run(addr string)
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
		key := KeyView{}
		if err := c.BindJSON(&key); err == nil {
			value := s.s.get(key.Key)
			if value != nil {
				c.JSON(http.StatusOK, key)
				return
			}
		}
		intKey, err := strconv.Atoi(key.Key)
		if err != nil {
			intKey = s.rand.Int()
			key.Key = strconv.Itoa(intKey)
		}
		current := sudoku.Generate(N, intKey)
		s.s.set(key.Key, sudoku.NewGameWithBoard(N, current))
		c.JSON(http.StatusOK, key)
	})
	s.r.POST("/api/new_board", func(c *gin.Context) {
		board := BoardView{}
		key := KeyView{}
		if err := c.BindJSON(&board); err == nil {
			value := s.s.get(board.Board)
			if value != nil {
				key.Key = value.(string)
				c.JSON(http.StatusOK, key)
				return
			}
		}
		if len(board.Board) >= N*N*N*N {
			var err error
			current := sudoku.NewBoard(N)
			count := 0
			for row := 0; row < N*N; row++ {
				for col := 0; col < N*N; col++ {
					current[row][col], err = strconv.Atoi(board.Board[count : count+1])
					if err != nil {
						break
					}
					count++
				}
				if err != nil {
					break
				}
			}
			if err == nil {
				key.Key = board.Board[0 : N*N*N*N]
				s.s.set(key.Key, sudoku.NewGameWithBoard(N, current))
				c.JSON(http.StatusOK, key)
				return
			}
		}
		intKey := s.rand.Int()
		key.Key = strconv.Itoa(intKey)
		s.s.set(key.Key, sudoku.NewGame(N, intKey))
		c.JSON(http.StatusOK, key)
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
	s.r.POST("/api/place", func(c *gin.Context) {
		pos := PosView{}
		if err := c.BindJSON(&pos); err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		value := s.s.get(pos.Key)
		if value == nil {
			c.JSON(http.StatusNotFound, nil)
			return
		}
		game := value.(sudoku.Game)
		if game == nil {
			panic("wrong")
		}
		game.Place(sudoku.PlacementView{
			Row: pos.Row,
			Col: pos.Col,
			Val: pos.Val,
		})
		c.JSON(http.StatusOK, nil)
	})
	s.r.POST("/api/undo", func(c *gin.Context) {
		pos := PosView{}
		if err := c.BindJSON(&pos); err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		value := s.s.get(pos.Key)
		if value == nil {
			c.JSON(http.StatusNotFound, nil)
			return
		}
		game := value.(sudoku.Game)
		if game == nil {
			panic("wrong")
		}
		ok, view := game.Undo()
		if !ok {
			c.JSON(http.StatusOK, nil)
			return
		}
		c.JSON(http.StatusOK, view)
	})
	s.r.POST("/api/implication", func(c *gin.Context) {
		pos := PosView{}
		if err := c.BindJSON(&pos); err != nil {
			c.JSON(http.StatusBadRequest, nil)
			return
		}
		value := s.s.get(pos.Key)
		if value == nil {
			c.JSON(http.StatusNotFound, nil)
			return
		}
		game := value.(sudoku.Game)
		if game == nil {
			panic("wrong")
		}
		ok, view := game.Implication()
		if !ok {
			c.JSON(http.StatusOK, nil)
			return
		}
		c.JSON(http.StatusOK, view)
	})
	s.r.POST("/api/access", func(c *gin.Context) {
		key := PosView{}
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
