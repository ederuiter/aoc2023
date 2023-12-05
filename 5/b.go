package main

import (
	"bufio"
	"cmp"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Range struct {
	from int64
	to   int64
	num  int64
}

func intersect(from Range, to Range) Range {
	fromStart := from.to
	fromEnd := from.to + from.num - 1
	toStart := to.from
	toEnd := to.from + to.num - 1
	if toStart > fromEnd || fromStart > toEnd {
		fmt.Printf("(%+v => %+v) => no intersection (%d > %d || %d > %d)\n", from, to, toStart, fromEnd, fromStart, toEnd)
		return Range{
			from: 0,
			to:   0,
			num:  0,
		}
	} else {

		intersectStart := max(fromStart, toStart)
		intersectEnd := min(fromEnd, toEnd)
		res := Range{
			from: intersectStart,
			to:   intersectStart + (to.to - to.from),
			num:  intersectEnd - intersectStart + 1,
		}
		fmt.Printf("(%+v => %+v) => intersection => %+v\n", from, to, res)
		return res
	}
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	seeds := []string{}
	category := "_"
	nextCategory := map[string]string{"_": "seed"}
	ranges := make(map[string][]Range)
	for scanner.Scan() {
		row := scanner.Text()
		colonPos := strings.Index(row, ":")
		if colonPos >= 0 {
			if row[0:colonPos] == "seeds" {
				seeds = strings.Split(row[colonPos+2:], " ")
				ranges[category] = []Range{}
				for index, seed := range seeds {
					if index%2 == 1 {
						continue
					}
					seedNum, _ := strconv.ParseInt(seed, 10, 64)
					num, _ := strconv.ParseInt(seeds[index+1], 10, 64)
					ranges[category] = append(ranges[category], Range{
						from: seedNum,
						to:   seedNum,
						num:  num,
					})
				}
			} else {
				parts := strings.Split(row[0:colonPos-4], "-to-")
				category = parts[0]
				nextCategory[category] = parts[1]
				ranges[category] = []Range{}
			}
		} else if row != "" {
			parts := strings.Split(row, " ")
			from, _ := strconv.ParseInt(parts[1], 10, 64)
			to, _ := strconv.ParseInt(parts[0], 10, 64)
			num, _ := strconv.ParseInt(parts[2], 10, 64)

			ranges[category] = append(ranges[category], Range{
				from: from,
				to:   to,
				num:  num,
			})
		}
	}

	sortedRanges := make(map[string][]Range)
	category = "_"
	for category != "location" {
		fmt.Println(category)
		slices.SortFunc(ranges[category], func(a Range, b Range) int {
			return cmp.Compare(a.from, b.from)
		})

		next := int64(0)
		res := []Range{}
		for _, currentRange := range ranges[category] {
			if currentRange.from > next {
				res = append(res, Range{
					from: next,
					to:   next,
					num:  currentRange.from - next,
				})
			}
			res = append(res, currentRange)
			next = currentRange.from + currentRange.num
		}
		last := res[len(res)-1].from + res[len(res)-1].num - 1
		res = append(res, Range{
			from: last + 1,
			to:   last + 1,
			num:  math.MaxInt64 - last - 1,
		})

		sortedRanges[category] = res
		category = nextCategory[category]
	}
	fmt.Printf("%+v\n", sortedRanges)
	fmt.Printf("%+v\n", ranges["_"])

	res := int64(math.MaxInt64)
	currentRanges := ranges["_"]
	category = "seed"
	for category != "location" {
		fmt.Printf("Category: %s\n", category)
		newRanges := []Range{}
		for _, currentRange := range currentRanges {
			last := currentRange.to + currentRange.num - 1
			for _, targetRange := range sortedRanges[category] {
				intersection := intersect(currentRange, targetRange)
				if intersection.num > 0 {
					newRanges = append(newRanges, intersection)
				} else if last < targetRange.from+targetRange.num-1 {
					break
				}
			}
		}

		currentRanges = newRanges
		category = nextCategory[category]
	}

	//fmt.Printf("Current ranges: %+v\n", currentRanges)
	for _, currentRange := range currentRanges {
		res = min(res, currentRange.to)
	}

	fmt.Printf("Nearest seed location is: %d\n", res)
}
