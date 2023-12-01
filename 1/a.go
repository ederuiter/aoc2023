package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var total uint64 = 0
	for scanner.Scan() {
		var start byte = 0x0
		var end byte = 0x0
		text := scanner.Text()
		for _, b := range scanner.Bytes() {
			if unicode.IsDigit(rune(b)) {
				if start == 0x0 {
					start = b
					end = b
				} else {
					end = b
				}
			}
		}
		str := string([]byte{start, end})
		value, _ := strconv.ParseUint(str, 10, 64)
		fmt.Printf("%s => %d\n", text, value)
		total = total + value
	}
	fmt.Printf("Total: %d\n", total)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
