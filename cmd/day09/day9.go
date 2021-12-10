package main

import (
    "fmt"
    "sort"
    "strconv"
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
    tracked map[Point]struct{}
    maxHeight int
    numMaxHeight map[int]int
}

func newBoard() *Board {
    heights := make(map[Point]int)
    tracked := make(map[Point]struct{})
    numMaxHeight := make(map[int]int)
    board := Board{heights: heights, tracked: tracked, maxHeight: 0, numMaxHeight: numMaxHeight}
    return &board
}

func boardFromInput(lines []string) *Board {
    board := newBoard()
    for i, line := range lines {
        for j, char := range line {
            point := Point{i, j}
            height, _ := strconv.Atoi(string(char))
            board.heights[point] = height
        }
    }
    return board
}

func (b *Board) IsLowPoint(point Point) bool {
    left := Point{point.x - 1, point.y}
    if _, ok := b.heights[left]; ok {
        if b.heights[left] <= b.heights[point] {
            //logger.Logs.Infof("Not low: Point %v (height %d) is at a greater than or equal height to comparison point %v (height %d)", point, b.heights[point], left, b.heights[left])
            return false
        }
    }
    up := Point{point.x, point.y - 1}
    if _, ok := b.heights[up]; ok {
        if b.heights[up] <= b.heights[point] {
            //logger.Logs.Infof("Not low: Point %v (height %d) is at a greater than or equal height to comparison point %v (height %d)", point, b.heights[point], up, b.heights[up])
            return false
        }
    }
    right := Point{point.x + 1, point.y}
    if _, ok := b.heights[right]; ok {
        if b.heights[right] <= b.heights[point] {
            //logger.Logs.Infof("Not low: Point %v (height %d) is at a greater than or equal height to comparison point %v (height %d)", point, b.heights[point], right, b.heights[right])
            return false
        }
    }
    down := Point{point.x, point.y + 1}
    if _, ok := b.heights[down]; ok {
        if b.heights[down] <= b.heights[point] {
            //logger.Logs.Infof("Not low: Point %v (height %d) is at a greater than or equal height to comparison point %v (height %d)", point, b.heights[point], down, b.heights[down])
            return false
        }
    }
    //logger.Logs.Infof("LOW: Point %v (height %d) is at a lesser height than all 4 comparison points", point, b.heights[point])
    return true
}

func (b *Board) ClearTracked() {
    b.tracked = make(map[Point]struct{})
}

func (b *Board) SearchBasinFrom(low Point) int {
    // logger.Logs.Infof("Starting new basin search at point %v", low)
    b.ClearTracked()
    b.tracked[low] = struct{}{}
    return b.SearchBasin(low, 1)
}

func (b *Board) SearchBasin(curr Point, accum int) int {
    // logger.Logs.Infof("Searching basin at point %v", curr)
    b.tracked[curr] = struct{}{}
    if b.heights[curr] == 9 {
        accum -= 1 // don't count peaks
        // logger.Logs.Infof("Stop recursing: point %v (height %d) is a peak. Accumulator at %d", curr, b.heights[curr], accum)
        return accum
    }
    left := Point{curr.x, curr.y - 1}
    if _, ok := b.heights[left]; ok {
        if _, ok := b.tracked[left]; ! ok {
            // logger.Logs.Infof("Recurse left from point %v: Point %v (height %d) exists and hasn't been visited yet", curr, left, b.heights[left])
            accum = b.SearchBasin(left, accum + 1)
            // logger.Logs.Infof("Return from recurse left at point %v, accumulator at %d", curr, accum)
        } /* else {
            logger.Logs.Infof("Skip point %v - already tracked.", left)
        } */
    } /* else {
        // logger.Logs.Infof("Skip point %v - no height.", left)
    } */
    up := Point{curr.x - 1, curr.y}
    if _, ok := b.heights[up]; ok {
        if _, ok := b.tracked[up]; ! ok {
            // logger.Logs.Infof("Recurse up from point %v: Point %v (height %d) exists and hasn't been visited yet", curr, up, b.heights[up])
            accum = b.SearchBasin(up, accum + 1)
            // logger.Logs.Infof("Return from recurse up at point %v, accumulator at %d", curr, accum)
        } /* else {
            // logger.Logs.Infof("Skip point %v - already tracked.", up)
        } */
    } /* else {
        // logger.Logs.Infof("Skip point %v - no height.", up)
    } */
    right := Point{curr.x, curr.y + 1}
    if _, ok := b.heights[right]; ok {
        if _, ok := b.tracked[right]; ! ok {
            // logger.Logs.Infof("Recurse right from point %v: Point %v (height %d) exists and hasn't been visited yet", curr, right, b.heights[right])
            accum = b.SearchBasin(right, accum + 1)
            // logger.Logs.Infof("Return from recurse right at point %v, accumulator at %d", curr, accum)
        } /* else {
            // logger.Logs.Infof("Skip point %v - already tracked.", right)
        } */
    } /* else {
        // logger.Logs.Infof("Skip point %v - no height.", right)
    } */
    down := Point{curr.x + 1, curr.y}
    if _, ok := b.heights[down]; ok {
        if _, ok := b.tracked[down]; ! ok {
            // logger.Logs.Infof("Recurse down from point %v: Point %v (height %d) exists and hasn't been visited yet", curr, down, b.heights[down])
            accum = b.SearchBasin(down, accum + 1)
            // logger.Logs.Infof("Return from recurse down at point %v, accumulator at %d", curr, accum)
        } /* else {
            // logger.Logs.Infof("Skip point %v - already tracked.", down)
        } */
    } /* else {
        // logger.Logs.Infof("Skip point %v - no height.", down)
    } */
    return accum
}

func part1() int {
    lines := reader.LinesFromFile("input.txt")
    board := boardFromInput(lines)
    sum := 0
    for point := range board.heights {
        if board.IsLowPoint(point) {
            sum += board.heights[point] + 1
        }
    }
    return sum
}

func part2() int {
    lines := reader.LinesFromFile("input.txt")
    board := boardFromInput(lines)
    basins := make([]int, 0)
    for point := range board.heights {
        if board.IsLowPoint(point) {
            basin := board.SearchBasinFrom(point)
            // logger.Logs.Infof("Searched basin at point %v, got area %d", point, basin)
            basins = append(basins, basin)
        }
    }
    // logger.Logs.Infof("Collected all basins info: %v", basins)
    sort.Sort(sort.Reverse(sort.IntSlice(basins)))
    return basins[0] * basins[1] * basins[2]
}
