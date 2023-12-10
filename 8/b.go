package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	name      string
	leftName  string
	left      *Node
	rightName string
	right     *Node
}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(integers []int) int {
	a, b, integers := integers[0], integers[1], integers[2:]
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM([]int{result, integers[i]})
	}

	return result
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

	ends := [][]int{}
	var newNode *Node
	for currentNodeIndex, currentNode := range start {
		s := currentNode
		res := "S"
		ends = append(ends, []int{})
		done := false
		i := 0
		restart := false
		for !done {
			i = i + 1
			for m, move := range moves {
				if move == "L" {
					newNode = currentNode.left
				} else {
					newNode = currentNode.right
				}

				if m == 0 {
					if restart {
						if s == newNode {
							done = true
							break
						} else {
							s = newNode
							restart = false
						}
					}
					res = res + "[" + newNode.name + "]"
				}
				if newNode.name[2] != 'Z' {
					res = res + "."
				} else {
					ends[currentNodeIndex] = append(ends[currentNodeIndex], i)
					res = res + "Z[" + newNode.name + "][" + strconv.FormatInt(int64(i), 10) + "]"
					restart = true
				}
				currentNode = newNode
			}
			//fmt.Println(res)
			res = " "
		}
		//fmt.Println(res)
	}
	fmt.Printf("%+v\n", ends)

	indexes := []int{}
	for _, end := range ends {
		indexes = append(indexes, end[0])
	}

	result := LCM(indexes) * len(moves)
	fmt.Printf("Result: %d\n", result)

	//
	//
	//
	//numMoves := int64(0)
	//numIterations := int64(0)
	//isEnd := false
	//nodes := start
	//for !isEnd {
	//	for _, move := range moves {
	//		isEnd = true
	//		newNodes := []*Node{}
	//		for _, currentNode := range nodes {
	//			if move == "L" {
	//				newNode = currentNode.left
	//			} else {
	//				newNode = currentNode.right
	//			}
	//			if newNode.name[2] != 'Z' {
	//				isEnd = false
	//			}
	//			newNodes = append(newNodes, newNode)
	//		}
	//		numMoves = numMoves + 1
	//		nodes = newNodes
	//		if isEnd {
	//			break
	//		}
	//	}
	//	//numIterations = numIterations + 1
	//	//if numIterations%100000 == 0 {
	//	//	fmt.Printf(" => next iteration (%d, %d)\n", numIterations, numMoves)
	//	//}
	//}
	//
	//fmt.Printf("target reached in %d moves (%d)\n", numMoves, numIterations)
}
