package main

import (
    "fmt"
    "regexp"
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

func splitReg(str string, sep string) []string {
    return regexp.MustCompile(sep).Split(str, -1)
}

// Set code lovingly stolen from https://gist.github.com/bgadrian/cb8b9344d9c66571ef331a14eb7a2e80
type Set struct {
	list map[string]struct{} //empty structs occupy 0 memory
}

func (s Set) String() string {
    keys := make([]string, len(s.list))
    i := 0
    for k:= range s.list {
        keys[i] = k
        i++
    }
    return strings.Join(keys, "")
}

func NewSet() *Set {
	s := &Set{}
	s.list = make(map[string]struct{})
    return s
}

func SetFrom(chars string) *Set {
    s := NewSet()
    for _, char := range chars {
        s.Add(string(char))
    }
    return s
}

func (s *Set) Has(v string) bool {
	_, ok := s.list[v]
	return ok
}

func (s *Set) Add(v string) {
	s.list[v] = struct{}{}
}

func (s *Set) Remove(v string) {
	delete(s.list, v)
}

func (s *Set) In(r *Set) bool {
    // Return true if s is a subset of r
    for k := range s.list {
        if ! r.Has(k) {
            return false
        }
    }
    return true
}

func (s *Set) Intersection(r *Set) *Set {
    // Non-destructive intersection for hopefully obvious reasons
    a, b := s, r
    if len(s.list) > len(r.list) {
        a, b = r, s
    } else {
        a, b = s, r
    }
    result := NewSet()
    for k := range a.list {
        if _, ok := b.list[k]; ok {
            result.Add(k)
        }
    }
    return result
}

func (s *Set) Union(r *Set) *Set {
    // Non-destructive union for hopefully obvious reasons
    result := NewSet()
    for k := range s.list {
        result.Add(k)
    }
    for k := range r.list {
        result.Add(k)
    }
    return result
}

func (s *Set) Subtract(o *Set) *Set {
    // Non-destructive subtract for hopefully obvious reasons
    result := NewSet()
    result = result.Union(s)
    for k := range o.list {
        result.Remove(k)
    }
    return result
}

func (s *Set) Equals(o *Set) bool {
    return fmt.Sprint(s.list) == fmt.Sprint(o.list)
}

type Decoder struct {
    digits map[int]*Set
}

func (d *Decoder) DecodeInputs(strings []string) {
    otherDigits := make(map[int][]*Set)
    for i := range strings {
        s := SetFrom(strings[i])
        switch numDigits := len(strings[i]); numDigits {
        // Special cases: digits 1, 4, 7, and 8 have unique counts of characters
        case 2:
            d.digits[1] = s
        case 3:
            d.digits[7] = s
        case 4:
            d.digits[4] = s
        case 7:
            d.digits[8] = s
        default:
            otherDigits[numDigits] = append(otherDigits[numDigits], s)
        }
    }
    // Lost my goddamn mind on the rest of this function, but it works
    out:
    for _, length5Set := range otherDigits[5] {
        for _, length6Set := range otherDigits[6] {
            if length5Set.In(length6Set) && ! d.digits[1].In(length6Set){
                d.digits[5] = length5Set
                d.digits[6] = length6Set
                d.digits[9] = d.digits[5].Union(d.digits[1])
                // logger.Logs.Infof("Setting digit for 5 from set: %v", length5Set)
                // logger.Logs.Infof("Setting digit for 6 from set: %v", length6Set)
                // logger.Logs.Infof("Setting digit for 9 from set: %v", d.digits[9])
                break out
            }
        }
    }
    for _, length6Set := range otherDigits[6] {
        if ! length6Set.Equals(d.digits[6]) && ! length6Set.Equals(d.digits[9]) {
            d.digits[0] = length6Set
            break
        }
    }
    for _, length5Set := range otherDigits[5] {
        if length5Set.Equals(d.digits[5]) {
            continue
        }
        diff := d.digits[8].Subtract(length5Set)
        if diff.In(d.digits[5]) {
            d.digits[2] = length5Set
        } else {
            d.digits[3] = length5Set
        }
    }
}

func (d *Decoder) DecodeOutputs(strings []string) []int {
    results := make([]int, 0)
    for i := range strings {
        s := SetFrom(strings[i])
        for j, set := range d.digits {
            if set.Equals(s) {
                results = append(results, j)
                break
            }
        }
    }
    return results
}


func NewDecoder() *Decoder {
    d := &Decoder{}
    d.digits = make(map[int]*Set)
    return d
}

func inputsOutputs(lines []string) []map[string][]string {
    result := make([]map[string][]string, 0)
    for _, line := range lines {
        //logger.Logs.Infof("Parsing line: %s", line)
        inputOutput := splitReg(line, " \\| ")
        //logger.Logs.Infof("Parsed line out of ' | ' into inputs, outputs: %s", inputOutput)
        inputs := splitReg(inputOutput[0], " ")
        outputs := splitReg(inputOutput[1], " ")
        entry := make(map[string][]string)
        entry["inputs"] = inputs
        entry["outputs"] = outputs
        result = append(result, entry)
    }
    return result
}

func part1() int {
    lines := reader.LinesFromFile("input.txt")
    sum := 0
    entries := inputsOutputs(lines)
    for i := range entries {
        // logger.Logs.Infof("%dth entry (previously a line of input): %s", i, entries[i])
        for _, output := range entries[i]["outputs"] {
            numDigits := len(output)
            // logger.Logs.Infof("Number of digits in output string '%s': %d", output, numDigits)
            switch numDigits {
            case 2:
                fallthrough
            case 3:
                fallthrough
            case 4:
                fallthrough
            case 7:
                sum += 1
            }
        }
    }
    return sum
}

func part2() int {
    lines := reader.LinesFromFile("input.txt")
    entries := inputsOutputs(lines)
    sum := 0
    for i := range entries {
        // logger.Logs.Infof("%dth entry (previously a line of input): %s", i, entries[i])
        decoder := NewDecoder()
        decoder.DecodeInputs(entries[i]["inputs"])
        //logger.Logs.Infof("Decoded inputs: %v", decoder.digits)
        results := decoder.DecodeOutputs(entries[i]["outputs"])
        // logger.Logs.Infof("Encoded outputs: %v", entries[i]["outputs"])
        // logger.Logs.Infof("Decoded outputs: %v", results)
        str := ""
        for num := range results {
            str += strconv.Itoa(results[num])
        }
        tmp, _ := strconv.Atoi(str)
        sum += tmp
    }
    return sum
}
