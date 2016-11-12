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
const color = true

var generation = 0

var board [rows][cols]int // value represents age of a cell

func cls() {
	out, _ := exec.Command("tput", "clear").Output()
	fmt.Printf("%s", out)
}

func pos(r, c int) {
	out, _ := exec.Command("tput", "cup",
		strconv.Itoa(r), strconv.Itoa(c)).Output()
	fmt.Printf("%s", out)
}

func cell(age int) string {

	if !color {
		if age > 0 {
			return "*"
		} else {
			return " "
		}
	}

	switch age {
	case 0:
		return " "
	case 1:
		return "\033[1;32m*\033[0m"
	case 2:
		return "\033[1;36m*\033[0m"
	case 3:
		return "\033[1;31m*\033[0m"
	case 4:
		return "\033[1;35m*\033[0m"
	default:
		return "\033[0;33m*\033[0m"
	}
	return " \033[0m"
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
			fmt.Print(cell(board[row][col]))
		}
		fmt.Println("|")
	}

	fmt.Print("+")
	for i := 0; i < cols; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")
}

func copyBoard(b [rows][cols]int) {
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			board[r][c] = b[r][c]
		}
	}
}

func neighbours(r, c int) int {
	count := 0
	if r > 0 && board[r-1][c] > 0 {
		count++
	}
	if r > 0 && c > 0 && board[r-1][c-1] > 0 {
		count++
	}
	if r > 0 && c < cols-1 && board[r-1][c+1] > 0 {
		count++
	}
	if c > 0 && board[r][c-1] > 0 {
		count++
	}
	if c < cols-1 && board[r][c+1] > 0 {
		count++
	}
	if r < rows-1 && board[r+1][c] > 0 {
		count++
	}
	if r < rows-1 && c > 0 && board[r+1][c-1] > 0 {
		count++
	}
	if r < rows-1 && c < cols-1 && board[r+1][c+1] > 0 {
		count++
	}
	return count
}

func life() {
	var new [rows][cols]int

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			var n = neighbours(row, col)
			if (n == 2 && board[row][col] > 0) || n == 3 {
				new[row][col] = board[row][col] + 1
			}
		}
	}

	copyBoard(new)
	generation++
}

func init_board() {
	rand.Seed(42) // want games to be repeatable, this static seed
	for i := 0; i < initCellCount; i++ {
		r := rand.Intn(rows)
		c := rand.Intn(cols)
		board[r][c] = 1
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
