package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func differences(nums []int64) (bool, []int64) {
	allZeros := true
	prev := nums[0]
	res := []int64{}
	for _, num := range nums[1:] {
		diff := num - prev
		if diff != 0 {
			allZeros = false
		}
		res = append(res, diff)
		prev = num
	}
	return allZeros, res
}

func getNext(nums []int64) int64 {
	prev := nums
	allZeros := false
	diff := []int64{}
	next := nums[len(nums)-1]
	dbg := [][]int64{nums}
	i := 1
	for !allZeros {
		allZeros, diff = differences(prev)
		dbg = append(dbg, diff)
		next = next + diff[len(diff)-1]
		prev = diff
		i = i + 1
	}

	p := int64(0)
	for k := i - 1; k >= 0; k-- {
		n := dbg[k][len(dbg[k])-1] + p
		dbg[k] = append(dbg[k], n)
		p = n
	}

	fmt.Println("-----")
	for _, d := range dbg {
		fmt.Printf("%+v\n", d)
	}

	if dbg[0][len(dbg[0])-1] != next {
		panic("huh?")
	}

	return next
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	total := int64(0)
	for scanner.Scan() {
		lineStr := strings.Split(scanner.Text(), " ")
		line := []int64{}
		for _, numStr := range lineStr {
			num, _ := strconv.ParseInt(numStr, 10, 64)
			line = append(line, num)
		}
		next := getNext(line)
		fmt.Printf("Line: %+v => %d\n", line, next)
		total += next
	}
	fmt.Printf("Sum of predicted next values is: %d\n", total)
}
