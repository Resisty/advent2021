package main

import (
    "regexp"
    logger "advent2021/adventlogger"
    reader "advent2021/adventreader"
)

type Polymer struct {
    pairs map[string]int
    instructs map[string]string
    counts map[string]int
}

func NewPolymer() *Polymer {
    pairs := make(map[string]int)
    instructs := make(map[string]string)
    counts := make(map[string]int)
    return &Polymer{pairs: pairs, instructs: instructs, counts: counts}
}

func stringToPairs(s string) map[string]int {
    pairs := make(map[string]int)
    left := s[:len(s) - 1] // all but last
    right := s[1:] // all but first
    for i, char := range left {
        pair := string(char) + string(right[i])
        pairs[pair] += 1
    }
    return pairs
}

func (p *Polymer) Step() {
    nextPairs := make(map[string]int)
    for pair, count := range p.pairs {
        insertChar := p.instructs[pair]
        p.counts[insertChar] += count
        triplet := string(pair[0]) + insertChar + string(pair[1])
        // split and count the two new pairs
        pairCount := stringToPairs(triplet)
        for newPair, newPairCount := range pairCount {
            // update the global counts for the new pairs
            nextPairs[newPair] += count * newPairCount
        }
    }
    p.pairs = nextPairs
}

func polymerFromInput(lines []string) *Polymer {
    template := ""
    gotTemplate := false
    polymer := NewPolymer()
    for _, line := range lines {
        if line == "" {
            continue
        }
        if ! gotTemplate {
            template = line
            for _, char := range line {
                polymer.counts[string(char)] += 1
            }
            polymer.pairs = stringToPairs(template)
            gotTemplate = true
        } else {
            reg := regexp.MustCompile(`([A-Z][A-Z])\s->\s([A-Z])`)
            result := reg.FindStringSubmatch(line)
            polymer.instructs[result[1]] = result[2]
        }
    }
    return polymer
}

func main() {
    result := part1()
    logger.Logs.Infof("Part one result: %d", result)
    result = part2()
    logger.Logs.Infof("Part two result: %d", result)
}

func part1() int {
    lines := reader.LinesFromFile("input.txt")
    polymer := polymerFromInput(lines)
    for i := 0; i < 10; i++ {
        polymer.Step()
    }
    min, max := -1, 0
    for _, count := range polymer.counts {
        if count > max {
            max = count
        }
        if count < min || min == -1 {
            min = count
        }
    }
    return max - min
}

func part2() int {
    lines := reader.LinesFromFile("input.txt")
    polymer := polymerFromInput(lines)
    for i := 0; i < 40; i++ {
        polymer.Step()
    }
    min, max := -1, 0
    for _, count := range polymer.counts {
        if count > max {
            max = count
        }
        if count < min || min == -1 {
            min = count
        }
    }
    return max - min
}
