package main

import (
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Cell struct {
	Row int
	Col int
}

type Board struct {
	Cells map[Cell]bool
}

var (
	cellSize = float64(40)
)

func newRandomBoard(width, height int, percentage float64) *Board {
	cells := map[Cell]bool{}
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if r.Float64() < percentage {
				cells[Cell{Row: i, Col: j}] = true
			}
		}
	}
	return &Board{
		Cells: cells,
	}
}

func (b *Board) isAlive(row, col int) bool {
	exists, alive := b.Cells[Cell{Row: row, Col: col}]
	return exists && alive
}

func (b *Board) makeAlive(row, col int) {
	b.Cells[Cell{Row: row, Col: col}] = true
}

func (b *Board) kill(row, col int) {
	b.Cells[Cell{Row: row, Col: col}] = false
}

func (b *Board) countNeighbours(row, col int) int {
	count := 0
	for i := row - 1; i <= row+1; i++ {
		for j := col - 1; j <= col+1; j++ {
			if i == row && j == col {
				continue
			}
			if b.isAlive(i, j) {
				count++
			}
		}
	}
	return count
}

func (b *Board) step() {
	buffer := map[Cell]bool{}
	for cell, _ := range b.Cells {
		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				row := cell.Row + i
				col := cell.Col + j
				neighbours := b.countNeighbours(row, col)
				if b.isAlive(row, col) && (neighbours == 2 || neighbours == 3) {
					buffer[Cell{Row: row, Col: col}] = true
				} else if !b.isAlive(row, col) && neighbours == 3 {
					buffer[Cell{Row: row, Col: col}] = true
				}
			}
		}
	}
	b.Cells = buffer
}

func draw(b *Board, imd *imdraw.IMDraw, win *pixelgl.Window) {
	imd.Clear()
	padding := float64(1)
	width := win.Bounds().Max.X
	height := win.Bounds().Max.Y
	imd.Color = pixel.Alpha(0.4) // colornames.White
	for x := 0.0; x <= width; x += cellSize {
		imd.Push(pixel.V(x, 0))
		imd.Push(pixel.V(x, height))
		imd.Line(2)
	}
	for y := height; y >= 0; y -= cellSize {
		imd.Push(pixel.V(0, y))
		imd.Push(pixel.V(width, y))
		imd.Line(2)
	}
	imd.Color = colornames.Black
	for cell, alive := range b.Cells {
		if alive {
			x := float64(cell.Col)*cellSize + padding
			y := height - float64(cell.Row)*cellSize - padding
			imd.Push(pixel.V(x, y))
			imd.Push(pixel.V(x+cellSize-padding*2, y-cellSize+padding*2))
			imd.Rectangle(0)
		}
	}
}

func posToCell(pos pixel.Vec, screenHeight float64) (row, col int) {
	col = int((pos.X - 2.0) / cellSize)
	row = int((screenHeight - pos.Y - 2.0) / cellSize)
	return
}

func emptyBoard() *Board {
	return &Board{Cells: map[Cell]bool{}}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:     "Game of life",
		Bounds:    pixel.R(0, 0, 1024, 768),
		VSync:     true,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	imd := imdraw.New(nil)
	run := false
	speed := int64(1)
	tick := time.Tick(time.Second / time.Duration(speed))
	board := emptyBoard()
	for !win.Closed() {
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			row, col := posToCell(win.MousePosition(), win.Bounds().Max.Y)
			if board.isAlive(row, col) {
				board.kill(row, col)
			} else {
				board.makeAlive(row, col)
			}
		}
		// Initialize
		if win.JustPressed(pixelgl.KeyI) {
			row, col := posToCell(pixel.Vec{win.Bounds().Max.X, 0}, win.Bounds().Max.Y)
			board = newRandomBoard(col, row, 0.15)
		}
		// Run/stop
		if win.JustPressed(pixelgl.KeyR) {
			run = !run
		}
		// Step
		if win.JustPressed(pixelgl.KeyS) {
			board.step()
		}
		// Clear
		if win.JustPressed(pixelgl.KeyC) {
			board = emptyBoard()
		}
		// Quit
		if win.JustPressed(pixelgl.KeyQ) {
			win.SetClosed(true)
		}
		// Speed up
		if win.JustPressed(pixelgl.KeyRightBracket) {
			speed += 1
			tick = time.Tick(time.Second / time.Duration(speed))
		}
		// Slow down
		if win.JustPressed(pixelgl.KeyLeftBracket) {
			// Speed zero is not allowed
			if speed > 1 {
				speed -= 1
				tick = time.Tick(time.Second / time.Duration(speed))
			}
		}
		// Zoom in
		if win.JustPressed(pixelgl.KeyZ) {
			cellSize *= 2
		}
		// Zoom out
		if win.JustPressed(pixelgl.KeyX) {
			cellSize /= 2
		}

		draw(board, imd, win)
		win.Clear(colornames.Skyblue)
		imd.Draw(win)
		win.Update()
		select {
		case <-tick:
			if run {
				board.step()
			}
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
