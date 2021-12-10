package main

import (
    "fmt"
    "sort"
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

var Openers = map[string]string{
    "[": "]",
    "(": ")",
    "{": "}",
    "<": ">",
}
var Closers = map[string]string{
    "]": "[",
    ")": "(",
    "}": "{",
    ">": "<",
}

var CloserPoints = map[string]int{
    ")": 3,
    "]": 57,
    "}": 1197,
    ">": 25137,
}

var CollapserPoints = map[string]int{
    ")": 1,
    "]": 2,
    "}": 3,
    ">": 4,
}

func IsCloser(char string) bool {
    _, ok := Closers[char]
    return ok
}

type Stack []string

func (s Stack) String() string {
    return strings.Join(s, "")
}

func (s *Stack) IsEmpty() bool {
    return len(*s) == 0
}

func (s *Stack) Push(char string) error {
    last := s.Tail()
    if IsCloser(char) && last == Closers[char] {
        s.Pop()
        return nil
    }
    if IsCloser(char) {
        return fmt.Errorf("Trying to close chunk (opener: '%s') with closer '%s'! Illegal!", last, char) 
    }
    *s = append(*s, char)
    return nil
}

func (s *Stack) Pop() (string, bool) {
    if s.IsEmpty() {
        return "", false
    } else {
        index := len(*s) - 1
        item := (*s)[index]
        *s = (*s)[:index]
        return item, true
    }
}

func (s *Stack) Tail() string {
    if s.IsEmpty() {
        return ""
    }
    return (*s)[len(*s) - 1]
}

func part1() int {
    lines := reader.LinesFromFile("input.txt")
    points := 0
    for _, line := range lines {
        s := make(Stack, 0)
        for _, char := range line {
            if err := s.Push(string(char)); err != nil {
                points += CloserPoints[string(char)]
                break
            }
        }
    }
    return points
}

func part2() int {
    lines := reader.LinesFromFile("input.txt")
    scores := make([]int, 0)
    points := 0
    LINES:
    for _, line := range lines {
        points = 0
        s := make(Stack, 0)
        for _, char := range line {
            if err := s.Push(string(char)); err != nil {
                // line is corrupt, move to next line
                continue LINES
            }
        }
        // complete the line
        for ! s.IsEmpty() {
            last := s.Tail()
            collapser := Openers[last]
            if err := s.Push(collapser); err != nil {
                panic(err)
            }
            points += (4 * points) + CollapserPoints[collapser]
        }
        scores = append(scores, points)
    }
    sort.Ints(scores)
    return scores[len(scores) / 2]
}
