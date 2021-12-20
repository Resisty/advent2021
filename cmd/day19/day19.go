package main

import (
    "fmt"
    "regexp"
    "strconv"
    logger "advent2021/adventlogger"
    reader "advent2021/adventreader"
)

type Point struct {
    x, y, z int
}

func (p Point) String() string {
    return fmt.Sprintf("%d,%d,%d", p.x, p.y, p.z)
}

func vector(a, b Point) Point {
    return Point{x: b.x - a.x, y: b.y - a.y, z: b.z - a.z}
}

type Scanner struct {
    points []Point
    uniqPoints map[Point]struct{}
    label string
}

func NewScanner(label string) *Scanner {
    points := make([]Point, 0)
    uniqPoints := make(map[Point]struct{})
    return &Scanner{points: points, uniqPoints: uniqPoints, label: label}
}

func (s *Scanner) dup() *Scanner {
    newScan := NewScanner(s.label)
    newScan.points = make([]Point, 0)
    for _, point := range s.points {
        newScan.points = append(newScan.points, point)
        newScan.uniqPoints[point] = struct{}{}
    }
    return newScan
}

func (s Scanner) String() string {
    return fmt.Sprintf("%s, Points: %s", s.label, fmt.Sprint(s.points))
}

func (s *Scanner) rotate90x() {
    // rotate 90 around "x" axis
    for i := range s.points {
        s.points[i] = Point{s.points[i].x, s.points[i].z, -(s.points[i]).y}
    }
}

func (s *Scanner) rotate90y() {
    // rotate 90 around "y" axis
    for i := range s.points {
        s.points[i] = Point{-(s.points[i].z), s.points[i].y, s.points[i].x}
    }
}

func (s *Scanner) rotate90z() {
    // rotate 90 around "z" axis
    for i := range s.points {
        s.points[i] = Point{s.points[i].y, -(s.points[i].x), s.points[i].z}
    }
}

func (s *Scanner) translate(offsetVector Point) {
    points := make([]Point, 0)
    uniqPoints := make(map[Point]struct{})
    for _, point := range s.points {
        translated := Point{point.x + offsetVector.x, point.y + offsetVector.y, point.z + offsetVector.z}
        points = append(points, translated)
        uniqPoints[translated] = struct{}{}
    }
    s.points = points
    s.uniqPoints = uniqPoints
}


func (s *Scanner) createOrientations() []*Scanner {
    // create all possible orientations of the scanner and return them
    allRotations := make([]*Scanner, 0)
    for i := 0; i < 6; i++ {
        if i == 4 {
            // We did 4 "lateral" faces, time for "top"
            s.rotate90z()
        }
        if i == 5 {
            // We did 4 "lateral" and the "top" face, time for the "bottom"
            s.rotate90z()
            s.rotate90z()
        }
        for j := 0; j < 4; j++ {
            s.rotate90x()
            allRotations = append(allRotations, s.dup())
        }
        if i < 4 {
            s.rotate90y() // next "lateral" face
        }
    }
    return allRotations
}

func scannersFromInput(lines []string) []*Scanner {
    scanners := make([]*Scanner, 0)
    scanner := NewScanner("Blank")
    for _, line := range lines {
        if line == "" {
            scanners = append(scanners, scanner)
            continue
            // complete the "scanner"
        }
        reg := regexp.MustCompile(`---\sscanner\s(\d+)\s---`)
        if result := reg.FindStringSubmatch(line); result != nil {
            scanner = NewScanner(result[1])
            continue
        }
        reg = regexp.MustCompile(`([0-9-]+),([0-9-]+),([0-9-]+)`)
        result := reg.FindStringSubmatch(line)
        if result != nil {
            x, y, z := result[1], result[2], result[3]
            xInt, _ := strconv.Atoi(x)
            yInt, _ := strconv.Atoi(y)
            zInt, _ := strconv.Atoi(z)
            p := Point{x: xInt, y: yInt, z: zInt}
            scanner.points = append(scanner.points, p)
            scanner.uniqPoints[p] = struct{}{}
        }
    }
    scanners = append(scanners, scanner) // loop exits before final scanner is added
    return scanners
}

func pairs(scanner, oriented *Scanner) [][]Point {
    allPairs := make([][]Point, 0)
    for _, scannerPoint := range scanner.points {
        for _, orientedPoint := range oriented.points {
            pair := []Point{scannerPoint, orientedPoint}
            allPairs = append(allPairs, pair)
        }
    }
    return allPairs
}

func overlapping(scanner, oriented *Scanner, offsetVector Point) int {
    count := 0
    for _, point := range scanner.points {
        adjusted := Point{point.x + offsetVector.x, point.y + offsetVector.y, point.z + offsetVector.z}
        if _, ok := oriented.uniqPoints[adjusted]; ok {
            count += 1
        }
    }
    return count
}

func main() {
    result := part1()
    logger.Logs.Infof("Part one result: %d", result)
    result = part2()
    logger.Logs.Infof("Part two result: %d", result)
}

func part1() int {
    lines := reader.LinesFromFile("input.txt")
    scanners := scannersFromInput(lines)
    totalScanners := len(scanners)
    origin := scanners[0]
    scanners = scanners[1:]
    unoriented := make(map[string][]*Scanner)
    for _, scanner := range scanners {
        unorientedList := make([]*Scanner, 0)
        unoriented[scanner.label] = append(unorientedList, scanner.createOrientations()...)
    }
    oriented := make([]*Scanner, 0)
    oriented = append(oriented, origin)
    minTolerance := 12
    translation := Point{0, 0, 0}
    for len(oriented) < totalScanners {
        max := 0
        acceptableOrientedIndex := -1
        removeLabel := ""
        for label, unorientedList := range unoriented {
            for i, scanner := range unorientedList {
                for _, orientedScanner := range oriented {
                    // get all pairs of points between orientedScanner and scanner
                    pointPairs := pairs(orientedScanner, scanner)
                    for _, pair := range pointPairs {
                        offsetVector := vector(pair[1], pair[0]) // translate from unoriented to oriented
                        nOverlap := overlapping(scanner, orientedScanner, offsetVector) // count the number of overlapping points if we apply the offset to scanner
                        if nOverlap > max {
                            max = nOverlap
                            if nOverlap >= minTolerance {
                                acceptableOrientedIndex = i
                                removeLabel = label
                                translation = offsetVector
                            }
                        }
                    }
                }
            }
        }
        if acceptableOrientedIndex > -1 {
            // we found an index (rotation of some scanner of label "foo") with a maximal overlap greater than minimal tolerance
            // count it as oriented
            nowOriented := unoriented[removeLabel][acceptableOrientedIndex]
            nowOriented.translate(translation)
            oriented = append(oriented, nowOriented)
            delete(unoriented, removeLabel)
        }
    }
    uniqPoints := make(map[Point]struct{})
    for _, scanner := range oriented {
        for _, point := range scanner.points {
            uniqPoints[point] = struct{}{}
        }
    }
    return len(uniqPoints)
}

func part2() int {
    // lines := reader.LinesFromFile("test.txt")
    return 4 
}
