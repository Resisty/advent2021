package main

import (
    "fmt"
    "strconv"
    "strings"
    logger "advent2021/adventlogger"
    reader "advent2021/adventreader"
)

const binStringLength = 9

type Point struct {
    x, y int
}

type Board struct {
    points map[Point]string
    height, width int
    infiniteFill, enhance string // [ types furiously ] "Enhance."
}

func (p Point) String() string {
    return fmt.Sprintf("%d,%d", p.x, p.y)
}

func NewBoard(h, w int, enhance string) *Board {
    points := make(map[Point]string)
    board := Board{points: points, height: h, width: w, infiniteFill: ".", enhance: enhance}
    return &board
}

func (b *Board) Print() {
    for i := 0; i < b.height; i++ {
        line := ""
        for j := 0; j < b.width; j++ {
            if char, ok := b.points[Point{j, i}]; ok {
                line += char
            }
        }
        fmt.Println(line)
    }
}

func (b *Board) expand() {
    points := make(map[Point]string)
    h, w := b.height + 2, b.width + 2
    for y := 0; y < h; y++ {
        for x := 0; x < w; x++ {
            if x == 0 || y == 0 || x == w - 1 || y == h - 1 {
                // new border
                points[Point{x, y}] = b.infiniteFill
            } else {
                if char, ok := b.points[Point{x - 1, y - 1}]; ok {
                    // shifting previous points to x+1, y+1 means subtract 1 from current coords
                    points[Point{x, y}] = char
                }
            }
        }
    }
    b.points = points
    b.height = h
    b.width = w
}

func lightsToBin(s string) string {
    binString := ""
    for _, char := range s {
        if string(char) == "." {
            binString += "0"
        } else {
            binString += "1"
        }
    }
    return binString
}

func binToInt(s string) int {
    i, _ := strconv.ParseInt(s, 2, 64)
    return int(i)
}

func (b *Board) enhanceAt(x, y int) string {
    movement := []Point{
        Point{-1, -1}, Point{0, -1}, Point{1, -1},
        Point{-1, 0},  Point{0, 0},  Point{1, 0},
        Point{-1, 1},  Point{0, 1},  Point{1, 1},
    }
    lights := ""
    for _, move := range movement {
        stepTo := Point{move.x + x, move.y + y}
        if val, ok := b.points[stepTo]; ok {
            lights += val
        } else {
            lights += b.infiniteFill
        }
    }
    return string(b.enhance[binToInt(lightsToBin(lights))])
}

func (b *Board) enhanceN(num int) {
    for i := 0; i < num; i++{
        b.expand()
        newPoints := make(map[Point]string)
        for point, char := range b.points {
            newPoints[point] = char
        }
        for y := 0; y < b.height; y++ {
            for x := 0; x < b.width; x++ {
                newPoints[Point{x, y}] = b.enhanceAt(x, y)
            }
        }
        b.points = newPoints
        b.infiniteFill = string(b.enhance[binToInt(lightsToBin(strings.Repeat(b.infiniteFill, binStringLength)))])
    }
}

func (b *Board) countLit() int {
    count := 0
    for y := 0; y < b.height; y++ {
        for x := 0; x < b.width; x++ {
            char := b.points[Point{x, y}]
            if char == "#" {
                count += 1
            }
        }
    }
    return count
}

func boardFromInput(enhance string, lines []string) *Board {
    height, width := len(lines), len(lines[0])
    board := NewBoard(height, width, enhance)
    for y, line := range lines {
        for x, char := range line {
            p := Point{x, y}
            board.points[p] = string(char)
        }
    }
    return board
}

func main() {
    result := part1()
    logger.Logs.Infof("Part one result: %d", result)
    result = part2()
    logger.Logs.Infof("Part two result: %d", result)
}

func part1() int {
    lines := reader.LinesFromFile("input.txt")
    enhancement := lines[0]
    image := lines[2:]
    board := boardFromInput(enhancement, image)
    board.enhanceN(2)
    return board.countLit()
}

func part2() int {
    lines := reader.LinesFromFile("input.txt")
    enhancement := lines[0]
    image := lines[2:]
    board := boardFromInput(enhancement, image)
    board.enhanceN(50)
    return board.countLit()
}
