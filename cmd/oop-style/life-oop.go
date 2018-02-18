package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// CLI parameters
var file = flag.String("f", "", "LIF 1.05/1.06 file")
var deltaXflag = flag.Int("deltax", 0, "X translation for loaded shape (default: auto)")
var deltaYflag = flag.Int("deltay", 0, "Y translation for loaded shape (default: auto)")

var nrows = flag.Int("rows", 22, "number of rows (default: 22)")
var ncols = flag.Int("cols", 78, "number of cols (default: 78)")

var turns = flag.Int("turns", 0, "number of generations to simulate")
var rate = flag.Int("r", 2, "Rate of generations per second (default: 2)")

var ageColor = flag.Bool("color", true, "use color to show cell age (default: true)")
var ageShape = flag.Bool("shape", true, "use shapes to show cell age (default: true)")

// end of CLI parameters

func cls() {
	fmt.Printf("\033[2J")
}

func pos(r, c int) {
	fmt.Printf("\033[" + strconv.Itoa(r) + ";" + strconv.Itoa(c) + "H")
}

type universe struct {
	rows, cols int
	prev       [][]int // previous generation of universe
	board      [][]int // value represents age of a cell
	alive      int     // cells alive
	gen        int     // current generation
}

func newBoard(r, c int) [][]int {
	b := make([][]int, r)
	for i := 0; i < r; i++ {
		b[i] = make([]int, c)
	}
	return b
}

// NewLife returns new initialised universe
func NewLife(r, c int) *universe {

	u := &universe{}

	u.rows = r
	u.cols = c

	u.board = newBoard(r, c)
	u.prev = newBoard(r, c)

	u.alive = 0
	u.gen = 0

	return u
}

func (u *universe) rPentomino() {
	// create r-pentomino pattern
	u.board[u.rows/2][u.cols/2] = 1
	u.board[u.rows/2+1][u.cols/2] = 1
	u.board[u.rows/2+2][u.cols/2] = 1
	u.board[u.rows/2][u.cols/2+1] = 1
	u.board[u.rows/2+1][u.cols/2-1] = 1

	u.alive = 5
}

func cellShape(age int) string {

	if age == 0 {
		return " "
	}

	shape := "*" // "▣"
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

func (u *universe) draw() {
	pos(0, 0)

	fmt.Printf("Conway's Life in Go | board %dx%d;", u.rows, u.cols)
	fmt.Printf(" rate %d/sec; alive = %3d; gen = %d\033[K",
		*rate, u.alive, u.gen)

	// to speed up we're going to update the frame only every 100 generation
	redrawFrame := u.gen%100 == 0
	if redrawFrame {
		fmt.Println()
		fmt.Print("+")
		for i := 0; i < u.cols; i++ {
			fmt.Print("-")
		}
		fmt.Println("+")
	}

	for r := 0; r < u.rows; r++ {
		if redrawFrame {
			fmt.Print("|")
		}
		for c := 0; c < u.cols; c++ {
			if u.board[r][c] != u.prev[r][c] {
				pos(r+3, c+2)

				fmt.Print(cellShape(u.board[r][c]))
			}
		}
		if redrawFrame {
			pos(r+3, u.cols+2)
			fmt.Println("|")
		}
	}
	if redrawFrame {
		fmt.Print("+")
		for i := 0; i < u.cols; i++ {
			fmt.Print("-")
		}
		fmt.Print("+")
	}
	pos(u.rows+3, 0)
}

func (u *universe) neighbours(r, c int) (count int) {
	if r > 0 && c > 0 && u.prev[r-1][c-1] > 0 {
		count++
	}
	if r > 0 && c < u.cols-1 && u.prev[r-1][c+1] > 0 {
		count++
	}
	if r > 0 && u.prev[r-1][c] > 0 {
		count++
	}

	if c > 0 && u.prev[r][c-1] > 0 {
		count++
	}
	if c < u.cols-1 && u.prev[r][c+1] > 0 {
		count++
	}

	if r < u.rows-1 && c > 0 && u.prev[r+1][c-1] > 0 {
		count++
	}
	if r < u.rows-1 && c < u.cols-1 && u.prev[r+1][c+1] > 0 {
		count++
	}
	if r < u.rows-1 && u.prev[r+1][c] > 0 {
		count++
	}

	return
}

func (u *universe) evolve() {
	u.prev = u.board

	u.board = newBoard(u.rows, u.cols)
	u.alive = 0

	for r := 0; r < u.rows; r++ {
		for c := 0; c < u.cols; c++ {
			n := u.neighbours(r, c)
			if (n == 2 && u.prev[r][c] > 0) || n == 3 {
				u.board[r][c] = u.prev[r][c] + 1
				u.alive++
			}
		}
	}

	u.gen++
}

func (u *universe) loadShape105(lines []string) {
	startX := u.cols/2 - 1
	startY := u.rows/2 - 1

	r := 0 // current row as we read the shape
	for ln, l := range lines {
		if len(l) > 1 && l[0] == '#' {
			if l[1] == 'P' {
				f := strings.Split(l, " ")
				atoi, err := strconv.Atoi(f[1])
				if err != nil {
					log.Fatalf("Error parsing shape in line %v: %v\n", ln, err)
				}
				startX += atoi
				atoi, err = strconv.Atoi(strings.TrimSuffix(f[2], "\r"))
				if err != nil {
					log.Fatalf("Error parsing shape in line %v: %v\n", ln, err)
				}
				startY += atoi
			}
			continue
		}

		for j, c := range l {
			if c == '*' {
				if startY+r < 0 || startY+r >= u.rows {
					fmt.Println("ERROR: Not enough rows to render the shape")
					os.Exit(2)
				}
				if startX+j < 0 || startX+j >= u.cols {
					fmt.Println("ERROR: Not enough cols to render the shape")
					os.Exit(2)
				}
				u.board[startY+r][startX+j] = 1
				u.alive++
			}
		}
		r++
	}

}

func (u *universe) loadShape106(lines []string) {
	centerX := u.cols/2 - 1
	centerY := u.rows/2 - 1

	for ln, l := range lines {
		if len(l) > 0 && l[0] == '#' {
			continue
		}

		f := strings.Split(l, " ")
		if len(f) < 2 {
			continue
		}
		atoi, err := strconv.Atoi(f[0])
		if err != nil {
			log.Fatalf("Error parsing shape in line %v: %v\n", ln, err)
		}
		x := atoi
		atoi, err = strconv.Atoi(strings.Trim(f[1], "\r"))
		if err != nil {
			log.Fatalf("Error parsing shape in line %v: %v\n", ln, err)
		}
		y := atoi

		u.board[centerY+y][centerX+x] = 1
		u.alive++
	}
}

func (u *universe) boundries() (startX, startY, endX, endY int) {
	startX, endX, startY, endY = u.cols, 0, u.rows, 0

	// calculate shape boundries
	for r := 0; r < u.rows; r++ {
		for c := 0; c < u.cols; c++ {
			if u.board[r][c] > 0 {
				if r < startY {
					startY = r
				}
				if r > endY {
					endY = r
				}
				if c < startX {
					startX = c
				}
				if c > endX {
					endX = c
				}
			}
		}
	}

	return
}

func (u *universe) translate() {

	startX, startY, endX, endY := u.boundries()

	deltaX := *deltaXflag
	deltaY := *deltaYflag

	if deltaX == 0 && deltaY == 0 {
		// recenter the shape
		// calculate translation
		currentCenterX := (startX + endX) / 2
		currentCenterY := (startY + endY) / 2

		centerX := u.cols/2 - 1
		centerY := u.rows/2 - 1

		deltaX = centerX - currentCenterX
		deltaY = centerY - currentCenterY
	}

	// translate board
	t := newBoard(u.rows, u.cols)
	for r := startY; r <= endY; r++ {
		for c := startX; c <= endX; c++ {
			if u.board[r][c] > 0 {
				t[r+deltaY][c+deltaX] = u.board[r][c]
			}
		}
	}
	u.board = t

	// log.Printf("%v", t)
	// log.Printf("D: deltaX = %v ; deltaY = %v", deltaX, deltaY)
	// log.Printf("D: startX = %v ; endX = %v", startX, endX)
	// log.Fatalf("D: startY = %v ; endY = %v", startY, endY)

}

func (u *universe) loadFrom(fn string) {

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")

	if len(lines) < 1 {
		log.Fatalf("File %s appears to be empty", fn)
	}

	// detect file format
	is105 := regexp.MustCompile("^\\s*#\\s*Life\\s+1\\.05")
	is106 := regexp.MustCompile("^\\s*#\\s*Life\\s+1\\.06")

	switch {
	case is105.MatchString(lines[0]):
		u.loadShape105(lines)
	case is106.MatchString(lines[0]):
		u.loadShape106(lines)
	default:
		log.Fatalf("Invalid file format: %v", fn)
	}

	u.translate()

}

func main() {

	flag.Parse()

	cls()

	life := NewLife(*nrows, *ncols)

	if *file == "" {
		life.rPentomino()
	} else {
		life.loadFrom(*file)
	}

	for {
		life.draw()
		life.evolve()
		if *turns > 0 && life.gen >= *turns {
			fmt.Printf("\nReached generation %v\n", life.gen)
			break
		}
		time.Sleep(time.Second / time.Duration(*rate))
	}
}