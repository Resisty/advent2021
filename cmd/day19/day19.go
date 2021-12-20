package main

import (
    "fmt"
    "math"
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

func GetDist3D(a, b Point) float64 {
    return math.Sqrt(float64((a.x - b.x) * (a.x - b.x) + (a.y - b.y) * (a.y - b.y) + (a.z - b.z) * (a.z - b.z)))
}

func vector(a, b Point) Point {
    return Point{x: b.x - a.x, y: b.y - a.y, z: b.z - a.z}
}

type ScannerDistsMap map[float64][][]int // map distance between (sets of) two indices (Points) in a Scanner

type Scanner struct {
    points []Point
    uniqPoints map[Point]struct{}
    distMap ScannerDistsMap
    label string
}

func NewScanner(label string) *Scanner {
    points := make([]Point, 0)
    uniqPoints := make(map[Point]struct{})
    distMap := make(ScannerDistsMap)
    return &Scanner{points: points, uniqPoints: uniqPoints, distMap: distMap, label: label}
}

func (s Scanner) Print() {
    fmt.Printf("Scanner %s has Points: %s\n", s.label, fmt.Sprint(s.uniqPoints))
    fmt.Println("Distances:")
    for dist, indexList := range s.distMap {
        for _, indices := range indexList {
            fmt.Printf("(%v -> %v) = %f\n", s.points[indices[0]], s.points[indices[1]], dist)
        }
    }
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

func (s *Scanner) InternalDistances() {
    uniqPoints := make(map[Point]struct{})
    distMap := make(map[float64][][]int)
    for _, point := range s.points {
        uniqPoints[point] = struct{}{}
    }
    reducePoints := make([]Point, 0)
    for point := range uniqPoints {
        reducePoints = append(reducePoints, point)
    }
    for i := 0; i < len(reducePoints); i++ {
        for j := i + 1; j < len(reducePoints); j++ {
            dist := GetDist3D(reducePoints[i], reducePoints[j])
            if _, ok := distMap[dist]; ok {
                distMap[dist] = append(distMap[dist], []int{i, j})
            } else {
                distMap[dist] = [][]int{{i, j}}
            }
        }
    }
    s.distMap = distMap
    s.uniqPoints = uniqPoints
    s.points = reducePoints
}

func (s *Scanner) reorient(matchVector Point, index1, index2 int) int {
    // Reorient all points in the scanner (around the scanner as origin) until the points at index1 and index2 match
    // the vector
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
            // check vec against index1->index2 here
            vec := vector(s.points[index1], s.points[index2])
            logger.Logs.Infof("Checking origin vector (%v) against 'forward' rotated vector (%s -> %s): (%v)", matchVector, s.points[index1], s.points[index2], vec)
            if matchVector == vec {
                // oriented correctly, return reference point for translation
                return index1
            }
            vec = vector(s.points[index2], s.points[index1])
            logger.Logs.Infof("Checking origin vector (%v) against 'reverse' rotated vector (%v)", matchVector, vec)
            if matchVector == vec {
                // oriented correctly, return reference point for translation
                return index2
            }
        }
        s.rotate90y() // next "lateral" face
    }
    panic("Didn't find any orientation among 24 orientations of a cube to match comparison vector " + fmt.Sprint(matchVector))
}

func (s *Scanner) absorb(other *Scanner, originPoint Point, refIndex int) {
    // s should always be the "origin"; absorb scanner's points relative to originPoint's origin
    logger.Logs.Infof("Origin scanner is absorbing scanner %s", other.label)
    s.Print()
    other.Print()
    referencePoint := other.points[refIndex]
    offsetVector := vector(originPoint, referencePoint)
    if offsetVector.x != 0 && offsetVector.y != 0 && offsetVector.z != 0 {
        panic(fmt.Sprintf("Did not reorient scanner %s correctly!", other.label))
    }
    for _, point := range s.points {
        newPoint := Point{offsetVector.x + point.x, offsetVector.y + point.y, offsetVector.z + point.z}
        logger.Logs.Infof("Scanner %s translating off of origin scanner's point %s (corrected: %s)", s.label, point, newPoint)
        s.points = append(s.points, newPoint)
    }
    s.InternalDistances()
}

func scannersFromInput(lines []string) []*Scanner {
    scanners := make([]*Scanner, 0)
    scanner := NewScanner("Blank")
    for _, line := range lines {
        if line == "" {
            scanner.InternalDistances()
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
        }
    }
    scanner.InternalDistances()
    scanners = append(scanners, scanner) // loop exits before final scanner is added
    return scanners
}

//func uniqPointsEqual(a, b map[Point]struct{}) bool {
//    if len(a) != len(b) {
//        return false
//    }
//    for point := range a {
//        if _, ok := b[point]; ! ok {
//            return false
//        }
//    }
//    for point := range b {
//        if _, ok := a[point]; ! ok {
//            return false
//        }
//    }
//    return true
//}

func main() {
    result := part1()
    logger.Logs.Infof("Part one result: %d", result)
    result = part2()
    logger.Logs.Infof("Part two result: %d", result)
}

func part1() int {
    lines := reader.LinesFromFile("test.txt")
    scanners := scannersFromInput(lines)
    logger.Logs.Infof("Got scanners:")
    for _, scanner := range scanners {
        scanner.Print()
    }
    origin := scanners[0]
    scanners = scanners[1:]
    for len(scanners) > 0 {
        // stop when we've successfully reoriented and absorbed all scanners to origin
        var scannerToReorientIndex int
        var dist float64
        DISTFOUND:
        for i, scanner := range scanners {
        // find a unique distance that only origin and scanner share; unique means within origin and scanner as well
            for distance := range origin.distMap {
                if len(origin.distMap[dist]) > 1 {
                    continue
                }
                if _, ok := scanner.distMap[distance]; ok {
                    if len(scanner.distMap[distance]) > 1 {
                        continue
                    }
                    dist = distance
                    scannerToReorientIndex = i
                    logger.Logs.Infof("Origin scanner and scanner %s each have a distance %f which spans ONLY two points per scanner", scanners[scannerToReorientIndex].label, dist)
                    break DISTFOUND
                }
            }
        }
        // reorient scanner to origin
        logger.Logs.Infof("origin.distMap[%f] == %v", dist, origin.distMap[dist])
        originPoint1Index := origin.distMap[dist][0][0]
        originPoint2Index := origin.distMap[dist][0][1]
        originP1 := origin.points[originPoint1Index]
        originP2 := origin.points[originPoint2Index]
        originVector := vector(originP1, originP2)
        logger.Logs.Infof("Setting origin vector off of points %s -> %s, value %s", originP1, originP2, originVector)
        scannerPoint1Index := scanners[scannerToReorientIndex].distMap[dist][0][0]
        scannerPoint2Index := scanners[scannerToReorientIndex].distMap[dist][0][1]
        scannerPoint1 := scanners[scannerToReorientIndex].points[scannerPoint1Index]
        scannerPoint2 := scanners[scannerToReorientIndex].points[scannerPoint2Index]
        logger.Logs.Infof("Reorienting scanner %s with reference points %s -> %s", scanners[scannerToReorientIndex].label, scannerPoint1, scannerPoint2)
        scannerPointIndex := scanners[scannerToReorientIndex].reorient(originVector, scannerPoint1Index, scannerPoint2Index)
        origin.absorb(scanners[scannerToReorientIndex], originP1, scannerPointIndex)
        scanners[scannerToReorientIndex] = scanners[len(scanners) - 1]
        scanners = scanners[:len(scanners) - 1]
    }
    logger.Logs.Infof("Unique points: %v", origin.uniqPoints)
    return len(origin.uniqPoints)
}

func part2() int {
    // lines := reader.LinesFromFile("test.txt")
    return 4 
}

