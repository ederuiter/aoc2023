package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"
)

type partNumber struct {
	number   string
	row      int
	colStart int
	length   int
}

type gridPos struct {
	row int
	col int
}

func getAdjacentParts(row int, col int, grid map[int]map[int]partNumber) []partNumber {
	var startRow = max(0, row-1)
	var endRow = row + 1
	var startCol = max(0, col-1)
	var endCol = col + 1

	res := []partNumber{}
	for r := startRow; r <= endRow; r++ {
		for c := startCol; c <= endCol; c++ {
			part, ok := grid[r][c]
			if ok {
				c = part.colStart + part.length
				res = append(res, part)
			}
		}
	}
	return res
}

func addToGrid(grid map[int]map[int]partNumber, part partNumber) map[int]map[int]partNumber {
	for c := part.colStart; c < part.colStart+part.length; c++ {
		grid[part.row][c] = part
	}
	return grid
}

func main() {
	var grid = make(map[int]map[int]partNumber)
	var ratios []gridPos

	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var rowIndex = 0
	var currentNumber partNumber
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		grid[rowIndex] = make(map[int]partNumber)
		row := []rune(scanner.Text())
		for col, char := range row {
			if unicode.IsDigit(char) {
				currentNumber.number = currentNumber.number + string(char)
				currentNumber.length = currentNumber.length + 1
				if currentNumber.length == 1 {
					currentNumber.colStart = col
					currentNumber.row = rowIndex
				}
			} else {
				if currentNumber.length > 0 {
					grid = addToGrid(grid, currentNumber)
					currentNumber = partNumber{
						number:   "",
						row:      0,
						colStart: 0,
						length:   0,
					}
				}
				if char == '*' {
					ratios = append(ratios, gridPos{
						row: rowIndex,
						col: col,
					})
				}
			}
		}
		if currentNumber.length > 0 {
			grid = addToGrid(grid, currentNumber)
			currentNumber = partNumber{
				number:   "",
				row:      0,
				colStart: 0,
				length:   0,
			}
		}
		rowIndex = rowIndex + 1
	}

	fmt.Printf("Ratios: %+v\n", ratios)
	fmt.Printf("Grid: %+v\n", grid)

	var total int64 = 0
	for _, ratio := range ratios {
		parts := getAdjacentParts(ratio.row, ratio.col, grid)
		if len(parts) == 2 {
			n1, _ := strconv.ParseInt(parts[0].number, 10, 64)
			n2, _ := strconv.ParseInt(parts[1].number, 10, 64)
			r := n1 * n2
			total = total + r
			fmt.Printf("Valid ratio at (%d, %d) => %d * %d = %d\n", ratio.row, ratio.col, n1, n2, r)
		} else {
			fmt.Printf("Invalid ratio at (%d, %d) => %d parts\n", ratio.row, ratio.col, len(parts))
		}
	}
	fmt.Printf("Sum of partnumbers: %d\n", total)
}
