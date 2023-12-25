package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sum := 0
	current := 0
	str := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		for _, chr := range line {
			if chr == ',' {
				fmt.Printf("%s => %d\n", str, current)
				sum = sum + current
				current = 0
				str = ""
			} else {
				current = ((current + int(chr)) * 17) % 256
				str = str + string(chr)
			}
		}
	}
	if str != "" {
		fmt.Printf("%s => %d\n", str, current)
		sum = sum + current
	}

	fmt.Printf("Sum of hashes: %d\n", sum)
}
