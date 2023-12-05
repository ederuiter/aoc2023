package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
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

	seeds := []string{}
	category := "_"
	mapping := make(map[string]map[int64]int64)
	nextCategory := map[string]string{"_": "seed"}
	prevCategory := "_"
	for scanner.Scan() {
		row := scanner.Text()
		colonPos := strings.Index(row, ":")
		if colonPos >= 0 {
			if row[0:colonPos] == "seeds" {
				seeds = strings.Split(row[colonPos+2:], " ")
				mapping["_"] = make(map[int64]int64)
				for _, seed := range seeds {
					seedNum, _ := strconv.ParseInt(seed, 10, 64)
					mapping["_"][seedNum] = seedNum
				}
			} else {
				parts := strings.Split(row[0:colonPos-4], "-to-")
				prevCategory = category
				category = parts[0]
				nextCategory[category] = parts[1]
				mapping[category] = make(map[int64]int64)
				for _, p := range mapping[prevCategory] {
					mapping[category][p] = p
				}
			}
		} else if row != "" {
			parts := strings.Split(row, " ")
			from, _ := strconv.ParseInt(parts[1], 10, 64)
			to, _ := strconv.ParseInt(parts[0], 10, 64)
			num, _ := strconv.ParseInt(parts[2], 10, 64)

			for m, _ := range mapping[category] {
				if m >= from && m < from+num {
					mapping[category][m] = to + m - from
				}
			}
		}
	}

	res := int64(math.MaxInt64)
	for _, seed := range mapping["_"] {
		fmt.Println("---------------")
		index := seed
		category = "seed"
		var ok bool
		var orig int64
		for category != "location" {
			orig = index
			index, ok = mapping[category][index]
			if !ok {
				index = orig
			}
			origCat := category
			category = nextCategory[category]
			fmt.Printf("%s:%d => %s:%d\n", origCat, orig, category, index)

		}
		res = min(res, index)
	}
	fmt.Printf("Mapping: %+v\n", mapping)
	fmt.Printf("Nearest seed location is: %d\n", res)
}
