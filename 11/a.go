package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

type Galaxy struct {
	x int
	y int
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	galaxies := []*Galaxy{}
	y := 0
	maxX := 0
	usedX := map[int]int{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Index(line, "#") == -1 {
			y = y + 2
			continue
		}

		nodes := strings.Split(line, "")
		maxX = max(maxX, len(nodes)-1)

		for x, node := range nodes {
			if node == "#" {
				galaxies = append(galaxies, &Galaxy{
					x: x,
					y: y,
				})
				usedX[x] = x
			}
		}
		y = y + 1
	}

	realX := 0
	for x := 0; x <= maxX; x++ {
		if _, ok := usedX[x]; !ok {
			realX = realX + 1
		}
		usedX[x] = realX
		realX = realX + 1
	}
	fmt.Printf("%+v\n", usedX)

	for i, galaxy := range galaxies {
		galaxy.x = usedX[galaxy.x]
		fmt.Printf("Galaxy %d: %+v\n", i+1, galaxy)
	}

	l := len(galaxies)
	total := int64(0)
	for i, g1 := range galaxies {
		for k := i + 1; k < l; k++ {
			g2 := galaxies[k]

			distance := int64(math.Abs(float64(g1.x-g2.x)) + math.Abs(float64(g1.y-g2.y)))
			total = total + distance

			fmt.Printf("%d, %d => %+v %+v => %d\n", i+1, k+1, g1, g2, distance)
		}
	}
	fmt.Printf("Total: %d\n", total)
}
