package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func replaceLast(text string) string {
	replacements := []string{
		"one", "1",
		"two", "2",
		"three", "3",
		"four", "4",
		"five", "5",
		"six", "6",
		"seven", "7",
		"eight", "8",
		"nine", "9",
	}

	last := 0
	lastReplacement := ""
	lastSearch := ""
	for i := 0; i < len(replacements); i = i + 1 {
		if i%2 == 1 {
			continue
		}
		search := replacements[i]
		found := strings.LastIndex(text, search)
		if found > 0 && found > last {
			last = found
			lastSearch = search
			lastReplacement = replacements[i+1]
		}
	}
	if last > 0 {
		return text[:last] + lastReplacement + text[last+len(lastSearch):]
	}
	return text
}

func firstDigit(text string) byte {
	index := strings.IndexAny(text, "1234567890")
	if index >= 0 {
		return text[index]
	}
	return 0x0
}

func lastDigit(text string) byte {
	index := strings.LastIndexAny(text, "1234567890")
	if index >= 0 {
		return text[index]
	}
	return 0x0
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	replacer := strings.NewReplacer(
		"one", "1",
		"two", "2",
		"three", "3",
		"four", "4",
		"five", "5",
		"six", "6",
		"seven", "7",
		"eight", "8",
		"nine", "9",
	)

	scanner := bufio.NewScanner(file)
	var total uint64 = 0
	for scanner.Scan() {
		var start byte = 0x0
		var end byte = 0x0
		orig := scanner.Text()
		start = firstDigit(replacer.Replace(orig))
		end = lastDigit(replaceLast(orig))
		end2 := lastDigit(replacer.Replace(orig))
		if end != end2 {
			fmt.Printf("difference: %s (%s => %s)\n", orig, string([]byte{start, end2}), string([]byte{start, end}))
			//break
		}

		str := string([]byte{start, end})
		value, _ := strconv.ParseUint(str, 10, 64)
		//fmt.Printf("%s => %d\n", orig, value)
		total = total + value
	}
	fmt.Printf("Total: %d\n", total)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
