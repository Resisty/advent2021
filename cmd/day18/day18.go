package main

import (
    "fmt"
    "regexp"
    "strconv"
    "strings"
    logger "advent2021/adventlogger"
    reader "advent2021/adventreader"
)

// lovingly stolen from https://stackoverflow.com/a/10030772/895246
func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func findSimplePairs(snailNum string) ([][]int, [][]int) {
    reg := regexp.MustCompile(`(\[[0-9]+,[0-9]+\])`)
    result := reg.FindAllStringIndex(snailNum, -1)
    actualPairs := make([][]int, 0)
    for _, slice := range result {
        numReg := regexp.MustCompile(`([0-9]+),([0-9]+)`)
        numResult := numReg.FindStringSubmatch(snailNum[slice[0]:slice[1]])
        lNum, _ := strconv.Atoi(numResult[1])
        rNum, _ := strconv.Atoi(numResult[2])
        actualPairs = append(actualPairs, []int{lNum, rNum})
    }
    return result, actualPairs
}

func findDepth(snailNum string, slice []int) int {
    depth := 0
    for i, char := range snailNum {
        switch string(char) {
        case `[`:
            depth += 1
        case `]`:
            depth -= 1
        }
        if i == slice[0] {
            break
        }
    }
    return depth
}

func findNext(snailNum string, index, offset int, left bool) ([]int, int) {
    reg := regexp.MustCompile(`([0-9]+)`)
    substring := snailNum[index + offset:]
    if left {
        substring = Reverse(snailNum[:index + offset])
    }
    result := reg.FindStringIndex(substring)
    if result == nil {
        return nil, -1
    }
    numString := substring[result[0]:result[1]]
    if left {
        numString = Reverse(numString)
        result = []int{index - result[1] + offset, index - result[0] + offset}
    } else {
        result = []int{index + result[0] + offset, index + result[1] + offset}
    }
    num, _ := strconv.Atoi(numString)
    return result, num
}

func tryExplode(snailNum string) (string, bool) {
    simpleNums, pairs := findSimplePairs(snailNum)
    for i, slice := range simpleNums {
        depth := findDepth(snailNum, slice)
        if depth >= 5 { // Don't include the `[` of the lowest nested number to reach depth "4"
            lNum, rNum := pairs[i][0], pairs[i][1]
            offsetLeft, offsetMid := 0, 0
            lefterNumSlice, lefterNum := findNext(snailNum, slice[0], 0, true)
            if lefterNumSlice != nil {
                lefterNum += lNum
                lefterNumString := strconv.Itoa(lefterNum)
                offsetLeft += len(lefterNumString) - (lefterNumSlice[1] - lefterNumSlice[0]) // new lNum might grow (9 -> 10) or shrink (10 -> 9)
                lChunk, rChunk := snailNum[:lefterNumSlice[0]], snailNum[lefterNumSlice[1]:]
                snailNum = lChunk + lefterNumString + rChunk
            }
            lChunk, rChunk := snailNum[:slice[0] + offsetLeft], snailNum[slice[1] + offsetLeft:]
            offsetMid += 1 - (slice[1] - slice[0]) // len("0") == 1 - len("[145,72]") == 8 === -7 (shrink by 7)
            snailNum = lChunk + "0" + rChunk
            righterNumSlice, righterNum := findNext(snailNum, slice[1], offsetLeft + offsetMid, false)
            if righterNumSlice != nil {
                righterNum += rNum
                righterNumString := strconv.Itoa(righterNum)
                lChunk, rChunk := snailNum[:righterNumSlice[0]], snailNum[righterNumSlice[1]:]
                snailNum = lChunk + righterNumString + rChunk
            }
            return snailNum, true
        }
    }
    return snailNum, false
}

func trySplit(snailNum string) (string, bool) {
    reg := regexp.MustCompile(`(\d\d+)`)
    result := reg.FindStringIndex(snailNum)
    if result == nil {
        return snailNum, false
    }
    numString := snailNum[result[0]:result[1]]
    num, _ := strconv.Atoi(numString)
    lNum, rNum := num / 2, num / 2
    if num % 2 == 1 {
        rNum += 1
    }
    snailNum = snailNum[:result[0]] + fmt.Sprintf("[%d,%d]", lNum, rNum) + snailNum[result[1]:]
    return snailNum, true
}

func add(lNum, rNum string) string {
    return fmt.Sprintf("[%s,%s]", lNum, rNum)
}

func pairMagnitude (pair []int) int {
    return 3 * pair[0] + 2 * pair[1]
}

func magnitude(lNum string) int {
    for strings.Contains(lNum, "[") {
        result, actualPairs := findSimplePairs(lNum)
        mag := pairMagnitude(actualPairs[0])
        magString := strconv.Itoa(mag)
        lChunk, rChunk := lNum[:result[0][0]], lNum[result[0][1]:]
        lNum = lChunk + magString + rChunk
    }
    result, _ := strconv.Atoi(lNum)
    return result
}

func main() {
    result := part1()
    logger.Logs.Infof("Part one result: %d", result)
    result = part2()
    logger.Logs.Infof("Part two result: %d", result)
}

func reduce(snailNum string) string {
    yes := true
    for yes {
        snailNum, yes = tryExplode(snailNum)
        if yes {
            continue
        }
        snailNum, yes = trySplit(snailNum)
    }
    return snailNum
}

func part1() int {
    lines := reader.LinesFromFile("input.txt")
    lNum := lines[0]
    lNum = reduce(lNum)
    for i := 1; i < len(lines); i++ {
        lNum = add(lNum, lines[i])
        lNum = reduce(lNum)
    }
    return magnitude(lNum)
}

func part2() int {
    lines := reader.LinesFromFile("input.txt")
    linePairs := make(map[string]struct{})
    for i := 0; i < len(lines); i++ {
        for j := i + 1; j < len(lines); j++ {
            linePairs[add(lines[i], lines[j])] = struct{}{}
            linePairs[add(lines[j], lines[i])] = struct{}{}
        }
    }
    max, mag := 0, 0
    for snailNum := range linePairs { 
        snailNum = reduce(snailNum)
        mag = magnitude(snailNum)
        if mag > max {
            max = mag
        }
    }
    return max
}
