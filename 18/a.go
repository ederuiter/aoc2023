package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	X         int
	Y         int
	Direction Direction
	Num       int
	Color     int
}

type Move struct {
	X int
	Y int
}

type Direction string

const (
	U Direction = "U"
	R           = "R"
	D           = "D"
	L           = "L"
)

var movement = map[Direction]Move{
	U: {0, -1},
	D: {0, 1},
	L: {-1, 0},
	R: {1, 0},
}

var opposite = map[Direction]Direction{
	U: D,
	L: R,
	R: L,
	D: U,
}

/*   --#######----############.
 *   --#.....######..........#.
 *   ..###---------#######---#.
 *   ....###########.....#####.
 */

func floodFill(startY int, startX int, color int, grid map[int]map[int]int, maxY int, maxX int) map[int]map[int]int {
	stack := []*Move{{startX, startY}}
	var move *Move
	for len(stack) > 0 {
		move, stack = stack[0], stack[1:]
		if _, ok := grid[move.Y][move.X]; ok {
			continue
		}

		x := move.X
		y := move.Y
		grid[y][x] = color
		//fmt.Printf("[%d, %d]\n", x, y)

		if x > 0 && grid[y][x-1] != color {
			stack = append(stack, &Move{x - 1, y})
		}
		if y > 0 && grid[y-1][x] != color {
			stack = append(stack, &Move{x, y - 1})
		}
		if x < maxX && grid[y][x+1] != color {
			stack = append(stack, &Move{x + 1, y})
		}
		if y < maxY && grid[y+1][x] != color {
			stack = append(stack, &Move{x, y + 1})
		}
	}
	return grid
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	prev := &Point{0, 0, L, 0, 0}
	points := []*Point{prev}
	scanner := bufio.NewScanner(file)
	minX := points[0].X
	minY := points[0].Y
	length := 0
	prevDirection := Direction('0')
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		direction, numStr, colorStr := Direction(line[0][0]), line[1], strings.Trim(line[2], "()#")
		color, _ := strconv.ParseInt(colorStr[:5], 16, 64)
		if direction == prevDirection || prevDirection == opposite[direction] {
			panic("Not simple!")
		}
		num, _ := strconv.ParseInt(numStr, 10, 64)
		length += int(num)
		current := &Point{prev.X + (int(num) * movement[direction].X), prev.Y + (int(num) * movement[direction].Y), U, 0, 0}
		prev.Num = int(num)
		prev.Direction = direction
		prev.Color = int(color)

		minX = min(minX, current.X)
		minY = min(minY, current.Y)
		points = append(points, current)
		prev = current
		prevDirection = direction
	}

	//points = append(points, points[1])

	fmt.Printf("%+v\n", points)

	minX -= 1
	minY -= 1
	maxX := 0
	maxY := 0
	for i, _ := range points {
		points[i].X -= minX
		points[i].Y -= minY
		maxX = max(maxX, points[i].X)
		maxY = max(maxY, points[i].Y)
		fmt.Printf("%+v\n", points[i])
	}
	maxX++
	maxY++

	//
	//fmt.Printf("%+v\n", points)

	grid := map[int]map[int]int{}
	for i := 0; i <= maxY; i++ {
		grid[i] = map[int]int{}
	}

	for _, point := range points {
		grid[point.Y][point.X] = point.Color
		for i := 1; i < point.Num; i++ {
			grid[point.Y+i*movement[point.Direction].Y][point.X+i*movement[point.Direction].X] = point.Color
		}
	}

	grid = floodFill(0, 0, 1, grid, maxY, maxX)

	numInside := 0
	for y := 0; y <= maxY; y++ {
		line := ""
		for x := 0; x <= maxX; x++ {
			if _, ok := grid[y][x]; ok {
				line += "#"
			} else {
				line += " "
				numInside++
			}
		}
		fmt.Println(line)
	}
	fmt.Printf("Inside: %d trench: %d lave: %d\n", numInside, length, numInside+length)
}
