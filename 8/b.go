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
	start := []*Node{}

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

		if name[2] == 'A' {
			start = append(start, nodeMap[name])
		}
	}

	for _, node := range nodeMap {
		node.left = nodeMap[node.leftName]
		node.right = nodeMap[node.rightName]
	}

	var newNode *Node
	numMoves := int64(0)
	numIterations := int64(0)
	isEnd := false
	nodes := start
	for !isEnd {
		for _, move := range moves {
			isEnd = true
			newNodes := []*Node{}
			for _, currentNode := range nodes {
				if move == "L" {
					newNode = currentNode.left
				} else {
					newNode = currentNode.right
				}
				if newNode.name[2] != 'Z' {
					isEnd = false
				}
				newNodes = append(newNodes, newNode)
			}
			numMoves = numMoves + 1
			nodes = newNodes
			if isEnd {
				break
			}
		}
		numIterations = numIterations + 1
		if numIterations%100000 == 0 {
			fmt.Printf(" => next iteration (%d, %d)\n", numIterations, numMoves)
		}
	}

	fmt.Printf("target reached in %d moves (%d)\n", numMoves, numIterations)
}
