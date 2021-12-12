package main

import (
    "fmt"
    "regexp"
    "strings"
    logger "advent2021/adventlogger"
    reader "advent2021/adventreader"
)

// Replace 'start' with '0' and 'end' with '1'; they could match node labels
var replacer = map[string]string{"start": "0", "end": "1"}

type Node struct {
    visited bool
    label string
    neighbors []*Node
}

func NewNode(label string) *Node {
    neighbors := make([]*Node, 0)
    node := Node{visited: false, label: label, neighbors: neighbors}
    return &node
}

func (n Node) String() string {
    return fmt.Sprintf(n.label)
}

func (n *Node) Walk(pathSoFar string) int {
    safeLabel := n.label
    if val, ok := replacer[safeLabel]; ok {
        safeLabel = val
    }
    if n.label == "end" {
        return 1
    }
    if n.IsLower() && strings.Contains(pathSoFar, safeLabel) {
        return 0
    }
    sum := 0
    pathSoFar += safeLabel
    for _, neigh := range n.neighbors {
        sum += neigh.Walk(pathSoFar)
    }
    return sum
}

type ShittyPath struct {
    path string
    visitedSmallTwice bool
}

func (n *Node) ShittyWalk(pathSoFar ShittyPath) int {
    safeLabel := n.label
    if val, ok := replacer[safeLabel]; ok {
        safeLabel = val
    }
    if n.label == "end" {
        return 1
    }
    if n.IsLower() && strings.Contains(pathSoFar.path, safeLabel) {
        if pathSoFar.visitedSmallTwice {
            // illegal path
            return 0
        }
        if n.label == "start" {
            // illegal path
            return 0
        }
        pathSoFar.visitedSmallTwice = true
    }
    sum := 0
    pathSoFar.path += safeLabel
    for _, neigh := range n.neighbors {
        sum += neigh.ShittyWalk(pathSoFar)
    }
    return sum
}

func (n *Node) Connect(other *Node) {
    n.neighbors = append(n.neighbors, other)
}

func (n *Node) IsLower() bool {
    return strings.ToLower(n.label) == n.label
}

type NodeMap struct {
    nodeMap map[string]*Node
}

func NewNodeMap() *NodeMap {
    nodeMap := make(map[string]*Node)
    return &NodeMap{nodeMap: nodeMap}
}

func (nm NodeMap) Print() {
    for _, node := range nm.nodeMap {
        neighbors := make([]string, 0)
        for _, neigh := range node.neighbors {
            neighbors = append(neighbors, neigh.label)
        }
        fmt.Printf("%s neighbors: %s\n", node, strings.Join(neighbors, ","))
    }
}


func LinesToNodes(lines []string) *NodeMap {
    nodeMap := NewNodeMap()
    for _, line := range lines {
        tokens := regexp.MustCompile("-").Split(line, -1)
        first, second := tokens[0], tokens[1]
        nodeA := NewNode(first)
        nodeB := NewNode(second)
        if node, ok := nodeMap.nodeMap[first]; ok {
            nodeA = node
        } else {
            nodeMap.nodeMap[first] = nodeA
        }
        if node, ok := nodeMap.nodeMap[second]; ok {
            nodeB = node
        } else {
            nodeMap.nodeMap[second] = nodeB
        }
        nodeA.Connect(nodeB)
        nodeB.Connect(nodeA)
    }
    return nodeMap
}

func main() {
    result := part1()
    logger.Logs.Infof("Part one result: %d", result)
    result = part2()
    logger.Logs.Infof("Part two result: %d", result)
}

func part1() int {
    lines := reader.LinesFromFile("input.txt")
    nodeMap := LinesToNodes(lines)
    // nodeMap.Print()
    numPaths := nodeMap.nodeMap["start"].Walk("")
    return numPaths
}

func part2() int {
    lines := reader.LinesFromFile("input.txt")
    nodeMap := LinesToNodes(lines)
    numPaths := nodeMap.nodeMap["start"].ShittyWalk(ShittyPath{"", false})
    return numPaths
}
