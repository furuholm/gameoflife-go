package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

type Board struct {
	Width  int
	Height int
	Cells  [][]bool
}

func makeCells(width, height int) [][]bool {
	cells := make([][]bool, height)
	for i := range cells {
		cells[i] = make([]bool, width)
	}
	return cells
}

func newBoard(width, height int, percentage float64) Board {
	cells := makeCells(width, height)
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if r.Float64() < percentage {
				cells[i][j] = true
			}
		}
	}
	return Board{
		Width:  width,
		Height: height,
		Cells:  cells,
	}
}

func (b *Board) print() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
	for i := 0; i < b.Height; i++ {
		for j := 0; j < b.Width; j++ {
			if b.isAlive(i, j) {
				fmt.Print(" X ")
			} else {
				fmt.Print("   ")
			}
		}
		fmt.Print("\n")
	}
}

func (b *Board) isAlive(row, col int) bool {
	if row < 0 || col < 0 || col >= b.Width || row >= b.Height {
		return false // out of bounds
	} else {
		return b.Cells[row][col]
	}
}

func (b *Board) makeAlive(row, col int) {
	b.Cells[row][col] = true
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
	buffer := makeCells(b.Width, b.Height)
	for i := 0; i < b.Height; i++ {
		for j := 0; j < b.Width; j++ {
			neighbours := b.countNeighbours(i, j)
			if b.isAlive(i, j) {
				buffer[i][j] = neighbours == 2 || neighbours == 3
			} else if !b.isAlive(i, j) {
				buffer[i][j] = neighbours == 3
			}
		}
	}
	b.Cells = buffer
}

func main() {
	board := newBoard(80, 80, 0.2)
	for {
		board.print()
		time.Sleep(1 * time.Second)
		board.step()
	}
}
