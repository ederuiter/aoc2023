package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Direction int

const (
	N Direction = iota
	S
	W
	E
)

type Grid struct {
	NumX int
	NumY int
	Load int
	Data [][]string
}

func (g Grid) ToString() string {
	res := ""
	for _, line := range g.Data {
		res = res + strings.Join(line, "")
	}
	return res
}

type GridPos struct {
	X int
	Y int
}

var tiltStart = map[Direction]GridPos{
	N: {X: 0, Y: 0},
	W: {X: 0, Y: 0},
	S: {X: 0, Y: 1},
	E: {X: 1, Y: 0},
}

var tiltNextLine = map[Direction]GridPos{
	N: {X: 1, Y: 0},
	W: {X: 0, Y: 1},
	S: {X: 1, Y: 0},
	E: {X: 0, Y: 1},
}

var tiltNextPos = map[Direction]GridPos{
	N: {X: 0, Y: 1},
	W: {X: 1, Y: 0},
	S: {X: 0, Y: -1},
	E: {X: -1, Y: 0},
}

var cacheMap = map[string]int{}

func Tilt(grid Grid, direction Direction) Grid {
	pos := GridPos{
		X: tiltStart[direction].X * (grid.NumX - 1),
		Y: tiltStart[direction].Y * (grid.NumY - 1),
	}
	grid.Load = 0
	for pos.X >= 0 && pos.X < grid.NumX && pos.Y >= 0 && pos.Y < grid.NumY {
		startLinePos := pos
		nextFree := pos
		for pos.X >= 0 && pos.X < grid.NumX && pos.Y >= 0 && pos.Y < grid.NumY {
			cur := grid.Data[pos.Y][pos.X]
			if cur == "O" {
				if nextFree != pos {
					grid.Data[nextFree.Y][nextFree.X] = cur
					grid.Data[pos.Y][pos.X] = "."
				}
				grid.Load = grid.Load + (grid.NumY - nextFree.Y)

				nextFree.X = nextFree.X + tiltNextPos[direction].X
				nextFree.Y = nextFree.Y + tiltNextPos[direction].Y
			}

			pos.X = pos.X + tiltNextPos[direction].X
			pos.Y = pos.Y + tiltNextPos[direction].Y
			if cur == "#" {
				nextFree = pos
			}
		}

		pos.X = startLinePos.X + tiltNextLine[direction].X
		pos.Y = startLinePos.Y + tiltNextLine[direction].Y
	}
	return grid
}

func TiltCycle(grid Grid, cycle int) (Grid, int) {
	cacheKey := grid.ToString()
	cycleTime := 0
	if cache, ok := cacheMap[cacheKey]; ok {
		//fmt.Printf("Cycle detected at %d; original cached at: %d\n", cycle, cache)
		cycleTime = cycle - cache
	} else {
		cacheMap[cacheKey] = cycle
	}

	grid = Tilt(grid, N)
	//fmt.Printf("N =>\n")
	//printGrid(grid)
	grid = Tilt(grid, W)
	//fmt.Printf("W =>\n")
	//printGrid(grid)
	grid = Tilt(grid, S)
	//fmt.Printf("S =>\n")
	//printGrid(grid)
	grid = Tilt(grid, E)
	//fmt.Printf("E =>\n")
	//printGrid(grid)

	return grid, cycleTime
}

func printGrid(grid Grid) {
	fmt.Printf("┌%s┐\n", strings.Repeat("─", grid.NumX))
	for _, line := range grid.Data {
		fmt.Printf("│%s│\n", strings.Join(line, ""))
	}
	fmt.Printf("└%s┘\n", strings.Repeat("─", grid.NumX))
	fmt.Printf("Total load: %d\n", grid.Load)
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//total := int64(0)
	grid := Grid{
		Data: [][]string{},
		NumX: 0,
		NumY: 0,
	}
	scanner := bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "")
		grid.Data = append(grid.Data, []string{})
		if y == 0 {
			grid.NumX = len(line)
		}
		for _, c := range line {
			grid.Data[y] = append(grid.Data[y], c)
		}
		y++
	}
	grid.NumY = y

	tilted := Tilt(grid, N)
	fmt.Printf("Part 1:\n")
	printGrid(tilted)

	cycleTime := 0
	numIterations := 1_000_000_000
	skipped := false
	for i := 0; i < numIterations; i++ {
		grid, cycleTime = TiltCycle(grid, i)
		if cycleTime > 0 && !skipped {
			i = numIterations - ((numIterations - i) % cycleTime)
			skipped = true
		}
	}

	fmt.Printf("Part 2:\n")
	printGrid(grid)
}
