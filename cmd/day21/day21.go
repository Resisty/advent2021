package main

import (
    "fmt"
    "regexp"
    "strconv"
    logger "advent2021/adventlogger"
    reader "advent2021/adventreader"
)

type Die interface {
    Roll() int
    Rolls() int
    String() string
}

type D100 struct {
    face, rolls int
}

func NewD100() *D100 {
    return &D100{0, 0}
}

func (d *D100) Roll() int {
    d.face = (d.face + 1) % 101
    if d.face == 0 {
        d.face += 1
    }
    d.rolls += 1
    return d.face
}

func (d *D100) Rolls() int {
    return d.rolls
}

func (d D100) String() string {
    return fmt.Sprintf("Face: %d, rolls: %d", d.face, d.rolls)
}

type Player struct {
    number, position, score int
}

func NewPlayer(number, start int) *Player {
    return &Player{number, start, 0}
}

func sum(nums ...int) int {
    res := 0
    for _, num := range nums {
        res += num
    }
    return res
}

func (p *Player) advance(spaces int) {
    p.position = (p.position + spaces) % 10
    if p.position == 0 {
        p.score += 10
    } else {
        p.score += p.position
    }
}

func (p *Player) Turn(dice Die) {
    rolls := make([]int, 0)
    for i := 0; i < 3; i++ {
        roll := dice.Roll()
        rolls = append(rolls, roll)
    }
    p.advance(sum(rolls...))
}

func (p Player) String() string {
    return fmt.Sprintf("Order: %d, Position: %d, score: %d", p.number, p.position, p.score)
}

type Game struct {
    players []*Player
    dice Die
    gameOver int
}

func NewGame(numPlayers, finalScore int, die Die) *Game {
    players := make([]*Player, numPlayers)
    dice := die
    return &Game{players, dice, finalScore}
}

func (g *Game) Round() {
    for _, player := range g.players {
        player.Turn(g.dice)
        if g.Over() {
            return
        }
    }
}

func (g *Game) Over() bool {
    for _, player := range g.players {
        if player.score >= g.gameOver {
            return true
        }
    }
    return false
}

func (g *Game) Loser() []*Player {
    if ! g.Over() {
        return nil
    }
    losers := make([]*Player, 0)
    for _, player := range g.players {
        if player.score < g.gameOver {
            losers = append(losers, player)
        }
    }
    return losers
}

func gameFromInput(lines []string, gameOver int, dice Die) *Game {
    numPlayers := len(lines)
    game := NewGame(numPlayers, gameOver, dice)
    for _, line := range lines {
        reg := regexp.MustCompile(`Player\s(\d+)\sstarting\sposition:\s(\d+)`)
        result := reg.FindStringSubmatch(line)
        playerNumString, positionString := result[1], result[2]
        playerNum, _ := strconv.Atoi(playerNumString)
        position, _ := strconv.Atoi(positionString)
        game.players[playerNum - 1] = NewPlayer(playerNum, position) // 0-indexed
    }
    return game
}

// ffs, start everything over for part 2

type BoardState struct {
    turnLabel, position, score, otherPos, otherScore int
}

type Cache map[BoardState][]int

func NewBoardState(p1, p2 Player, roll int) BoardState {
    position := ((p1.position + roll - 1 ) % 10) + 1
    score := p1.score + position
    return BoardState{p1.number, position, score, p2.position, p2.score}
}

func winners(gameOver int, turnPlayer, otherPlayer Player, cache Cache) (int, int){
    sum1, sum2 := 0, 0
    possibleRolls := map[int]int{3: 1, 4: 3, 5: 6, 6: 7, 7: 6, 8: 3, 9: 1}
    for roll, numUniverses := range possibleRolls {
        copyPlayer := Player{turnPlayer.number, turnPlayer.position, turnPlayer.score}
        boardState := NewBoardState(copyPlayer, otherPlayer, roll)
        copyPlayer.position = boardState.position
        copyPlayer.score = boardState.score
        if copyPlayer.score >= gameOver {
            sum1 += numUniverses
        } else {
            if wins, ok := cache[boardState]; ok {
                sum1 += (wins[0] * numUniverses) 
                sum2 += (wins[1] * numUniverses)
            } else {
                otherPlayerWins, turnPlayerWins := winners(gameOver, otherPlayer, copyPlayer, cache)
                sum2 += (otherPlayerWins * numUniverses)
                sum1 += (turnPlayerWins * numUniverses)
                cache[boardState] = []int{turnPlayerWins, otherPlayerWins}
            }
        }
    }
    return sum1, sum2
}
// back to our regularly scheduled programming

func main() {
    result := part1()
    logger.Logs.Infof("Part one result: %d", result)
    result = part2()
    logger.Logs.Infof("Part two result: %d", result)
}

func part1() int {
    lines := reader.LinesFromFile("input.txt")
    game := gameFromInput(lines, 1000, NewD100())
    for ! game.Over() {
        game.Round()
    }
    rolls :=  game.dice.Rolls()
    loserScore := game.Loser()[0].score
    logger.Logs.Infof("Game over with losing score %d, die rolls %d", loserScore, rolls)
    return rolls * loserScore
}

func part2() int {
    lines := reader.LinesFromFile("input.txt")
    game := gameFromInput(lines, 1000, NewD100())
    player1 := game.players[0]
    player2 := game.players[1]
    for i:= 1; i < 22; i++ {
        cache := make(Cache)
        player1wins, player2wins := winners(i, *player1, *player2, cache)
        logger.Logs.Infof("Target score %d: %d %d", i, player1wins, player2wins)
    }
    return 4
}
