package main

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	logger "advent2021/adventlogger"
	reader "advent2021/adventreader"
)

func main() {
	result := part1()
	logger.Logs.Infof("Part one result: %d", result)
	result = part2()
	logger.Logs.Infof("Part two result: %d", result)
}

type Point struct {
	x, y int
}

func (p Point) String() string {
	return fmt.Sprintf("%d,%d", p.x, p.y)
}

type Board struct {
	heights map[Point]int
	maxHeight int
	numMaxHeight map[int]int
}

// greatest common divisor (GCD) via Euclidean algorithm
// Lovingly stolen from https://siongui.github.io/2017/05/14/go-gcd-via-euclidean-algorithm/
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func slope(x1, y1, x2, y2 int) (int, int) {
	delta_x := x2 - x1
	delta_y := y2 - y1
    absx, absy := delta_x, delta_y
    if delta_x < 0 {
        absx = -delta_x
    }
    if delta_y < 0 {
        absy = -delta_y
    }
    if delta_x == 0 {
        return 0, delta_y / absy
    }
    if delta_y == 0 {
        return delta_x / absx, 0
    }
	gcd := GCD(absx, absy)
	delta_x /= gcd
	delta_y /= gcd
	return delta_x, delta_y
}

func newBoard() *Board {
	heights := make(map[Point]int)
	numMaxHeight := make(map[int]int)
	board := Board{heights: heights, maxHeight: 0, numMaxHeight: numMaxHeight}
	return &board
}

func (b *Board) trackMaxHeight(point Point) {
	height := b.heights[point]
	if height >= b.maxHeight {
		b.maxHeight = height
		b.numMaxHeight[height] += 1
	}
}

func (b *Board) AddHorizVertLine(p1, p2 Point) {
    if p1.x == p2.x || p1.y == p2.y {
        // logger.Logs.Infof("Required horiz|vert line only; satisfied by points %s -> %s", p1, p2)
        b.AddLine(p1, p2)
        return
    }
    // logger.Logs.Infof("Required horiz|vert line only; NOT SATISFIED by points %s -> %s", p1, p2)
}

func (b *Board) AddHVDLine(p1, p2 Point) {
    // Add lines if they are horizontal, vertical, or pi/4 radians.
	dx, dy := slope(p1.x, p1.y, p2.x, p2.y)
    if p1.x == p2.x || p1.y == p2.y || dx == dy || dx + dy == 0 {
        // logger.Logs.Infof("Required horiz|vert line only; satisfied by points %s -> %s", p1, p2)
        b.AddLine(p1, p2)
        return
    }
    // logger.Logs.Infof("Required horiz|vert line only; NOT SATISFIED by points %s -> %s", p1, p2)
}

func (b *Board) AddLine(p1, p2 Point) {
	dx, dy := slope(p1.x, p1.y, p2.x, p2.y)
	i, j := p1.x, p1.y
	for i != p2.x || j != p2.y {
		point := Point{i, j}
        //logger.Logs.Infof("Tracking point %s on line connecting %s to %s", point, p1, p2)
		b.heights[point] += 1
		b.trackMaxHeight(point)
		i += dx
		j += dy
	}
	// Account for the final point
	point := Point{i, j}
    // logger.Logs.Infof("Tracking final point %s on line connecting %s to %s", point, p1, p2)
	b.heights[point] += 1
	b.trackMaxHeight(point)
}

func lineToPoints (lineString string, separator string) []Point {
	tokens := regexp.MustCompile(separator).Split(strings.TrimSpace(lineString), -1)
	points := make([]Point, 0)
	for _, coordString := range tokens {
		xyStringList := regexp.MustCompile(",").Split(coordString, -1)
		x, _ := strconv.Atoi(xyStringList[0])
		y, _ := strconv.Atoi(xyStringList[1])
		points = append(points, Point{x, y})
	}
	return points
}

func boardFromInput(lines []string, choice ...string) *Board {
	board := newBoard()
	for _, line := range lines {
		// logger.Logs.Infof("Adding line to board; line = %s", line)
		points := lineToPoints(line, `\s->\s`)
        if len(choice) > 0 {
            switch linesTypes := choice[0]; linesTypes {
            case "hv": 
                board.AddHorizVertLine(points[0], points[1])
            case "hvd": 
                board.AddHVDLine(points[0], points[1])
            }
        } else {
            board.AddLine(points[0], points[1])
        }
	}
	return board
}

func inputLines(part int) []string {
	input := bytes.NewBuffer(reader.FromFile("p"+fmt.Sprintf("%d", part)))
    scanner := bufio.NewScanner(input)
	var lines []string
    for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}


func part1() int {
	lines := inputLines(2)
	board := boardFromInput(lines, "hv")
    sum := 0
    for _, height := range board.heights {
        if height < 2 {
            continue
        }
        sum += 1
    }
	return sum
}

func part2() int {
	lines := inputLines(2)
	board := boardFromInput(lines, "hvd")
    sum := 0
    for _, height := range board.heights {
        if height < 2 {
            continue
        }
        sum += 1
    }
	return sum
}
