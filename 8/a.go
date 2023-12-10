package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Node struct {
	name      string
	leftName  string
	left      *Node
	rightName string
	right     *Node
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	nodeMap := map[string]*Node{}
	moves := []string{}
	for scanner.Scan() {
		if len(moves) == 0 {
			moves = strings.Split(scanner.Text(), "")
			continue
		}

		if scanner.Text() == "" {
			continue
		}

		parts := strings.Split(scanner.Text(), " = ")
		name := parts[0]
		lr := parts[1]
		lrParts := strings.Split(lr[1:len(lr)-1], ", ")
		nodeMap[name] = &Node{
			name:      name,
			leftName:  lrParts[0],
			left:      nil,
			rightName: lrParts[1],
			right:     nil,
		}
	}

	for _, node := range nodeMap {
		node.left = nodeMap[node.leftName]
		node.right = nodeMap[node.rightName]
	}

	currentNode := nodeMap["AAA"]
	numMoves := 0
	numIterations := 0
	for currentNode.name != "ZZZ" {
		for _, move := range moves {
			if currentNode.name == "ZZZ" {
				break
			}
			if move == "L" {
				currentNode = currentNode.left
			} else {
				currentNode = currentNode.right
			}
			numMoves = numMoves + 1
		}
		numIterations = numIterations + 1
	}

	fmt.Printf("target reached in %d moves (%d iterations)\n", numMoves, numIterations)
}
