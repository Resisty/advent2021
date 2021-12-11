package main

import (
    "fmt"
    "strconv"
    logger "advent2021/adventlogger"
    reader "advent2021/adventreader"
)

type Point struct {
    x, y int
}

type Octopus struct {
    energy int
    flash bool
}

func (p Point) String() string {
    return fmt.Sprintf("%d,%d", p.x, p.y)
}

func (o Octopus) String() string {
    return fmt.Sprintf("Energy: %b, flashed: %t", o.energy, o.flash)
}

func (o *Octopus) Increment() {
    o.energy += 1
}

func (o *Octopus) Flash() {
    o.flash = true
}

func (o *Octopus) Reset() {
    if o.flash {
        o.energy = 0
    }
    o.flash = false
}

type Board struct {
    octopuses map[Point]*Octopus
    flashes, stepFlashes int
    width, height int
}

func (b *Board) Print() {
    for i := 0; i < b.width; i++ {
        s := ""
        for j := 0; j < b.height; j++ {
            s += strconv.Itoa(b.octopuses[Point{i, j}].energy)
        }
        fmt.Println(s)
    }
}

func newBoard() *Board {
    octopuses := make(map[Point]*Octopus)
    flashes := 0
    board := Board{octopuses: octopuses, flashes: flashes}
    return &board
}

func boardFromInput(lines []string) *Board {
    board := newBoard()
    board.width = len(lines)
    board.height = len(lines)
    for i, line := range lines {
        for j, char := range line {
            energy, _ := strconv.Atoi(string(char))
            point := Point{x: i, y: j}
            octo := &Octopus{energy: energy}
            board.octopuses[point] = octo
        }
    }
    return board
}

func (b *Board) Increment() {
    for _, octopus := range b.octopuses {
        octopus.Increment()
    }
}

func (b *Board) Reset() {
    for _, octopus := range b.octopuses {
        if octopus.energy > 9 {
            octopus.Reset()
        }
    }
}

func (b *Board) Step() bool {
    b.stepFlashes = 0
    b.Increment()
    for point := range b.octopuses {
        b.ResolveFlash(point)
    }
    b.Reset()
    return b.stepFlashes != b.width * b.height
}


func (b *Board) ResolveFlash(p Point) {
    octopus := b.octopuses[p]
    if octopus.energy <= 9 || octopus.flash {
        return
    }
    octopus.Flash()
    b.flashes += 1
    b.stepFlashes += 1
    dirs := [][]int{
        {-1, 0}, {-1, -1}, {0, -1}, {1, -1}, 
        {1, 0}, {1, 1}, {0, 1}, {-1, 1}, 
    }
    for _, dir := range dirs {
        neighborPoint := Point{p.x+dir[0], p.y+dir[1]}
        if neighbor, ok := b.octopuses[neighborPoint]; ok {
            neighbor.energy += 1
            b.ResolveFlash(neighborPoint)
        }
    }
}

func main() {
    result := part1()
    logger.Logs.Infof("Part one result: %d", result)
    result = part2()
    logger.Logs.Infof("Part two result: %d", result)
}

func part1() int {
    lines := reader.LinesFromFile("test.txt")
    board := boardFromInput(lines)
    for i := 0; i < 100 ; i++ {
        board.Step()
    }
    return board.flashes
}

func part2() int {
    lines := reader.LinesFromFile("input.txt")
    board := boardFromInput(lines)
    board.Print()
    i := 0
    for board.Step() {
        i+= 1
    }
    return i + 1 // last call to board.Step() was false so i didn't increment
}
