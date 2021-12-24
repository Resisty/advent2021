package main

import (
    "fmt"
    "regexp"
    "strconv"
    "strings"
    logger "advent2021/adventlogger"
    reader "advent2021/adventreader"
)

type Point struct {
    x, y int
}

func (p Point) String() string {
    return fmt.Sprintf("%d,%d", p.x, p.y)
}

type Board struct {
    points map[Point]string
    roomsXPos map[int]string
    amphipods []*Amphipod
    hallwayLen, energy int
}

func (b Board) Hash() string {
    boardString := strings.Repeat("#", b.hallwayLen)
    boardString += "\n"
    hallwayStr := ""
    for x := 0; x < b.hallwayLen; x++ {
        hallwayStr += b.points[Point{x, 0}]
    }
    boardString += hallwayStr + "\n"
    for i := 1; i < 3; i++ {
        hallwayStr = ""
        for x := 0; x < b.hallwayLen; x++ {
            if spot, ok := b.points[Point{x, i}]; ok {
                hallwayStr += spot
            } else {
                hallwayStr += "#"
            }
        }
        boardString += hallwayStr + "\n"
    }
    return boardString
}

func (b Board) Sprint() string {
    boardString := b.Hash()
    boardString += "\nEnergy: " + strconv.Itoa(b.energy) + "\n"
    return boardString
}

func (b *Board) Dup(selectedPod *Amphipod) (*Board, *Amphipod) {
    points := make(map[Point]string)
    pods := make([]*Amphipod, 0)
    for _, pod := range b.amphipods {
        if pod.position.x == selectedPod.position.x && pod.position.y == selectedPod.position.y {
            pods = append(pods, selectedPod)
        } else {
            pods = append(pods, pod)
        }
    }
    for point, mark := range b.points {
        points[point] = mark
    }
    return &Board{points: points, roomsXPos: b.roomsXPos, amphipods: pods, hallwayLen: b.hallwayLen, energy: b.energy}, selectedPod
}

func (b *Board) Finished() bool {
    moves := make([]Point, 0)
    allCorrectRooms := true
    for _, pod := range b.amphipods {
        moves = append(moves, pod.LegalMoves(b)...)
        point := Point{pod.roomXPos, 1}
        top := b.points[point]
        if top != pod.name {
            allCorrectRooms = false
        }
        point = Point{pod.roomXPos, 2}
        bot := b.points[point]
        if bot != pod.name {
            allCorrectRooms = false
        }
    }
    return allCorrectRooms && len(moves) == 0
}

func (b *Board) Dead() bool {
    moves := make([]Point, 0)
    allCorrectRooms := true
    for _, pod := range b.amphipods {
        moves = append(moves, pod.LegalMoves(b)...)
        point := Point{pod.roomXPos, 1}
        top := b.points[point]
        if top != pod.name {
            allCorrectRooms = false
        }
        point = Point{pod.roomXPos, 2}
        bot := b.points[point]
        if bot != pod.name {
            allCorrectRooms = false
        }
    }
    return ! allCorrectRooms && len(moves) == 0
}

func (b *Board) Move(a *Amphipod, p Point) {
    origPos := a.position
    nrg := a.Move(p)
    b.points[origPos] = "."
    b.points[a.position] = a.name
    b.energy += nrg
}


type Amphipod struct {
    name string
    position Point
    roomXPos int
}

func (a *Amphipod) Dup() *Amphipod{
    return &Amphipod{a.name, a.position, a.roomXPos}
}


func (a *Amphipod) Move(p Point) int {
    // y energy is distance into hallway + distance out of hallway
    dy := a.position.y + p.y
    dx := a.position.x - p.x
    if dx < 0 {
        dx = -dx
    }
    a.position.x = p.x
    a.position.y = p.y
    switch a.name {
    case "A":
        return dx + dy
    case "B":
        return 10 * dx + 10 * dy
    case "C":
        return 100 * dx + 100 * dy
    case "D":
        return 1000 * dx + 1000 * dy
    }
    panic(fmt.Sprintf("%s doesn't match legal names to tally energy!", a))
}

func (a Amphipod) String() string {
    return fmt.Sprintf("N(%s), P(%s), R(%d)", a.name, a.position, a.roomXPos)
}

func (a *Amphipod) LegalMoves(b *Board) []Point {
    top := Point{a.roomXPos, 1}
    bot := Point{a.roomXPos, 2}
    if b.points[top] == b.points[bot] && b.points[top] == a.name {
        // This 'pod and its sibling are both in their room
        return make([]Point, 0)
    }
    if a.position.x == bot.x && a.position.y == 2 {
        // This 'pod is at the bottom of its room, no reason to go anywhere
        return make([]Point, 0)
    }

    if a.position.y == 2 && b.points[Point{a.position.x, a.position.y - 1}] != "." {
        // amphipod blocking exit to hallway, no legal moves
        return make([]Point, 0)
    }
    // get moves to the left
    leftPoints := make([]Point, 0)
    passableLeft := make([]Point, 0)
    for i := a.position.x - 1; i >= 0; i-- {
        if b.points[Point{i, 0}] == "." {
            passableLeft = append(passableLeft, Point{i, 0})
        }
        if _, ok := b.roomsXPos[i]; ok {
                // can't stop in front of a room
                continue
        }
        if b.points[Point{i, 0}] != "." {
            // blocked, can't move further left
            break
        }
        // can move left
        leftPoints = append(leftPoints, Point{i, 0})
    }
    // get moves to the right
    rightPoints := make([]Point, 0)
    passableRight := make([]Point, 0)
    for i := a.position.x + 1; i < b.hallwayLen; i++ {
        if b.points[Point{i, 0}] == "." {
            passableRight = append(passableRight, Point{i, 0})
        }
        if _, ok := b.roomsXPos[i]; ok {
                // can't stop in front of a room
                continue
            }
        if b.points[Point{i, 0}] != "." {
            // blocked, can't move furtherright 
            break
        }
        // can move right 
        rightPoints = append(rightPoints, Point{i, 0})
    }
    lenHallwayLeft, lenHallwayRight := len(passableLeft), len(passableRight)
    // in the hallway, can only move to empty/sibling-occupied room
    if lenHallwayLeft > 0 && passableLeft[lenHallwayLeft - 1].x <= top.x  && a.position.x >= top.x {
        // leftmost open hallway space (if there is one) is at or left of our room (and we're to the right of it)
        if b.points[top] == "." && b.points[bot] == "." {
            // empty room, take it
            return []Point{bot}
        }
        if b.points[top] == "." && b.points[bot] == a.name {
            // sibling is in bottom room, safe to take top
            return []Point{top}
        }
    }
    if lenHallwayRight > 0 && passableRight[lenHallwayRight - 1].x >= top.x && a.position.x <= top.x{
        // rightmost open hallway space (if there is one) is at or right of our room (and we're to the left of it)
        if b.points[top] == "." && b.points[bot] == "." {
            // empty room, take it
            return []Point{bot}
        }
        if b.points[top] == "." && b.points[bot] == a.name {
            // sibling is in bottom room, safe to take top
            return []Point{top}
        }
    }
    if a.position.y == 0 {
        // we're in the hallway but can't move into a room; locked
        return make([]Point, 0)
    }
    return append(leftPoints, rightPoints...)
}

func cheapestBoard(bs []*Board) *Board {
    min := int(^uint(0) >> 1)
    mIndex := 0
    for i, board := range bs {
        if board.energy < min {
            min = board.energy
            mIndex = i
        }
    }
    return bs[mIndex]
}

type Cache map[string][]*Board

func allBoardStates(board *Board, cache, finished Cache) []*Board {
    boards := make([]*Board, 0)
    logger.Logs.Infof("Top of allBoarStates, printing board...")
    fmt.Print(board.Sprint())
    for _, amphipod := range board.amphipods {
        moves := amphipod.LegalMoves(board)
        for _, move := range moves {
            copyPod := amphipod.Dup()
            copyBoard, copyPod := board.Dup(copyPod) // copy board around moving 'pod
            copyBoard.Move(copyPod, move)
            if copyBoard.Dead() {
                // we don't care about dead-end boards
                logger.Logs.Infof("Board is dead, discard and check next move.")
                continue
            }
            if copyBoard.Finished() {
                // we found an answer, check the cache
                minFinished := copyBoard
                if finishedBoards, ok := finished[board.Hash()]; ok {
                    minFinished = cheapestBoard(append(finishedBoards, minFinished))
                }
                finished[board.Hash()] = []*Board{minFinished}
                return []*Board{minFinished}
            } else {
                // check the cache
                if _, ok := cache[copyBoard.Hash()]; ok {
                    // hit, we've been here before -> loop
                    continue
                } else {
                    // miss, recurse and cache
                    cache[copyBoard.Hash()] = []*Board{copyBoard}
                    nextBoards := allBoardStates(copyBoard, cache, finished)
                    boards = append(boards, nextBoards...)
                }
            }
        }
    }
    return boards
}

func boardFromInput(lines []string) (*Board) {
    hallwayLen := len(lines[1]) - 2
    roomsXPos := make(map[int]string)
    amphipods := make([]*Amphipod, 0)
    points := make(map[Point]string)
    amphNames := []string{"A", "B", "C", "D"}
    for i := 2; i < 4; i++ {
        amphLine := lines[i]
        reg := regexp.MustCompile(`([A-D])`)
        amphs := reg.FindAllStringIndex(amphLine, -1)
        for j, slice := range amphs {
            amphChar := amphLine[slice[0]:slice[1]]
            amphipod := &Amphipod{string(amphChar), Point{slice[0] - 1, i - 1}, 0}
            amphipods = append(amphipods, amphipod)
            roomsXPos[slice[0] - 1] = amphNames[j]
            points[Point{slice[0] - 1, i - 1}] = string(amphChar)
        }
    }
    // correct amphipods' room assignments
    for xPos, name := range roomsXPos {
        for _, amphipod := range amphipods {
            if amphipod.name == name {
                amphipod.roomXPos = xPos
            }
        }
    }
    for i := 0; i < hallwayLen; i++ {
        points[Point{i, 0}] = "."
    }
    board := &Board{points: points, roomsXPos: roomsXPos, amphipods: amphipods, hallwayLen: hallwayLen}
    logger.Logs.Infof("Dumping 'pods: %v", amphipods)
    return board
}

func main() {
    result := part1()
    logger.Logs.Infof("Part one result: %d", result)
    result = part2()
    logger.Logs.Infof("Part two result: %d", result)
}

func part1() int {
    lines := reader.LinesFromFile("test.txt")
    start := boardFromInput(lines)
    allBoards := allBoardStates(start, make(Cache), make(Cache))
    cheapest := cheapestBoard(allBoards) 
    fmt.Print(cheapest.Sprint())
    return cheapest.energy
}

func part2() int {
    //lines := reader.LinesFromFile("test.txt")
    return 4
}
