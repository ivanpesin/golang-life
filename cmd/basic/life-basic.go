package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"time"
)

// Global settings
const initCellCount = 250
const refreshRate = 8
const rows = 30
const cols = 80

var generation = 0

var board [rows][cols]bool

func cls() {
	out, _ := exec.Command("tput", "clear").Output()
	fmt.Printf("%s", out)
}

func pos(r, c int) {
	out, _ := exec.Command("tput", "cup",
		strconv.Itoa(r), strconv.Itoa(c)).Output()
	fmt.Printf("%s", out)
}

func draw() {
	pos(0, 0)

	fmt.Printf("Conway's Life in Go | board %dx%d;", rows, cols)
	fmt.Printf(" rate %d/sec; gen = %d\n", refreshRate, generation)
	fmt.Print("+")
	for i := 0; i < cols; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")

	for row := 0; row < rows; row++ {
		fmt.Print("|")
		for col := 0; col < cols; col++ {
			if board[row][col] {
				fmt.Printf("*")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println("|")
	}

	fmt.Print("+")
	for i := 0; i < cols; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")
}

func copy_board(b [rows][cols]bool) {
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			board[r][c] = b[r][c]
		}
	}
}

func neighbours(r, c int) int {
	count := 0
	if r > 0 && board[r-1][c] {
		count++
	}
	if r > 0 && c > 0 && board[r-1][c-1] {
		count++
	}
	if r > 0 && c < cols-1 && board[r-1][c+1] {
		count++
	}
	if c > 0 && board[r][c-1] {
		count++
	}
	if c < cols-1 && board[r][c+1] {
		count++
	}
	if r < rows-1 && board[r+1][c] {
		count++
	}
	if r < rows-1 && c > 0 && board[r+1][c-1] {
		count++
	}
	if r < rows-1 && c < cols-1 && board[r+1][c+1] {
		count++
	}
	return count
}

func life() {
	var new [rows][cols]bool

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			var n = neighbours(row, col)
			if (n == 2 && board[row][col]) || n == 3 {
				new[row][col] = true
			}
		}
	}

	copy_board(new)
	generation++
}

func init_board() {
	rand.Seed(42) // want games to be repeatable, this static seed
	for i := 0; i < initCellCount; i++ {
		r := rand.Intn(rows)
		c := rand.Intn(cols)
		board[r][c] = true
	}
}

func main() {
	cls()
	init_board()
	for {
		draw()
		life()
		time.Sleep(time.Second / time.Duration(refreshRate))
	}
}
