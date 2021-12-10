package main

import (
    "bufio"
    "bytes"
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

func part1() int {
    input := bytes.NewBuffer(reader.FromFile("input.txt"))
    start, prev, count := false, 0, 0
    scanner := bufio.NewScanner(input)
    for scanner.Scan() {
        entry, _ := strconv.Atoi(scanner.Text())
        if ! start {
            // logger.Logs.Infof("Started: %t. Setting prev to %d", start, entry)
            prev = entry
            start = true
            continue
        }
        // logger.Logs.Infof("Comparing entry (%d) -ge prev (%d)", intEntry, intPrev)
        if entry > prev {
            count += 1
            // logger.Logs.Infof("Entry (%d) IS GREATER than prev (%d)", entry, prev)
            // logger.Logs.Infof("Incremented count: %d", count)
        }
        // logger.Logs.Infof("Updating prev (%d) to entry (%d)", prev, entry)
        prev = entry
    }
    return count
}

func part2() int {
    input := bytes.NewBuffer(reader.FromFile("input.txt"))
    scanner := bufio.NewScanner(input)
    var entries []int
    for scanner.Scan() {
        entry, _ := strconv.Atoi(scanner.Text())
        entries = append(entries, entry)
    }
    start, prev, count := false, 0, 0
    for i := 0; i < len(entries) - 2; i++ {
        entry := entries[i] + entries[i+1] + entries[i+2]
        if ! start { 
            start = true
            prev = entry
            // logger.Logs.Infof("Started: %t. Setting prev to %d", start, entry)
            continue
        }
        // logger.Logs.Infof("Comparing entry (%d) -ge prev (%d)", entry, prev)
        if entry > prev { 
            // logger.Logs.Infof("Entry (%d) IS GREATER than prev (%d)", entry, prev)
            // logger.Logs.Infof("Incremented count: %d", count)
            count += 1
        }
        // logger.Logs.Infof("Updating prev (%d) to entry (%d)", prev, entry)
        prev = entry
    }
    return count
}
