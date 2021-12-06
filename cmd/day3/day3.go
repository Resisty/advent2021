package main

import (
    "bufio"
    "bytes"
    "fmt"
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

type commonBinary struct {
    lines []string
    zeroes []int
    ones []int
}

func (cb commonBinary) String() string {
    return fmt.Sprintf("Zeroes: %d, Ones: %d", cb.zeroes, cb.ones)
}

func (cb *commonBinary) CountCommon() {
    // This should only be called once per set of lines, so start fresh
    cb.zeroes = nil
    cb.ones = nil
    for _, line := range(cb.lines) {
        for index, char := range(line) {
            if len(cb.zeroes) <= index {
                cb.zeroes = append(cb.zeroes, 0)
                cb.ones = append(cb.ones, 0)
            }
            switch string(char) {
            case "0":
                cb.zeroes[index] += 1
            case "1":
                cb.ones[index] += 1
            }
        }
    }
}

func (cb *commonBinary) Reduce(index int, uncommon bool) {
    // logger.Logs.Infof("Checking index %d for collection of lines: %s", index, cb.lines)
    // Determine common digit for all lines
    cb.CountCommon()

    // Determine which lines to keep with bit rating
    var result []string
    for _, line := range cb.lines {
        if uncommon {
            // CO2 Rating == keep lines with less common digit OR 0 at index
            if string(line[index]) == cb.UncommonAtIndex(index) {
                result = append(result, line)
            }
        } else {
            // O2 Rating == keep lines with common digit OR 1 at index 
            if string(line[index]) == cb.CommonAtIndex(index) {
                result = append(result, line)
            }
        }
    }
    // logger.Logs.Infof("Reduced collection of lines to %s", result)
    cb.lines = result
}

func (cb *commonBinary) Commonality() (string, string) {
    more, less := "", ""
    for index := range cb.zeroes {
        if cb.zeroes[index] >= cb.ones[index] {
            more += "0"
            less += "1"
        } else {
            more += "1"
            less += "0"
        }
    }
    return more, less
}

func (cb *commonBinary) CommonAtIndex(index int) string {
    if cb.zeroes[index] > cb.ones[index] {
        return "0"
    } else {
        return "1"
    }
}

func (cb *commonBinary) UncommonAtIndex(index int) string {
    if cb.zeroes[index] <= cb.ones[index] {
        return "0"
    } else {
        return "1"
    }
}

func part1() int {
    input := bytes.NewBuffer(reader.FromFile("p1"))
    var lines []string
    scanner := bufio.NewScanner(input)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    cb := commonBinary{lines: lines}
    cb.CountCommon()
    more, less := cb.Commonality()
    logger.Logs.Infof("Gamma rate: %s, epsilon rate: %s", more, less)
    moreInt, _ := strconv.ParseInt(more, 2, 64)
    lessInt, _ := strconv.ParseInt(less, 2, 64)
    return int(moreInt) * int(lessInt)
}

func part2() int {
    input := bytes.NewBuffer(reader.FromFile("p2"))
    var lines []string
    scanner := bufio.NewScanner(input)
    for scanner.Scan() {
        line := scanner.Text()
        lines = append(lines, line)
    }

    // Get O2 rating
    cb := commonBinary{lines: lines}
    index := 0
    cb.Reduce(index, false)
    for ok := true; ok; ok = len(cb.lines) > 1 {
        index += 1
        cb.Reduce(index, false)
    }
    o2RatingInt, _ := strconv.ParseInt(cb.lines[0], 2, 64)

    // Get CO2 rating
    index = 0
    cb = commonBinary{lines: lines}
    cb.Reduce(index, true)
    for ok := true; ok; ok = len(cb.lines) > 1 {
        index += 1
        cb.Reduce(index, true)
    }
    co2RatingInt, _ := strconv.ParseInt(cb.lines[0], 2, 64)

    // Multiply result
    logger.Logs.Infof("O2 Rating: %d, CO2 Rating: %d", o2RatingInt, co2RatingInt)
    return int(o2RatingInt) * int(co2RatingInt)
}
