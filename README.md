# Conway's Life simulation in Go

Small and simple Conway's Life simpulator project exploring Go languague. Three different implementations with different features:

- Quick and simple with approx 100 lines of code
- More features with old-facioned classical procedural approach
- Most feature rich with interface-style approach (number of features have nothing to do with approaches :) )

Examples below use oop-style version.

## Build

```
git clone ...
cd golang-life/cmd/oop-style
go build -ldfages="-s -w" life-oop.go -o life
```

## Usage

```
$ ./life -h
Conway's Life simulator in Go:
  Text-mode or animated gif simulation according to B2/S23 rules.

Usage:
  -color
    	In text-mode use color to show cell age
  -cols int
    	Number of cols (default 78)
  -deltax int
    	X translation for loaded shape (centered by default)
  -deltay int
    	Y translation for loaded shape (centered by default)
  -f string
    	Load life pattern from LIF 1.05/1.06 file
  -gif string
    	Instead of text-mode, generate animated GIF file with
    	specified name containing the evolution.
  -r int
    	Rate of generations per second (default 2)
  -rows int
    	Number of rows (default 22)
  -shape
    	In text-mode use shapes to show cell age
  -turns int
    	Number of generations to simulate
```

## Text-mode demonstation

[![asciicast](https://asciinema.org/a/UtQosxric9ff5DggaombI0zBW.png)](https://asciinema.org/a/UtQosxric9ff5DggaombI0zBW)

## GIF demonstation

![](images/piheptomino.gif)