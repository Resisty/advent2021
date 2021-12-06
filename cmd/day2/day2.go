package main

import (
    "bufio"
    "bytes"
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

func part1() int {
    x_pos, y_pos := 0, 0
    input := bytes.NewBuffer(reader.FromFile("p1"))
    scanner := bufio.NewScanner(input)
    for scanner.Scan() {
        tokens := strings.Split(scanner.Text(), " ")
        switch direction := tokens[0]; direction {
        case "forward":
            unit, _ := strconv.Atoi(tokens[1])
            x_pos += unit
        case "down":
            unit, _ := strconv.Atoi(tokens[1])
            y_pos += unit
        case "up":
            unit, _ := strconv.Atoi(tokens[1])
            y_pos -= unit
        }
    }
    return x_pos * y_pos
}

func part2() int {
    x_pos, y_pos, aim := 0, 0, 0
    input := bytes.NewBuffer(reader.FromFile("p2"))
    scanner := bufio.NewScanner(input)
    for scanner.Scan() {
        tokens := strings.Split(scanner.Text(), " ")
        switch direction := tokens[0]; direction {
        case "forward":
            unit, _ := strconv.Atoi(tokens[1])
            x_pos += unit
            y_pos += (aim * unit)
        case "down":
            unit, _ := strconv.Atoi(tokens[1])
            aim += unit
        case "up":
            unit, _ := strconv.Atoi(tokens[1])
            aim -= unit
        }
    }
    return x_pos * y_pos
}
