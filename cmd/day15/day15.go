package main

import (
    "fmt"
    "sort"
    "strconv"
    logger "advent2021/adventlogger"
    reader "advent2021/adventreader"
)

const MaxInt = int(^uint(0) >> 1)

type Point struct {
    x, y int
}

type CostPoint struct {
    Point
    cost int
    dist int
}

type CostPoints []*CostPoint

func (costs CostPoints) Len() int {
    return len(costs)
}

func (costs CostPoints) Swap(i, j int) {
    costs[i], costs[j] = costs[j], costs[i]
}

func (costs CostPoints) Less(i, j int) bool {
    if costs[i].dist == costs[j].dist {
        if costs[i].x != costs[j].x {
            return costs[i].y < costs[j].y
        } else {
            return costs[i].x < costs[j].x
        }
    }
    return costs[i].dist < costs[j].dist
}

type Board struct {
    points map[Point]*CostPoint
    height, width int
}

func (p Point) String() string {
    return fmt.Sprintf("%d,%d", p.x, p.y)
}

func (c CostPoint) String() string {
    return fmt.Sprintf("(%d,%d):%d", c.x, c.y, c.dist)
}


func NewBoard(h, w int) *Board {
    points := make(map[Point]*CostPoint)
    board := Board{points: points, height: h, width: w}
    return &board
}

func (b *Board) KayakDotCom() int {
    dx := []int{-1, 0, 1, 0}
    dy := []int{0, 1, 0, -1}
    costList := make([]*CostPoint, 0)
    costList = append(costList, b.points[Point{0, 0}])
    for len(costList) != 0 {
        curr := costList[0]
        costList = costList[1:]
        for i := 0; i < 4; i++ {
            x := curr.x + dx[i]
            y := curr.y + dy[i]
            point, ok := b.points[Point{x, y}]
            if ! ok {
                continue
            }
            if point.dist > curr.dist + point.cost {
                point.dist = curr.dist + point.cost
                costList = append(costList, point)
            }
        }
        sort.Sort(CostPoints(costList))
    }
    return b.points[Point{b.width - 1, b.height - 1}].dist
}

func boardsFromInput(lines []string) *Board {
    height := len(lines)
    width := len(lines[0])
    board := NewBoard(height, width)
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            point := Point{x, y}
            dist := MaxInt
            cost, _ := strconv.Atoi(string(lines[x][y]))
            if x == 0 && y == 0 {
                dist = 0 // origin has no cost to enter and therefore no dist
            }
            board.points[point] = &CostPoint{point, cost, dist}
        }
    }
    return board
}

func (b *Board) Embiggen() *Board {
    bigBoard := NewBoard(b.height * 5, b.width * 5)
    count := 0
    for point, costPoint := range b.points {
        for y := 0; y < 5; y++ {
            for x := 0; x < 5; x++ {
                translate := Point{x * b.width + point.x, y * b.height + point.y}
                newCost := costPoint.cost + x + y
                dist := MaxInt
                if translate.x == 0 && translate.y == 0 {
                    dist = 0
                }
                if newCost >= 10 {
                    newCost = newCost % 10 + 1
                }
                newCostPoint := &CostPoint{translate, newCost, dist}
                bigBoard.points[translate] = newCostPoint
                count += 1
            }
        }
    }
    return bigBoard
}

func main() {
    result := part1()
    logger.Logs.Infof("Part one result: %d", result)
    result = part2()
    logger.Logs.Infof("Part two result: %d", result)
}

func part1() int {
    lines := reader.LinesFromFile("input.txt")
    board := boardsFromInput(lines)
    return board.KayakDotCom()
}

func part2() int {
    lines := reader.LinesFromFile("input.txt")
    board := boardsFromInput(lines)
    board = board.Embiggen()
    return board.KayakDotCom()
}
