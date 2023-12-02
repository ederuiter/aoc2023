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

	available := map[string]int64{
		"red":   12,
		"green": 13,
		"blue":  14,
	}

	scanner := bufio.NewScanner(file)
	var total int64 = 0
	for scanner.Scan() {
		text := scanner.Text()
		game, rest := split2(text, ": ")
		_, gameIdStr := split2(game, " ")
		subsets := strings.Split(rest, "; ")
		possible := true
		for _, subset := range subsets {
			marbles := strings.Split(subset, ", ")
			for _, marble := range marbles {
				numStr, color := split2(marble, " ")
				num, _ := strconv.ParseInt(numStr, 10, 64)
				fmt.Printf("Checking %s %d > %d\n", color, num, available[color])
				if num > available[color] {
					possible = false
					break
				}
			}
			if !possible {
				break
			}
		}
		gameId, _ := strconv.ParseInt(gameIdStr, 10, 64)
		if possible {
			total = total + gameId
			fmt.Printf("Game %d is possible (%d): %s\n", gameId, total, text)
		} else {
			fmt.Printf("Game %d is impossible: %s\n", gameId, text)
		}
	}
	fmt.Printf("Sum of possible games: %d\n", total)
}
