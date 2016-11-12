package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// Global settings
var refreshRate = flag.Int("r", 8, "refresh rate per sec")
var initCellCount = flag.Int("i", 100, "initial number of alive cells")

var ageShape = flag.Bool("shape", true, "use shapes to denote cell age")
var ageColor = flag.Bool("color", true, "use color to denote cell age")

var rows = flag.Int("rows", 23, "number of rows in universe")
var cols = flag.Int("cols", 78, "number of columns in universe")

var generation = 0
var alive = 0

var board [][]int // value represents age of a cell

func cls() {
	out, _ := exec.Command("tput", "clear").Output()
	fmt.Printf("%s", out)
}

func pos(r, c int) {
	// out, _ := exec.Command("tput", "cup",
	// 	strconv.Itoa(r), strconv.Itoa(c)).Output()
	fmt.Printf("\033[" + strconv.Itoa(r) + ";" + strconv.Itoa(c) + "H")
}

func initSlice(r, c int) [][]int {
	new := make([][]int, r)
	for i := 0; i < r; i++ {
		new[i] = make([]int, c)
	}

	return new
}

func cell(age int) string {

	if age == 0 {
		return " "
	}

	shape := "*"
	if *ageShape {
		switch age {
		case 1:
			shape = "."
		case 2:
			shape = "∘"
		case 3:
			shape = "∙"
		default:
			shape = "*"
		}
	}

	if *ageColor {
		switch age {
		case 1:
			return "\033[1;32m" + shape + "\033[0m"
		case 2:
			return "\033[1;36m" + shape + "\033[0m"
		case 3:
			return "\033[1;31m" + shape + "\033[0m"
		case 4:
			return "\033[1;35m" + shape + "\033[0m"
		default:
			return "\033[0;33m" + shape + "\033[0m"
		}
	}

	return shape
}

func draw() {
	pos(0, 0)

	fmt.Printf("Conway's Life in Go | board %dx%d;", *rows, *cols)
	fmt.Printf(" rate %d/sec; alive = %2d; gen = %d   \n",
		*refreshRate, alive, generation)
	fmt.Print("+")
	for i := 0; i < *cols; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")

	for row := 0; row < *rows; row++ {
		fmt.Print("|")
		for col := 0; col < *cols; col++ {
			fmt.Print(cell(board[row][col]))
		}
		fmt.Println("|")
	}

	fmt.Print("+")
	for i := 0; i < *cols; i++ {
		fmt.Print("-")
	}
	fmt.Println("+")
}

func copyBoard(b [][]int) {
	for r := 0; r < *rows; r++ {
		for c := 0; c < *cols; c++ {
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
	if r > 0 && c < *cols-1 && board[r-1][c+1] > 0 {
		count++
	}
	if c > 0 && board[r][c-1] > 0 {
		count++
	}
	if c < *cols-1 && board[r][c+1] > 0 {
		count++
	}
	if r < *rows-1 && board[r+1][c] > 0 {
		count++
	}
	if r < *rows-1 && c > 0 && board[r+1][c-1] > 0 {
		count++
	}
	if r < *rows-1 && c < *cols-1 && board[r+1][c+1] > 0 {
		count++
	}
	return count
}

func life() {
	new := initSlice(*rows, *cols)
	alive = 0

	for row := 0; row < *rows; row++ {
		for col := 0; col < *cols; col++ {
			var n = neighbours(row, col)
			if (n == 2 && board[row][col] > 0) || n == 3 {
				new[row][col] = board[row][col] + 1
				alive++
			}
		}
	}

	copyBoard(new)
	generation++
}

func initBoard() {

	// board[*rows/2][*cols/2] = 1
	// board[*rows/2+1][*cols/2] = 1
	// board[*rows/2+2][*cols/2] = 1
	// board[*rows/2][*cols/2+1] = 1
	// board[*rows/2+1][*cols/2-1] = 1

	// board[*rows/2][*cols/2] = 1
	// board[*rows/2][*cols/2+1] = 1
	// board[*rows/2][*cols/2+2] = 1
	// board[*rows/2+1][*cols/2+2] = 1
	// board[*rows/2+1][*cols/2-1] = 1
	// board[*rows/2+2][*cols/2-2] = 1
	// board[*rows/2+3][*cols/2-3] = 1

	//return

	rand.Seed(42) // want games to be repeatable, this static seed
	for i := 0; i < *initCellCount; i++ {
		r := rand.Intn(*rows)
		c := rand.Intn(*cols)
		board[r][c] = 1
	}
}

func pause() {
	if *refreshRate > 0 {
		time.Sleep(time.Second / time.Duration(*refreshRate))
	} else {
		fmt.Print("Press Enter to advance to next generation")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func main() {
	flag.Parse()

	board = initSlice(*rows, *cols)
	cls()
	initBoard()
	for {
		draw()
		life()
		pause()
	}
}
