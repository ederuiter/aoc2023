package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	r := regexp.MustCompile(" +")

	var times []string
	var distances []string
	var score int64 = 1
	scores := []int64{}
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")
		t := parts[0]

		splitted := r.Split(strings.Trim(parts[1], " "), -1)
		if t == "Time" {
			times = splitted
		} else {
			distances = splitted
		}

	}
	for index, timeStr := range times {
		time, _ := strconv.ParseInt(timeStr, 10, 64)
		distance, _ := strconv.ParseInt(distances[index], 10, 64)

		sq := math.Sqrt(float64((time * time) - (4 * distance)))
		x1 := (float64(-time) - sq) / -2
		x2 := (float64(-time) + sq) / -2

		end := int64(math.Ceil(x1)) - 1
		start := int64(math.Floor(x2)) + 1

		num := int64(end-start) + 1
		score = score * num

		fmt.Printf("(time: %d, distance: %d) => (sq: %f, x1: %f, x2: %f)\n", time, distance, sq, x1, x2)

		scores = append(scores, num)
	}
	fmt.Printf("Times: %+v\nScore: %d\nScores: %+v\n", times, score, scores)
}
