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

// NewLife returns new initialised universe
func NewLife(r, c int) *universe {

	u := &universe{}

	u.rows = r
	u.cols = c

	u.board = make([][]int, r)
	u.prev = make([][]int, r)

	for i := 0; i < r; i++ {
		u.board[i] = make([]int, c)
		u.prev[i] = make([]int, c)
	}

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

func (u *universe) draw() {
	pos(0, 0)

	fmt.Printf("Conway's Life in Go | board %dx%d;", u.rows, u.cols)
	fmt.Printf(" rate %d/sec; alive = %3d; gen = %d\033[K",
		2, u.alive, u.gen)

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
			switch {
			case u.board[r][c] == 0 && u.prev[r][c] != 0:
				pos(r+3, c+2)
				fmt.Printf(" ")
			case u.board[r][c] > 0 && u.prev[r][c] == 0:
				pos(r+3, c+2)
				fmt.Printf("*")
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

	u.board = make([][]int, u.rows)
	u.alive = 0

	for r := 0; r < u.rows; r++ {
		u.board[r] = make([]int, u.cols)
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
	}
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

}

func main() {

	file := flag.String("f", "", "LIF 1.05/1.06 file")
	flag.Parse()

	cls()

	life := NewLife(22, 78)

	if *file == "" {
		life.rPentomino()
	} else {
		life.loadFrom(*file)
	}

	for {
		life.draw()
		life.evolve()
		time.Sleep(125 * time.Millisecond)
	}
}
