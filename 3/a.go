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

func hasSymbolAdjacent(number partNumber, symbols map[int]map[int]bool) bool {
	var startRow = max(0, number.row-1)
	var endRow = number.row + 1
	var startCol = max(0, number.colStart-1)
	var endCol = number.colStart + number.length

	for r := startRow; r <= endRow; r++ {
		for c := startCol; c <= endCol; c++ {
			_, ok := symbols[r][c]
			if ok {
				return true
			}
		}
	}
	return false
}

func main() {
	var symbols = make(map[int]map[int]bool)
	var numbers []partNumber

	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var rowIndex = 0
	var currentNumber partNumber
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		symbols[rowIndex] = make(map[int]bool)
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
					numbers = append(numbers, currentNumber)
					currentNumber = partNumber{
						number:   "",
						row:      0,
						colStart: 0,
						length:   0,
					}
				}
				if char != '.' {
					symbols[rowIndex][col] = true
				}
			}
		}
		if currentNumber.length > 0 {
			numbers = append(numbers, currentNumber)
			currentNumber = partNumber{
				number:   "",
				row:      0,
				colStart: 0,
				length:   0,
			}
		}
		rowIndex = rowIndex + 1
	}

	fmt.Printf("Numbers: %+v\n", numbers)
	fmt.Printf("Symbols: %+v\n", symbols)

	var total int64 = 0
	for _, number := range numbers {
		if !hasSymbolAdjacent(number, symbols) {
			fmt.Printf("Invalid partnumber: %s\n", number.number)
		} else {
			n, _ := strconv.ParseInt(number.number, 10, 64)
			fmt.Printf("Valid partnumber: %d\n", n)
			total = total + n
		}
	}
	fmt.Printf("Sum of partnumbers: %d\n", total)
}
