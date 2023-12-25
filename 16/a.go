package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Direction int

const (
	N Direction = iota
	E
	S
	W
)

type Grid = [][]byte

type LightPos struct {
	Y         int
	X         int
	Direction Direction
}

var step = map[Direction]LightPos{
	N: {Y: -1, X: 0, Direction: N},
	E: {Y: 0, X: 1, Direction: E},
	S: {Y: 1, X: 0, Direction: S},
	W: {Y: 0, X: -1, Direction: W},
}

var oppositeDirection = map[Direction]Direction{
	N: S,
	S: N,
	E: W,
	W: E,
}

var reflections = map[byte]map[Direction][]Direction{
	'/': {
		N: {E},
		S: {W},
		E: {N},
		W: {S},
	},
	'\\': {
		N: {W},
		S: {E},
		E: {S},
		W: {N},
	},
	'-': {
		N: {E, W},
		S: {E, W},
		E: {},
		W: {},
	},
	'|': {
		E: {N, S},
		W: {N, S},
		N: {},
		S: {},
	},
}

func Init2D[T any](n, m int) [][]T {
	matrix := make([][]T, n)
	rows := make([]T, n*m)
	for i, startRow := 0, 0; i < n; i, startRow = i+1, startRow+m {
		endRow := startRow + m
		matrix[i] = rows[startRow:endRow:endRow]
	}
	return matrix
}

func Energize(mirrors Grid, start LightPos) int {
	stack := []LightPos{}
	current := start
	numEnergized := 0
	flow := Init2D[[]bool](len(mirrors), len(mirrors[0]))
	energized := Init2D[bool](len(mirrors), len(mirrors[0]))

	for {
		// walk this path in the current direction until we get to a mirror
		for current.Y >= 0 && current.Y < len(mirrors) && current.X >= 0 && current.X < len(mirrors[current.Y]) {
			//fmt.Printf("Currentpos: [%d, %d] => %d\n", current.Y, current.X, current.Direction)
			if len(flow[current.Y][current.X]) == 0 {
				flow[current.Y][current.X] = []bool{false, false, false, false}
			}

			if flow[current.Y][current.X][oppositeDirection[current.Direction]] {
				break
			}
			if !energized[current.Y][current.X] {
				numEnergized = numEnergized + 1
				energized[current.Y][current.X] = true
			}

			flow[current.Y][current.X][oppositeDirection[current.Direction]] = true
			mirror := mirrors[current.Y][current.X]
			if mirror != '.' && len(reflections[mirror][current.Direction]) > 0 {
				for _, reflectionDirection := range reflections[mirror][current.Direction] {
					next := LightPos{
						Y:         current.Y + step[reflectionDirection].Y,
						X:         current.X + step[reflectionDirection].X,
						Direction: reflectionDirection,
					}
					//fmt.Printf("[%d, %d] => %d\n", next.Y, next.X, reflectionDirection)
					stack = append(stack, next)
					flow[current.Y][current.X][reflectionDirection] = true
				}
				break
			} else {
				flow[current.Y][current.X][current.Direction] = true
			}
			current.Y = current.Y + step[current.Direction].Y
			current.X = current.X + step[current.Direction].X
		}

		if len(stack) == 0 {
			break
		}
		current, stack = stack[0], stack[1:]
	}
	return numEnergized
}

func MaximizeEnergy(mirrors Grid, flow [][][]bool, energized [][]bool) int {
	startPositions := map[Direction]LightPos{
		N: {Y: 0, X: 0, Direction: S},
		S: {Y: len(mirrors) - 1, X: 0, Direction: N},
		E: {Y: 0, X: len(mirrors[0]) - 1, Direction: W},
		W: {Y: 0, X: 0, Direction: E},
	}
	sideStep := map[Direction]LightPos{
		N: {Y: 0, X: 1, Direction: N},
		S: {Y: 0, X: 1, Direction: S},
		E: {Y: 1, X: 0, Direction: E},
		W: {Y: 1, X: 0, Direction: W},
	}

	numEnergized := 0
	for side, start := range startPositions {
		current := start
		for current.Y >= 0 && current.Y < len(mirrors) && current.X >= 0 && current.X < len(mirrors[current.Y]) {
			numCurrent := Energize(mirrors, current)
			numEnergized = max(numEnergized, numCurrent)
			current.Y = current.Y + sideStep[side].Y
			current.X = current.X + sideStep[side].X
		}
	}
	return numEnergized
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	mirrors := Grid{}
	flow := [][][]bool{}
	energized := [][]bool{}
	scanner := bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		mirrors = append(mirrors, []byte{})
		for _, c := range line {
			mirrors[y] = append(mirrors[y], byte(c))
		}
		y++
	}

	numEnergized := Energize(mirrors, LightPos{
		Y:         0,
		X:         0,
		Direction: E,
	})

	fmt.Printf("Energized: %d\n", numEnergized)
	fmt.Printf("Max energized: %d\n", MaximizeEnergy(mirrors, flow, energized))
}
