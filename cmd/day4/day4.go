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

const BingoLen = 5

func main() {
	result := part1()
	logger.Logs.Infof("Part one result: %d", result)
	result = part2()
	logger.Logs.Infof("Part two result: %d", result)
}

type Point struct {
	value int
	marked bool
}

func (p *Point) Mark() {
	p.marked = true
	// logger.Logs.Infof("Marked point with value %d to true: %t", p.value, p.marked)
}

func (p Point) String() string {
	return fmt.Sprintf("Value: %d, Marked: %t", p.value, p.marked)
}

type Board struct {
	points [][]*Point
	valuePoints map[int]*Point
}

func (b Board) String() string {
	var lines []string
	for _, row := range b.points {
		var points []string
		for _, p :=  range row {
			points = append(points, p.String())
		}
		lines = append(lines, strings.Join(points, " "))
	}
	s := fmt.Sprint(strings.Join(lines, "\n"))
	return s
}
		

func (b *Board) AddRow(values []int) {
	var row []*Point
	// logger.Logs.Infof("Adding row to board, values: %d", values)
	for _, value := range values {
		point := Point{value: value, marked: false}
		row = append(row, &point)
		b.valuePoints[value] = &point
	}
	b.points = append(b.points, row)
}

func (b *Board) MarkPoint(value int) {
	if point, ok := b.valuePoints[value]; ok {
		// logger.Logs.Infof("Value %d exists in board; marking point.", value)
		point.Mark()
		// logger.Logs.Infof("Updated board: %s", b)
	}
}

func (b *Board) CheckRun() bool {
	for i := range b.points {
		got_run := 0
		for j := range b.points {
			if ! b.points[i][j].marked {
				got_run =0
				break
			} else {
				got_run += 1
			}
		}
		if got_run == BingoLen {
			return true
		}
	}
	for i := range b.points {
		got_run := 0
		for j := range b.points {
			if ! b.points[j][i].marked {
				got_run = 0
				break
			} else {
				got_run += 1
			}
		}
		if got_run == BingoLen {
			return true
		}
	}
	return false
}

func (b *Board) SumEmpties() int {
	sum := 0
	for i := range b.points {
		for j := range b.points {
			if ! b.points[i][j].marked {
				sum += b.points[i][j].value
			}
		}
	}
	return sum
}


func stringsToInts (numbersString string, separator string) []int {
	var numbers []int
	numberStrings := regexp.MustCompile(separator).Split(strings.TrimSpace(numbersString), -1)
	for _, numberString := range numberStrings {
		number, _ := strconv.Atoi(strings.TrimSpace(numberString))
		numbers = append(numbers, number)
	}
	return numbers
}

func gameFromInput(lines []string) ([]int, []Board) {
	bingoCalls := stringsToInts(lines[0], ",")
	lines = lines[1:]
	var boards []Board
	board := Board{make([][]*Point, 0), make(map[int]*Point)}
	count := 0
	for i := 0; i < len(lines); i++ {
		if lines[i] == "" {
			continue
		} else {
			board.AddRow(stringsToInts(lines[i], `\s+`))
			count += 1
		}
		if count == BingoLen {
			boards = append(boards, board)
			board = Board{make([][]*Point, 0), make(map[int]*Point)}
			count = 0
		}
	}	
	return bingoCalls, boards
}


func part1() int {
	input := bytes.NewBuffer(reader.FromFile("p1"))
    scanner := bufio.NewScanner(input)
	var lines []string
    for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	bingoCalls, boards := gameFromInput(lines)
	// logger.Logs.Infof("Bingo numbers to call: %d", bingoCalls)
	// logger.Logs.Infof("%d Boards to play", len(boards))
	// for i := range boards {
	// 	logger.Logs.Infof("Board %d: %s", i, boards[i])
	// }
	for _, number := range bingoCalls {
		for _, board := range boards {
			board.MarkPoint(number)
			if board.CheckRun() {
				emptySum := board.SumEmpties()
				return emptySum * number
			}
		}
	}
	logger.Logs.Infof("You screwed up!")
	return 4
}

func part2() int {
	input := bytes.NewBuffer(reader.FromFile("p2"))
    scanner := bufio.NewScanner(input)
	var lines []string
    for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	bingoCalls, boards := gameFromInput(lines)
	// logger.Logs.Infof("Bingo numbers to call: %d", bingoCalls)
	// logger.Logs.Infof("%d Boards to play", len(boards))
	// for i := range boards {
	// 	logger.Logs.Infof("Board %d: %s", i, boards[i])
	// }
	lastScore := -1
	winningBoards := make(map[int]bool)
	for i := range boards {
		winningBoards[i] = false
	}
	for _, number := range bingoCalls {
		for j, board := range boards {
			board.MarkPoint(number)
			if board.CheckRun() && ! winningBoards[j] {
				winningBoards[j] = true
				emptySum := board.SumEmpties()
				lastScore = emptySum * number
			}
		}
	}
	return lastScore
}
