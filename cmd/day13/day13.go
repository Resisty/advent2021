package main

import (
    "fmt"
    "regexp"
    "strconv"
    logger "advent2021/adventlogger"
    reader "advent2021/adventreader"
)

type Point struct {
    x, y int
}

type Board struct {
    points map[Point]struct{}
    height, width int
}

type Instruct struct {
    axis string
    coord int
}

func (p Point) String() string {
    return fmt.Sprintf("%d,%d", p.x, p.y)
}

func NewBoard() *Board {
    points := make(map[Point]struct{})
    board := Board{points: points, height: 0, width: 0}
    return &board
}

func (b *Board) Print() {
    for i := 0; i < b.height; i++ {
        line := ""
        for j := 0; j < b.width; j++ {
            char := "."
            if _, ok := b.points[Point{i, j}]; ok {
                char = "#"
            }
            line += char
        }
        fmt.Println(line)
    }
}

func (b *Board) AddLine(line string) {
    yx := regexp.MustCompile(",").Split(line, -1)
    y, _ := strconv.Atoi(string(yx[0]))
    x, _ := strconv.Atoi(string(yx[1]))
    point := Point{x, y}
    b.points[point] = struct{}{}
    if y >= b.width {
        b.width = y + 1
    }
    if x >= b.height {
        b.height = x + 1
    }
}

func boardFromInput(lines []string) (*Board, []Instruct) {
    board := NewBoard()
    coordsDone := false
    instructs := make([]Instruct, 0)
    for _, line := range lines {
        if line == "" {
            coordsDone = true
            continue
        }
        if ! coordsDone {
            board.AddLine(line)
        } else {
            tokens := regexp.MustCompile(" ").Split(line, -1)
            axisCoord := regexp.MustCompile("=").Split(tokens[len(tokens) - 1], -1)
            coord, _ := strconv.Atoi(axisCoord[1])
            instructs = append(instructs, Instruct{axis: axisCoord[0], coord: coord})
        }
    }
    return board, instructs
}

func (b *Board) FoldVert(x int) {
    for point := range b.points {
        if point.x > x {
            b.points[Point{2 * x - point.x,point.y}] = struct{}{}
            delete(b.points, point)
        }
    }
    b.height /= 2
}

func (b *Board) FoldHoriz(y int) {
    for point := range b.points {
        if point.y > y {
            b.points[Point{point.x, 2 * y - point.y}] = struct{}{}
            delete(b.points, point)
        }
    }
    b.width /= 2
}

func (b *Board) DoFold(ins Instruct) {
    decision := map[string]func(*Board, int){
        "y": (*Board).FoldVert,
        "x": (*Board).FoldHoriz,
    }
    decision[ins.axis](b, ins.coord)
}

func main() {
    result := part1()
    logger.Logs.Infof("Part one result: %d", result)
    result = part2()
    logger.Logs.Infof("Part two result: %d", result)
}

func part1() int {
    lines := reader.LinesFromFile("input.txt")
    board, instructions := boardFromInput(lines)
    board.DoFold(instructions[0])
    return len(board.points)
}

func part2() int {
    lines := reader.LinesFromFile("input.txt")
    board, instructions := boardFromInput(lines)
    for _, ins := range instructions {
        board.DoFold(ins)
    }
    board.Print()
    return 4
}
