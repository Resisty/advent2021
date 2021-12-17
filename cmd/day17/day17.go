package main

    // "sort"
    // "strings"
import (
    "fmt"
    "regexp"
    "strconv"
    logger "advent2021/adventlogger"
    reader "advent2021/adventreader"
)

type Area struct {
    minX, maxX, minY, maxY int
}

func (a Area) String() string {
    return fmt.Sprintf("%d < x < %d, %d < y < %d", a.minX, a.maxX, a.minY, a.maxY)
}

func targetArea(lines []string) Area {
    reg := regexp.MustCompile(`target area: x=([0-9-]+)..([0-9-]+), y=([0-9-]+)..([0-9-]+)`)
    result := reg.FindStringSubmatch(lines[0])
    results := make([]int, 0)
    for i := 1; i < 5; i++ {
        tmp, _ := strconv.Atoi(result[i])
        results = append(results, tmp)
    }
    area := Area{results[0], results[1], results[2], results[3]}
    return area
}

func Sigma(n int) int {
    sum := 0
    mult := 1
    if n < 0 {
        n = -n
        mult = -1
    }
    for i := 0; i <= n; i++ {
        sum += i
    }
    return sum * mult
}

func minXVelocity(a Area) int {
    i := 0
    dir := 1
    minX, maxX := a.minX, a.maxX
    if a.minX < 0 && a.maxX < 0 {
        dir = -1
        minX = -minX
        maxX = -maxX
    }
    for {
        sigma := Sigma(i)
        if minX <= dir * sigma && dir * sigma <= maxX {
            break
        }
        i += dir
    }
    return i
}

func maxYVelocity(a Area) int {
    i := 0
    dir := 1
    if a.minY < 0 {
        dir = -1
        for {
            next := i + dir * 1
            next2 := i + dir * 2
            if a.minY <= next && next <= a.maxY && (a.minY > next2 || next2 > a.maxY) {
                break
            }
            i += dir
        }
    } else {
        for {
            next := i + 1 * dir
            if a.minY <= i && i <= a.maxY && (a.minY > next || next > a.maxY) {
                break
            }
            i += dir
        }
    }
    return i * dir
}

func pointInArea(a Area, x, y int) bool {
    return a.minX <= x && x <= a.maxX && a.minY <= y && y <= a.maxY
}

func xyVelInArea(a Area, x, y int) bool {
    xPos, yPos := 0, 0
    xDir := 1
    if x < 0 {
        xDir = -1
    }
    for {
        xPos += xDir * x
        if x != 0 {
            x -= xDir * 1
        }
        yPos += y
        y -= 1
        if pointInArea(a, xPos, yPos) {
            return true
        }
        if y < 0 && yPos < a.minY {
            // falling and below bottom of target; will never reach
            return false
        }
    }
}

func main() {
    result := part1()
    logger.Logs.Infof("Part one result: %d", result)
    result = part2()
    logger.Logs.Infof("Part two result: %d", result)
}

func part1() int {
    lines := reader.LinesFromFile("input.txt")
    area := targetArea(lines)
    logger.Logs.Infof("Got area: %v", area)
    logger.Logs.Infof("Minimal x-velocity to reach target: %d", minXVelocity(area))
    logger.Logs.Infof("Maximal y-velocity to reach target: %d", maxYVelocity(area))
    return Sigma(maxYVelocity(area))
}

func part2() int {
    lines := reader.LinesFromFile("input.txt")
    area := targetArea(lines)
    logger.Logs.Infof("Got area: %v", area)
    minXVel := minXVelocity(area)
    maxXVel := area.maxX
    minYVel := area.minY
    maxYVel := maxYVelocity(area)
    count := 0
    for xVel := minXVel; xVel <= maxXVel; xVel++ {
        for yVel := minYVel; yVel <= maxYVel; yVel++ {
            if xyVelInArea(area, xVel, yVel) {
                count++
            }
        }
    }
    return count
}

