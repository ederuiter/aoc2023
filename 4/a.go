package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strings"
)

func intersection(s1, s2 []string) (inter []string) {
	hash := make(map[string]bool)
	for _, e := range s1 {
		hash[e] = true
	}
	for _, e := range s2 {
		// If elements present in the hashmap then append intersection list.
		if hash[e] {
			inter = append(inter, e)
		}
	}
	//Remove dups from slice.
	inter = removeDups(inter)
	return
}

// Remove dups from slice.
func removeDups(elements []string) (nodups []string) {
	encountered := make(map[string]bool)
	for _, element := range elements {
		if !encountered[element] {
			nodups = append(nodups, element)
			encountered[element] = true
		}
	}
	return
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	r := regexp.MustCompile(" +")
	score := int64(0)
	for scanner.Scan() {
		row := scanner.Text()
		start := strings.Index(row, ":")
		split := strings.Index(row, "|")
		winning := r.Split(strings.Trim(row[start+2:split-1], " "), -1)
		numbers := r.Split(strings.Trim(row[split+2:], " "), -1)
		winners := len(intersection(winning, numbers))
		fmt.Printf("%s => %d winning numbers\n", row, winners)
		if winners > 0 {
			score = score + int64(math.Pow(2, float64(winners-1)))
		}
	}
	fmt.Printf("Total score: %d\n", score)
}
