package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Lens struct {
	Label       string
	FocalLength int
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	boxes := map[int][]Lens{}
	numBoxes := 256
	for i := 0; i < numBoxes; i++ {
		boxes[i] = []Lens{}
	}

	box := 0
	label := ""
	operation := byte(0)
	focalLength := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		line = append(line, ',')
		for _, chr := range line {
			if chr == ',' {
				if operation == '-' {
					for x, lens := range boxes[box] {
						if lens.Label == label {
							boxes[box] = append(boxes[box][:x], boxes[box][x+1:]...)
							break
						}
					}
				} else {
					found := false
					for x, lens := range boxes[box] {
						if lens.Label == label {
							boxes[box][x].FocalLength = focalLength
							found = true
							break
						}
					}
					if !found {
						boxes[box] = append(boxes[box], Lens{
							Label:       label,
							FocalLength: focalLength,
						})
					}
				}

				box = 0
				label = ""
				operation = 0
				focalLength = 0
			} else if chr == '=' || chr == '-' {
				operation = chr
			} else if chr >= '0' && chr <= '9' {
				focalLength = int(chr - byte('0'))
			} else {
				box = ((box + int(chr)) * 17) % 256
				label = label + string(chr)
			}
		}
	}

	total := 0
	for i := 0; i < numBoxes; i++ {
		if len(boxes[i]) > 0 {
			fmt.Printf("Box %d:", i)
		}
		for slot, lens := range boxes[i] {
			power := (i + 1) * (slot + 1) * lens.FocalLength
			total = total + power
			fmt.Printf(" [%s %d] => %d", lens.Label, lens.FocalLength, power)
		}
		if len(boxes[i]) > 0 {
			fmt.Println()
		}
	}
	fmt.Printf("Total power: %d\n", total)
}
