package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func generateBitPatterns(numbers []int, latitude int, start bool) []string {
	res := []string{}
	var number int
	if len(numbers) == 0 {
		return res
	} else if len(numbers) > 1 {
		number, numbers = numbers[0], numbers[1:]
	} else {
		number, numbers = numbers[0], []int{}
	}
	for i := latitude; i >= 0; i-- {
		left := latitude - i
		zeros := i
		if !start {
			zeros++
		}
		bits := strings.Repeat("0", zeros) + strings.Repeat("1", number)
		if left > 0 {
			for _, pattern := range generateBitPatterns(numbers, left, false) {
				res = append(res, bits+pattern)
			}
		} else {
			for _, n := range numbers {
				bits = bits + "0" + strings.Repeat("1", n)
			}
			res = append(res, bits)
		}
	}
	return res
}

func validPatterns(patterns []string, maskFixed string, maskGap string) int64 {
	num := int64(0)

	f, _ := strconv.ParseUint(maskFixed, 2, 64)
	g, _ := strconv.ParseUint(maskGap, 2, 64)
	for _, str := range patterns {
		pattern, _ := strconv.ParseUint(str, 2, 64)
		if pattern&f == f && pattern&g == 0 {
			num++
		}
	}
	return num
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	total := int64(0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		parts := strings.Split(strings.Trim(line[0], "."), ".")
		str := ""
		for _, part := range parts {
			if len(part) > 0 {
				str = str + part + "."
			}
		}

		fixedMask := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(str, "#", "1"), ".", "0"), "?", "0")
		gapMask := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(str, ".", "1"), "#", "0"), "?", "0")

		numbersStr := strings.Split(line[1], ",")
		numbers := []int{}
		sum := 0
		for _, number := range numbersStr {
			n, _ := strconv.ParseInt(number, 10, 64)
			sum = sum + int(n)
			numbers = append(numbers, int(n))
		}
		numbers = append(numbers, 0)

		latitude := len(str) - (sum + len(numbers) - 1)

		patterns := generateBitPatterns(numbers, latitude, true)

		fmt.Printf("%s %s => %d, %d\n", str, line[1], latitude, len(patterns))
		fmt.Printf("%s\n%s\n\n", fixedMask, gapMask)
		for _, p := range patterns {
			fmt.Printf("%s\n", p)
		}

		numValid := validPatterns(patterns, fixedMask, gapMask)
		fmt.Printf("=> %d valid patterns\n", numValid)
		total = total + numValid
	}
	fmt.Printf("Total %d combinations to check\n", total)
}
