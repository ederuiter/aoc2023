package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Node struct {
	Type      rune
	x         int
	y         int
	isLoop    bool
	isOutside bool
	isPad     bool
}

type Grid [][]*Node

func getNeighbours(node *Node, grid Grid) (*Node, *Node) {
	var n1 *Node
	var n2 *Node
	switch node.Type {
	case 'F':
		n1 = grid[node.y][node.x+1]
		n2 = grid[node.y+1][node.x]
		break
	case 'L':
		n1 = grid[node.y-1][node.x]
		n2 = grid[node.y][node.x+1]
		break
	case '|', '‖':
		n1 = grid[node.y-1][node.x]
		n2 = grid[node.y+1][node.x]
		break
	case '-', '=':
		n1 = grid[node.y][node.x-1]
		n2 = grid[node.y][node.x+1]
		break
	case 'J':
		n1 = grid[node.y-1][node.x]
		n2 = grid[node.y][node.x-1]
		break
	case '7':
		n1 = grid[node.y][node.x-1]
		n2 = grid[node.y+1][node.x]
		break
	case '.':
		n1 = nil
		n2 = nil
		break
	default:
		panic("Invalid type: " + string(node.Type))
	}

	return n1, n2
}

func getPadX(a *Node, b *Node) *Node {
	res := &Node{
		Type:      '.',
		x:         a.x,
		y:         a.y,
		isLoop:    false,
		isOutside: false,
		isPad:     true,
	}
	if a.isLoop && b != nil && b.isLoop &&
		(a.Type == 'F' || a.Type == 'L' || a.Type == '-') &&
		(b.Type == '7' || b.Type == 'J' || b.Type == '-') {
		res.Type = '-'
		res.isLoop = true
	}
	return res
}

func getPadY(a *Node, b *Node) *Node {
	res := &Node{
		Type:      '.',
		x:         a.x,
		y:         a.y,
		isLoop:    false,
		isOutside: false,
		isPad:     true,
	}
	if a.isLoop && b != nil && b.isLoop &&
		(a.Type == 'F' || a.Type == '7' || a.Type == '|') &&
		(b.Type == 'J' || b.Type == 'L' || b.Type == '|') {
		res.Type = '|'
		res.isLoop = true
	}
	return res
}

func padGrid(grid Grid) Grid {
	res := Grid{}
	for y, row := range grid {
		res = append(res, []*Node{})
		res = append(res, []*Node{})
		for x, node := range row {
			var nextX *Node
			var nextY *Node
			if x+1 < len(row) {
				nextX = row[x+1]
			}
			if y+1 < len(grid) {
				nextY = grid[y+1][x]
			}
			res[y*2] = append(res[y*2], node, getPadX(node, nextX))
			res[(y*2)+1] = append(res[(y*2)+1], getPadY(node, nextY), &Node{
				Type:      '.',
				x:         x,
				y:         y,
				isLoop:    false,
				isOutside: false,
				isPad:     true,
			})
			res[y*2][x*2].y = y * 2
			res[y*2][x*2].x = x * 2
			res[y*2][(x*2)+1].y = y * 2
			res[y*2][(x*2)+1].x = (x * 2) + 1
			res[(y*2)+1][x*2].y = (y * 2) + 1
			res[(y*2)+1][x*2].x = x * 2
			res[(y*2)+1][(x*2)+1].y = (y * 2) + 1
			res[(y*2)+1][(x*2)+1].x = (x * 2) + 1
		}
	}
	return res
}

func printGrid(grid Grid) {
	charMap := map[rune]string{
		'F': "┌",
		'J': "┘",
		'L': "└",
		'|': "│",
		'-': "─",
		'7': "┐",
		'‖': "┃",
		'=': "━",
	}
	res := ""
	for _, row := range grid {
		for _, node := range row {
			if node.isLoop {
				res = res + charMap[node.Type]
			} else if node.isOutside {
				res = res + "."
			} else {
				res = res + " "
			}
		}
		res = res + "\n"
	}
	fmt.Print(res)
}

func floodFill(y int, x int, grid Grid) {
	maxY := len(grid) - 1
	maxX := len(grid[y]) - 1
	stack := []*Node{grid[y][x]}

	var node *Node
	for len(stack) > 0 {
		node, stack = stack[0], stack[1:]
		if node.isOutside {
			continue
		}
		node.isOutside = true
		x := node.x
		y := node.y

		if x > 0 && !grid[y][x-1].isLoop && !grid[y][x-1].isOutside {
			stack = append(stack, grid[y][x-1])
		}
		if y > 0 && !grid[y-1][x].isLoop && !grid[y-1][x].isOutside {
			stack = append(stack, grid[y-1][x])
		}
		if x < maxX && !grid[y][x+1].isLoop && !grid[y][x+1].isOutside {
			stack = append(stack, grid[y][x+1])
		}
		if y < maxY && !grid[y+1][x].isLoop && !grid[y+1][x].isOutside {
			stack = append(stack, grid[y+1][x])
		}
	}
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	y := 1
	grid := Grid{[]*Node{}}
	var start *Node
	for scanner.Scan() {
		grid = append(grid, []*Node{&Node{
			Type:   '.',
			x:      0,
			y:      y,
			isLoop: false,
		}})
		line := strings.Split(scanner.Text(), "")
		x := 0
		for i, s := range line {
			x = i + 1
			grid[y] = append(grid[y], &Node{Type: rune(s[0]), x: x, y: y, isLoop: s[0] == 'S'})
			if s[0] == 'S' {
				start = grid[y][x]
			}
		}
		grid[y] = append(grid[y], &Node{
			Type:   '.',
			x:      x + 1,
			y:      y,
			isLoop: false,
		})
		y = y + 1
	}
	grid = append(grid, []*Node{})
	last := len(grid) - 1
	for x, _ := range grid[1] {
		grid[0] = append(grid[0], &Node{
			Type:   '.',
			x:      x,
			y:      0,
			isLoop: false,
		})
		grid[last] = append(grid[last], &Node{
			Type:   '.',
			x:      x,
			y:      last,
			isLoop: false,
		})
	}

	fmt.Printf("Start: %+v\n", start)

	loop := []*Node{start}
	start.Type = 'L' //yes this is cheating
	current := start
	prev := start
	for {
		var next *Node
		n1, n2 := getNeighbours(current, grid)
		//fmt.Printf("Neighbours: %+v %+v\n", n1, n2)
		if n1 != prev && n1 != nil {
			next = n1
		} else {
			next = n2
		}
		if next == start {
			break
		}
		next.isLoop = true
		loop = append(loop, next)
		//fmt.Printf("Loop %+v\n", loop)
		prev = current
		current = next
	}

	printGrid(grid)
	paddedGrid := padGrid(grid)
	printGrid(paddedGrid)
	floodFill(0, 0, paddedGrid)
	printGrid(paddedGrid)

	inside := 0
	for _, row := range paddedGrid {
		for _, node := range row {
			if node.isPad == false && node.isOutside == false && node.isLoop == false {
				inside = inside + 1
			}
		}
	}

	fmt.Printf("Loop length: %d max distance: %d\n", len(loop), len(loop)/2)
	fmt.Printf("Tiles inside the loop: %d\n", inside)
}
