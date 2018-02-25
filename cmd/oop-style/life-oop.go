package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// GIF parameters
var palette = []color.Color{color.White, color.Black}
var gifCellSize = 8
var gifCellPadding = 2

// Config struct holds CLI parameters values
var Config struct {
	file       string
	deltaXflag int
	deltaYflag int
	nrows      int
	ncols      int
	turns      int
	rate       int
	ageColor   bool
	ageShape   bool
	genGIF     string
}

// end of CLI parameters

func cls() {
	fmt.Printf("\033[2J")
}

func pos(r, c int) {
	fmt.Printf("\033[" + strconv.Itoa(r) + ";" + strconv.Itoa(c) + "H")
}

// Universe struct
type Universe struct {
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
func NewLife(r, c int) *Universe {

	u := &Universe{}

	u.rows = r
	u.cols = c

	u.board = newBoard(r, c)
	u.prev = newBoard(r, c)

	u.alive = 0
	u.gen = 1

	return u
}

func (u *Universe) rPentomino() {
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
	// can't use string slice and index, because of UTF chars
	if Config.ageShape && age < 4 {
		switch age {
		case 1:
			shape = "."
		case 2:
			shape = "∘"
		case 3:
			shape = "∙"
		}
	}

	if Config.ageColor {
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

// paints cell at specified (row, column) location
func drawCell(img *image.Paletted, r, c int) {
	for y := 0; y < gifCellSize; y++ {
		for x := 0; x < gifCellSize; x++ {
			img.SetColorIndex(c*(gifCellSize+gifCellPadding)+x, r*(gifCellSize+gifCellPadding)+y, 1)
		}
	}
}

func addLabel(img *image.Paletted, x, y int, label string) {

	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}

	d := &font.Drawer{
		Dst: img,
		Src: image.NewUniform(palette[1]),
		//Face: inconsolata.Regular8x16,
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

// return gif image of current generation
func (u *Universe) image() *image.Paletted {
	rect := image.Rect(0, 0, u.cols*(gifCellSize+gifCellPadding), u.rows*(gifCellSize+gifCellPadding))
	img := image.NewPaletted(rect, palette)

	for i := 0; i < u.cols*(gifCellSize+gifCellPadding); i++ {
		img.SetColorIndex(i, 0, 1)
		img.SetColorIndex(i, u.rows*(gifCellSize+gifCellPadding)-1, 1)
	}
	for i := 0; i < u.rows*(gifCellSize+gifCellPadding); i++ {
		img.SetColorIndex(0, i, 1)
		img.SetColorIndex(u.cols*(gifCellSize+gifCellPadding)-1, i, 1)
	}

	for r := 0; r < u.rows; r++ {
		for c := 0; c < u.cols; c++ {
			if u.board[r][c] > 0 {
				drawCell(img, r, c)
			}
		}
	}

	l := fmt.Sprintf("Conway's Life in Go | board %dx%d;", u.rows, u.cols)
	l = l + fmt.Sprintf(" rate %d/sec; alive = %3d; gen = %d/%d", Config.rate, u.alive, u.gen, Config.turns)
	addLabel(img, 10, 20, l)

	return img
}

func (u *Universe) draw() {
	pos(0, 0)

	fmt.Printf("Conway's Life in Go | board %dx%d;", u.rows, u.cols)
	fmt.Printf(" rate %d/sec; alive = %3d; gen = %d\033[K",
		Config.rate, u.alive, u.gen)

	// to speed up we're going to update the frame only every 100 generation
	redrawFrame := u.gen%100 == 1
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

// returns number of alive cells surrounding the given one
func (u *Universe) neighbours(r, c int) (count int) {
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

// evolves the universe 1 step according to standard B3/S23
func (u *Universe) evolve() {
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

// shape loading from LIF 1.05
func (u *Universe) loadShape105(lines []string) {
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

// shape loading from LIF 1.06
func (u *Universe) loadShape106(lines []string) {
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

// function returns top left and bottom right coords of a rectangle
// that enframes the shape. Used to translate the shape to new position.
func (u *Universe) boundries() (startX, startY, endX, endY int) {
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

// translate moves the shape to a new position.
// if no CLI parameters specified, the shape is centered
func (u *Universe) translate() {
	startX, startY, endX, endY := u.boundries()

	deltaX := Config.deltaXflag
	deltaY := Config.deltaYflag

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

// load shape from a file
func (u *Universe) loadFrom(fn string) {
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

// CLI parameters parsing
func init() {
	flag.StringVar(&Config.file, "f", "", "Load life pattern from LIF 1.05/1.06 file")

	flag.IntVar(&Config.deltaXflag, "deltax", 0, "X translation for loaded shape (default: auto)")
	flag.IntVar(&Config.deltaYflag, "deltay", 0, "Y translation for loaded shape (default: auto)")

	flag.IntVar(&Config.nrows, "rows", 22, "number of rows")
	flag.IntVar(&Config.ncols, "cols", 78, "number of cols")

	flag.IntVar(&Config.turns, "turns", 0, "number of generations to simulate")
	flag.IntVar(&Config.rate, "r", 2, "Rate of generations per second")

	flag.BoolVar(&Config.ageColor, "color", false, "use color to show cell age")
	flag.BoolVar(&Config.ageShape, "shape", false, "use shapes to show cell age")

	flag.StringVar(&Config.genGIF, "gif", "", "generate GIF file with evolution")
}

// main cycle
func main() {

	flag.Parse()

	anim := gif.GIF{LoopCount: Config.turns}
	var outgif *os.File
	if Config.genGIF != "" {
		if Config.turns == 0 {
			log.Fatal("option -gif requires number of generations to simulate (-turns)")
		}
		var err error
		outgif, err = os.Create(Config.genGIF)
		if err != nil {
			log.Fatalf("failed to create file: %v", err)
		}
	} else {
		cls()
	}

	life := NewLife(Config.nrows, Config.ncols)
	if Config.file == "" {
		life.rPentomino()
	} else {
		life.loadFrom(Config.file)
	}

	for {
		if Config.genGIF != "" {
			anim.Image = append(anim.Image, life.image())
			if Config.turns > 0 && life.gen >= Config.turns {
				anim.Delay = append(anim.Delay, 300+100/Config.rate)
				fmt.Printf("\nReached generation %v\n", life.gen)
				break
			} else {
				anim.Delay = append(anim.Delay, 100/Config.rate)
			}
			life.evolve()
		} else {
			life.draw()
			if Config.turns > 0 && life.gen >= Config.turns {
				fmt.Printf("\nReached generation %v\n", life.gen)
				break
			}
			life.evolve()
			time.Sleep(time.Second / time.Duration(Config.rate))
		}
	}

	if Config.genGIF != "" {
		fmt.Printf("Generating GIF ... ")
		gif.EncodeAll(outgif, &anim)
		fmt.Printf("done.\n[%v]\n", Config.genGIF)
	}
}
