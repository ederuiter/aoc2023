package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func getReflections(lines []string) int {
	left := ""
	lineLength := len(lines[0])
	numLines := len(lines)
	for i, line := range lines {
		num := min(i+1, numLines-(i+1))
		if num == 0 {
			break
		}
		left = left + line
		right := ""
		for k := i + 1; k <= i+num; k++ {
			right = lines[k] + right
		}

		if len(left) > num*lineLength {
			left = left[len(left)-(num*lineLength):]
		}

		if left == right {
			return i + 1
		}
	}

	return 0
}

func calc(rows []string, cols []string) int64 {
	reflectionX := getReflections(cols)
	reflectionY := getReflections(rows)

	if reflectionX > 0 {
		fmt.Printf("  %s><\n", strings.Repeat(" ", reflectionX-1))
	}

	for y, line := range rows {
		match := "  "
		if reflectionY > 0 && y == reflectionY-1 {
			match = "v "
		} else if reflectionY > 0 && y == reflectionY {
			match = "^ "
		}
		fmt.Printf("%s%s\n", match, line)
	}

	res := int64((100 * reflectionY) + reflectionX)
	fmt.Printf(" => %d\n\n", res)
	return res
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	total := int64(0)
	scanner := bufio.NewScanner(file)
	rows := []string{}
	cols := []string{}
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			total = total + calc(rows, cols)
			y = 0
			rows = []string{}
			cols = []string{}
		} else {
			rows = append(rows, "")
			for x, cell := range strings.Split(line, "") {
				if y == 0 {
					cols = append(cols, "")
				}
				rows[y] = rows[y] + cell
				cols[x] = cols[x] + cell
			}
			y++
		}
	}
	total = total + calc(rows, cols)

	fmt.Printf("Reflection sum: %d\n", total)
}
