package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func generateBitPatterns(numbers []int, latitude int, maxZeros int, maskFixed string, maskGap string, start bool) int64 {
	res := int64(0)
	var number int
	if len(numbers) == 0 {
		return res
	} else if len(numbers) > 1 {
		number, numbers = numbers[0], numbers[1:]
	} else {
		number, numbers = numbers[0], []int{}
	}
	for i := min(maxZeros, latitude); i >= 0; i-- {
		left := latitude - i
		zeros := i
		if !start {
			zeros++
		}
		bits := strings.Repeat("0", zeros) + strings.Repeat("1", number)
		if len(bits) > 0 && !validatePattern(bits, maskFixed[0:len(bits)], maskGap[0:len(bits)]) {
			continue
		}
		if left > 0 {
			res = res + generateBitPatterns(numbers, left, maxZeros, maskFixed[len(bits):], maskGap[len(bits):], false)
		} else {
			for _, n := range numbers {
				bits = bits + "0" + strings.Repeat("1", n)
			}
			if validatePattern(bits, maskFixed[0:len(bits)], maskGap[0:len(bits)]) {
				res++
			}
		}
	}
	return res
}

func validatePattern(str string, maskFixed string, maskGap string) bool {
	//fmt.Printf("validate: \n%s\n%s\n%s\n", str, maskFixed, maskGap)
	f, _ := strconv.ParseUint(maskFixed, 2, 64)
	g, _ := strconv.ParseUint(maskGap, 2, 64)
	pattern, _ := strconv.ParseUint(str, 2, 64)
	return pattern&f == f && pattern&g == 0
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
		unfolded := line[0] + "?" + line[0] + "?" + line[0] + "?" + line[0] + "?" + line[0]
		numbersUnfolded := line[1] + "," + line[1] + "," + line[1] + "," + line[1] + "," + line[1]
		parts := strings.Split(strings.Trim(unfolded, "."), ".")
		str := ""
		for _, part := range parts {
			if len(part) > 0 {
				str = str + part + "."
			}
		}

		fixedMask := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(str, "#", "1"), ".", "0"), "?", "0")
		gapMask := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(str, ".", "1"), "#", "0"), "?", "0")

		prevFixed := strings.Index(fixedMask, "1")
		maxZeros := prevFixed
		index := prevFixed
		for prevFixed >= 0 {
			nextFixed := strings.Index(fixedMask[index+1:], "1")
			if nextFixed >= 0 {
				maxZeros = max(maxZeros, nextFixed)
				index = index + nextFixed + 1
			}
			prevFixed = nextFixed
		}
		if maxZeros < 0 {
			maxZeros = len(fixedMask)
		}

		numbersStr := strings.Split(numbersUnfolded, ",")
		numbers := []int{}
		sum := 0
		for _, number := range numbersStr {
			n, _ := strconv.ParseInt(number, 10, 64)
			sum = sum + int(n)
			numbers = append(numbers, int(n))
		}
		numbers = append(numbers, 0)

		latitude := len(str) - (sum + len(numbers) - 1)

		fmt.Printf("%s %s => %d %d %d\n", str, numbersUnfolded, len(str), latitude, maxZeros)
		patterns := generateBitPatterns(numbers, latitude, maxZeros, fixedMask, gapMask, true)
		fmt.Printf(" => %d\n", patterns)

		//fmt.Printf("%s\n%s\n\n", fixedMask, gapMask)
		//for _, p := range patterns {
		//	fmt.Printf("%s\n", p)
		//}

		//numValid := validPatterns(patterns, fixedMask, gapMask)
		//fmt.Printf("=> %d valid patterns\n", numValid)
		total = total + patterns
	}
	fmt.Printf("Total %d combinations to check\n", total)
}
