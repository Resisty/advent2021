package main

    // "sort"
    // "strings"
import (
    "fmt"
    "strconv"
    logger "advent2021/adventlogger"
    reader "advent2021/adventreader"
)

type Packet struct {
    hexString, binString string
    version, packType, value int
    subPackets []*Packet
}

func NewPacket(hex, bin string) *Packet {
    subs := make([]*Packet, 0)
    return &Packet{hexString: hex, binString: bin, version: 0, packType: 0, value: 0, subPackets: subs}
}

func (p Packet) String() string {
    return fmt.Sprintf("Hex: '%s', Bin: '%s', Version: '%d', Type: '%d', Value: '%d', SubPackets: '%v'", p.hexString, p.binString, p.version, p.packType, p.value, p.subPackets)
}

func (p *Packet) LShift(num int) string {
    result := p.binString[:num]
    p.binString = p.binString[num:]
    return result
}

func (p *Packet) VersionSum() int {
    sum := p.version
    for _, subPack := range p.subPackets {
        sum += subPack.VersionSum()
    }
    return sum
}

func (p *Packet) Sum() int {
    sum := 0
    for _, packet := range p.subPackets {
        sum += packet.Resolve()
    }
    return sum
}

func (p *Packet) Product() int {
    result := 1
    for _, packet := range p.subPackets {
        result *= packet.Resolve()
    }
    return result
}

func (p *Packet) Min() int {
    min := int(^uint(0) >> 1)
    for _, packet := range p.subPackets {
        val := packet.Resolve()
        if val < min {
            min = val
        }
    }
    return min
}

func (p *Packet) Max() int {
    max := 0 - int(^uint(0) >> 1) - 1
    for _, packet := range p.subPackets {
        val := packet.Resolve()
        if val > max {
            max = val
        }
    }
    return max
}

func (p *Packet) GreaterThan() int {
    if len(p.subPackets) != 2 {
        panic(fmt.Sprintf("Packet (%v) should be operator: GreaterThan with exactly two subpackets!", p))
    }
    left := p.subPackets[0].Resolve()
    right := p.subPackets[1].Resolve()
    if left > right {
        return 1
    } else {
        return 0
    }
}

func (p *Packet) LessThan() int {
    if len(p.subPackets) != 2 {
        panic(fmt.Sprintf("Packet (%v) should be operator: GreaterThan with exactly two subpackets!", p))
    }
    left := p.subPackets[0].Resolve()
    right := p.subPackets[1].Resolve()
    if left < right {
        return 1
    } else {
        return 0
    }
}

func (p *Packet) EqualTo() int {
    if len(p.subPackets) != 2 {
        panic(fmt.Sprintf("Packet (%v) should be operator: GreaterThan with exactly two subpackets!", p))
    }
    left := p.subPackets[0].Resolve()
    right := p.subPackets[1].Resolve()
    if left == right {
        return 1
    } else {
        return 0
    }
}

func (p *Packet) Resolve() int {
    switch p.packType {
    case 0:
        return p.Sum()
    case 1:
        return p.Product()
    case 2:
        return p.Min()
    case 3:
        return p.Max()
    case 4:
        return p.value
    case 5:
        return p.GreaterThan()
    case 6:
        return p.LessThan()
    case 7:
        return p.EqualTo()
    default:
        panic("You got a type that doesn't line up!")
    }
}

func (p *Packet) Parse() int{
    totalShift := 0
    p.version = binToInt(p.binString[:3])
    p.packType = binToInt(p.binString[3:6])
    p.LShift(6)
    totalShift += 6
    if p.packType == 4 {
        // literal
        valueString := ""
        for string(p.binString[0]) == "1" {
            tmp := p.LShift(5)
            totalShift += 5
            valueString += tmp[1:]
        }
        // once more for final literal starting with 0
        tmp := p.LShift(5)
        totalShift += 5
        valueString += tmp[1:]
        p.value = binToInt(valueString)
    } else {
        // operator
        lengthType := binToInt(p.LShift(1))
        totalShift += 1
        if lengthType == 0 {
            // next 15 bits => total length in bits of subpackets
            subPacketsBitLength := binToInt(p.LShift(15))
            totalShift += 15
            subPacketBits := p.LShift(subPacketsBitLength)
            totalShift += subPacketsBitLength
            for {
                subPacket := NewPacket("", subPacketBits)
                subPacket.Parse()
                p.subPackets = append(p.subPackets, subPacket)
                subPacketBits = subPacket.binString
                if binToInt(subPacketBits) == 0 {
                    break
                }
            }
        } else {
            // next 11 bits => number of subpackets
            subPacketsLength := binToInt(p.LShift(11))
            totalShift += 11
            for i := 0; i < subPacketsLength; i++ {
                subPacket := NewPacket("", p.binString) // copy binstring to subpacket
                shifted := subPacket.Parse()
                totalShift += shifted
                p.LShift(shifted) // update original to match shifted bits in parsed copy
                p.subPackets = append(p.subPackets, subPacket)
            }
        }
    }
    return totalShift
}

func binToInt(s string) int {
    i64, _ := strconv.ParseInt(s, 2, 64)
    return int(i64)
}

func expandToBin(s string) string {
    hex2bin := map[string]string{
        "0": "0000",
        "1": "0001",
        "2": "0010",
        "3": "0011",
        "4": "0100",
        "5": "0101",
        "6": "0110",
        "7": "0111",
        "8": "1000",
        "9": "1001",
        "A": "1010",
        "B": "1011",
        "C": "1100",
        "D": "1101",
        "E": "1110",
        "F": "1111",
    }
    binString := ""
    for _, char := range s {
        binString += hex2bin[string(char)]
    }
    return binString
}

func packetsFromInput(lines []string) []*Packet {
    packets := make([]*Packet, 0)
    for _, line := range lines {
        binStr := expandToBin(line)
        packet := NewPacket(line, binStr)
        packets = append(packets, packet)
    }
    return packets
}

func main() {
    result := part1()
    logger.Logs.Infof("Part one result: %d", result)
    result = part2()
    logger.Logs.Infof("Part two result: %d", result)
}

func part1() int {
    lines := reader.LinesFromFile("input.txt")
    packets := packetsFromInput(lines)
    for _, packet := range packets {
        packet.Parse()
    }
    return packets[0].VersionSum()
}

func part2() int {
    lines := reader.LinesFromFile("input.txt")
    packets := packetsFromInput(lines)
    for _, packet := range packets {
        packet.Parse()
    }
    return packets[0].Resolve()
}
