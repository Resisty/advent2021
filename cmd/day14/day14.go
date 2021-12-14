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
    min, max int
}

func NewPolymer() *Polymer {
    pairs := make(map[string]int)
    instructs := make(map[string]string)
    counts := make(map[string]int)
    return &Polymer{pairs: pairs, instructs: instructs, counts: counts, min: -1, max: 0}
}

func stringToPairs(s string) map[string]int {
    pairs := make(map[string]int)
    left := s[:len(s) - 1] // all but last
    right := s[1:] // all but first
    for i, char := range left {
        pairs[string(char) + string(right[i])] += 1
    }
    return pairs
}

func (p *Polymer) Step() {
    nextPairs := make(map[string]int)
    for pair, count := range p.pairs {
        insertChar := p.instructs[pair]
        p.counts[insertChar] += count
        triplet := string(pair[0]) + insertChar + string(pair[1])
        pairCount := stringToPairs(triplet)
        for newPair, newPairCount := range pairCount {
            nextPairs[newPair] += count * newPairCount
        }
    }
    p.pairs = nextPairs
}

func (p *Polymer) MinMax() {
    for _, count := range p.counts {
        if count > p.max {
            p.max = count
        }
        if count < p.min || p.min == -1 {
            p.min = count
        }
    }
}

func polymerFromInput(lines []string) *Polymer {
    polymer := NewPolymer()
    template, _, instructLines := lines[0], lines[1], lines[2:]
    for _, char := range template {
        polymer.counts[string(char)] += 1
    }
    polymer.pairs = stringToPairs(template)
    for _, line := range instructLines {
        reg := regexp.MustCompile(`([A-Z][A-Z])\s->\s([A-Z])`)
        result := reg.FindStringSubmatch(line)
        polymer.instructs[result[1]] = result[2]
    }
    return polymer
}

func main() {
    result := run(10)
    logger.Logs.Infof("Part one result: %d", result)
    result = run(40)
    logger.Logs.Infof("Part two result: %d", result)
}

func run(iterations int) int {
    lines := reader.LinesFromFile("input.txt")
    polymer := polymerFromInput(lines)
    for i := 0; i < iterations; i++ {
        polymer.Step()
    }
    polymer.MinMax()
    return polymer.max - polymer.min
}
