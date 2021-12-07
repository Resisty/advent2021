package main

import (
    "regexp"
    "strconv"
    logger "advent2021/adventlogger"
    reader "advent2021/adventreader"
)

const DaysToAnalyze = 80
const AnglerCycle = 8
const AnglerRefractory = 2

func main() {
    result := part1()
    logger.Logs.Infof("Part one result: %d", result)
    result = part2()
    logger.Logs.Infof("Part two result: %d", result)
}

func anglersTick(anglers map[int]int) map[int]int {
    replace := make(map[int]int)
    for dayNum := 0; dayNum < AnglerCycle + 1; dayNum++ {
        numAnglers := anglers[dayNum]
        if dayNum == 0 {
            replace[AnglerCycle - AnglerRefractory] = numAnglers
            replace[AnglerCycle] = numAnglers
            // logger.Logs.Infof("Anglers reproduced and reset (%d): %d. Also set %d new anglers to %d (max timer)", AnglerCycle -  AnglerRefractory, numAnglers, numAnglers, AnglerCycle)
        } else {
            replace[dayNum - 1] += numAnglers
            // logger.Logs.Infof("Anglers at timer %d ticked down to %d: %d", dayNum, dayNum - 1, anglers[dayNum])
        }
    }
    return replace
}

func makeAnglersFromLines(lines []string) map[int]int {
    anglers := make(map[int]int)
    for _, line := range lines {
        intList := regexp.MustCompile(",").Split(line, -1)
        for _, intString := range intList {
            integer, _ := strconv.Atoi(intString)
            anglers[integer] += 1
        }
    }
    return anglers
}

func part1() int {
    lines := reader.LinesFromFile("input.txt")
    anglers := makeAnglersFromLines(lines)
    logger.Logs.Infof("Anglers: %d", anglers)
    for i := 0; i < DaysToAnalyze; i++ {
        anglers = anglersTick(anglers)
        // logger.Logs.Infof("Anglers after day %d: %d", i + 1, anglers)
    }
    sum := 0
    for _, numAnglers := range anglers {
        sum += numAnglers
    }
    return sum
}

func part2() int {
    lines := reader.LinesFromFile("input.txt")
    anglers := makeAnglersFromLines(lines)
    logger.Logs.Infof("Anglers: %d", anglers)
    for i := 0; i < 256; i++ {
        anglers = anglersTick(anglers)
        // logger.Logs.Infof("Anglers after day %d: %d", i + 1, anglers)
    }
    sum := 0
    for _, numAnglers := range anglers {
        sum += numAnglers
    }
    return sum
}

