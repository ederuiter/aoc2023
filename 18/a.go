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
	X int
	Y int
}

type Direction uint8

const (
	U Direction = 'U'
	R           = 'R'
	D           = 'D'
	L           = 'L'
)

var movement = map[Direction]Point{
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

func main() {
	file, err := os.Open("test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	prev := Point{0, 0}
	points := []Point{prev}
	scanner := bufio.NewScanner(file)
	minX := points[0].X
	minY := points[0].Y
	length := 0
	prevDirection := Direction('0')
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		direction, numStr, _ := Direction(line[0][0]), line[1], strings.Trim(line[2], "()")
		if direction == prevDirection || prevDirection == opposite[direction] {
			panic("Not simple!")
		}
		num, _ := strconv.ParseInt(numStr, 10, 64)
		length += int(num)
		current := Point{prev.X + (int(num) * movement[direction].X), prev.Y + (int(num) * movement[direction].Y)}
		minX = min(minX, current.X)
		minY = min(minY, current.Y)
		points = append(points, current)
		prev = current
		prevDirection = direction
	}

	points = append(points, points[1])

	fmt.Printf("%+v\n", points)
	//
	//minX -= 1
	//minY -= 1
	//for i, _ := range points {
	//	points[i].X -= minX
	//	points[i].Y -= minY
	//}

	fmt.Printf("%+v\n", points)

	area := 0
	n := len(points) - 1
	for i := 0; i < n; i++ {
		a := (points[i].X * points[i+1].Y) - (points[i+1].X * points[i].Y)
		fmt.Printf("[%d %d]:[%d %d] => %d\n", points[i].X, points[i].Y, points[i+1].X, points[i+1].Y, a)
		area += a
	}

	fmt.Printf("Area: %d Length: %d\n", area/2, length)
}
