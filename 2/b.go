package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func split2(s string, sep string) (string, string) {
	parts := strings.SplitN(s, sep, 2)
	return parts[0], parts[1]
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var total int64 = 0
	for scanner.Scan() {
		text := scanner.Text()
		_, rest := split2(text, ": ")
		subsets := strings.Split(rest, "; ")

		needed := map[string]int64{
			"red":   0,
			"green": 0,
			"blue":  0,
		}

		for _, subset := range subsets {
			marbles := strings.Split(subset, ", ")
			for _, marble := range marbles {
				numStr, color := split2(marble, " ")
				num, _ := strconv.ParseInt(numStr, 10, 64)
				needed[color] = max(needed[color], num)
			}
		}
		fmt.Printf("%+v\n", needed)
		total = total + (needed["red"] * needed["green"] * needed["blue"])
	}
	fmt.Printf("Sum of possible games: %d\n", total)
}
