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

func getPrevious(nums []int64) int64 {
	last := nums
	allZeros := false
	diff := []int64{}
	prev := nums[0]
	dbg := [][]int64{nums}
	i := 1
	for !allZeros {
		allZeros, diff = differences(last)
		dbg = append(dbg, diff)
		prev = prev - diff[0]
		last = diff
		i = i + 1
	}

	p := int64(0)
	for k := i - 1; k >= 0; k-- {
		n := dbg[k][0] - p
		dbg[k] = append([]int64{n}, dbg[k]...)
		p = n
	}

	fmt.Println("-----")
	for _, d := range dbg {
		fmt.Printf("%+v\n", d)
	}

	return dbg[0][0]
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
		prev := getPrevious(line)
		fmt.Printf("Line: %+v => %d\n", line, prev)
		total += prev
	}
	fmt.Printf("Sum of predicted previous values is: %d\n", total)
}
