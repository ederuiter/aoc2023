package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
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
	U Direction = "3"
	R           = "0"
	D           = "1"
	L           = "2"
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

func compactAxis(items []int) ([]int, map[int]int) {
	slices.Sort(items)
	items = slices.Compact(items)

	m := map[int]int{}
	for index, x := range items {
		m[x] = index
	}

	return items, m
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
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		colorStr := strings.Trim(line[2], "()#")
		num, _ := strconv.ParseInt(colorStr[:5], 16, 64)
		direction := Direction(colorStr[5])

		length += int(num)
		current := &Point{prev.X + (int(num) * movement[direction].X), prev.Y + (int(num) * movement[direction].Y), U, 0, 0}
		prev.Num = int(num)
		prev.Direction = direction
		prev.Color = 255

		minX = min(minX, current.X)
		minY = min(minY, current.Y)
		points = append(points, current)
		prev = current
	}

	//points = append(points, points[1])

	fmt.Printf("%+v\n", points)

	minX -= 1
	minY -= 1
	xItems := []int{0}
	yItems := []int{0}
	for i, _ := range points {
		points[i].X -= minX
		points[i].Y -= minY
		xItems = append(xItems, points[i].X)
		xItems = append(xItems, points[i].X-1)
		xItems = append(xItems, points[i].X+1)
		yItems = append(yItems, points[i].Y)
		yItems = append(yItems, points[i].Y-1)
		yItems = append(yItems, points[i].Y+1)
		fmt.Printf("%+v\n", points[i])
	}

	xItems, xMap := compactAxis(xItems)
	yItems, yMap := compactAxis(yItems)

	maxY := len(yItems)
	maxX := len(xItems)

	//
	//fmt.Printf("%+v\n", points)

	grid := map[int]map[int]int{}
	for i := 0; i <= maxY+1; i++ {
		grid[i] = map[int]int{}
	}

	for _, point := range points {
		grid[yMap[point.Y]][xMap[point.X]] = point.Color
		for i := 1; i < point.Num; i++ {
			newY := point.Y + i*movement[point.Direction].Y
			newX := point.X + i*movement[point.Direction].X
			if _, yOk := yMap[newY]; !yOk {
				continue
			}
			if _, xOk := xMap[newX]; !xOk {
				continue
			}
			grid[yMap[newY]][xMap[newX]] = point.Color
		}
	}

	grid = floodFill(0, 0, 1, grid, maxY+1, maxX+1)

	fmt.Printf("%d %d %d %d\n", len(yItems), len(xItems), maxY, maxX)
	numInside := 0
	for y := 0; y < maxY; y++ {
		line := ""
		lineInside := 0
		for x := 0; x < maxX; x++ {
			if _, ok := grid[y][x]; ok {
				line += "#"
			} else {
				line += " "
				lineInside += xItems[x+1] - xItems[x]
			}
		}
		if lineInside > 0 {
			numInside += (yItems[y+1] - yItems[y]) * lineInside
		}
		fmt.Println(line)
	}
	fmt.Printf("Inside: %d trench: %d lava: %d\n", numInside, length, numInside+length)
}
