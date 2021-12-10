package main

import (
    "math"
    "regexp"
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

func linesToInts(lines []string) []int {
    ints := make([]int, 0)
    for _, line := range lines {
        intList := regexp.MustCompile(",").Split(line, -1)
        for _, intString := range intList {
            integer, _ := strconv.Atoi(intString)
            ints = append(ints, integer)
        }
    }
    return ints
}

func avg(a, b int) int {
    return (a + b) / 2
}

func median(ints []int) int {
    sort.Ints(ints)
    if len(ints) % 2 == 1 {
        return ints[len(ints) / 2]
    } else {
        return avg(ints[(len(ints) / 2) - 1], ints[len(ints) / 2])
    }
}

func fuelSum(dist int) int {
    return dist * (dist + 1) / 2
}

func max(ints []int) int {
    max := ints[0]
    for _, val := range ints {
        if val > max {
            max = val
        }
    }
    return max
}

func part1() int {
    lines := reader.LinesFromFile("input.txt")
    crabXs := linesToInts(lines)
    logger.Logs.Infof("Crab positions: %d", crabXs)
    finalPosition, fuel := median(crabXs), 0
    for _, position := range crabXs {
        fuel += int(math.Abs(float64(position) - float64(finalPosition)))
    }
    return fuel
}

func part2() int {
    lines := reader.LinesFromFile("input.txt")
    crabXs := linesToInts(lines)
    logger.Logs.Infof("Crab positions: %d", crabXs)
    fuel, minFuel := 0, -1
    for position := 0; position < max(crabXs); position++ {
        fuel = 0
        for _, otherPos := range crabXs {
            fuel += fuelSum(int(math.Abs(float64(position) - float64(otherPos))))
        }
        //logger.Logs.Infof("Total fuel to move all crabs to position %d: %d", position, fuel)
        if fuel < minFuel || minFuel < 0 {
            minFuel = fuel
        }
    }
    return minFuel
}
